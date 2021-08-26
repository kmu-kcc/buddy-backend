// Package fee provides access to the club fee of the Buddy System.
package fee

import (
	"context"
	"errors"
	"sort"

	"github.com/kmu-kcc/buddy-backend/config"
	"github.com/kmu-kcc/buddy-backend/pkg/member"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	payment = iota
	deposit
	exemption
)

var (
	ErrDuplicatedFee   = errors.New("duplicated fee")
	ErrAlreadyExempted = errors.New("already exempted")
)

// Fee represents a club fee state.
type Fee struct {
	Year      int                  `json:"year" bson:"year"`
	Semester  int                  `json:"semester" bson:"semester"`
	CarryOver int                  `json:"carry_over" bson:"carry_over"`
	Amount    int                  `json:"amount" bson:"amount"`
	Logs      []primitive.ObjectID `json:"logs" bson:"logs"`
}

// New returns a new club fee.
func New(year, semester, carryOver, amount int) *Fee {
	return &Fee{
		Year:      year,
		Semester:  semester,
		CarryOver: carryOver,
		Amount:    amount,
		Logs:      []primitive.ObjectID{},
	}
}

// Create creates a new fee history.
//
// NOTE:
//
// It is a privileged operation:
//	Only the club managers can access to this operation.
func (f Fee) Create() (err error) {
	var year, semester int

	if f.Semester == 1 {
		year, semester = f.Year-1, 2
	} else {
		year, semester = f.Year, 1
	}

	_, _, f.CarryOver, err = New(year, semester, 0, 0).Search()
	f.Logs = []primitive.ObjectID{}
	if err != nil {
		return
	}

	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MongoURI))
	if err != nil {
		return
	}
	defer client.Disconnect(ctx)

	collection := client.Database("club").Collection("fees")
	fee := new(Fee)

	if err = collection.FindOne(ctx, bson.D{bson.E{Key: "year", Value: f.Year}, bson.E{Key: "semester", Value: f.Semester}}).Decode(fee); err == mongo.ErrNoDocuments {
		_, err = collection.InsertOne(ctx, f)
		return
	} else if err == nil {
		return ErrDuplicatedFee
	}
	return
}

// Amount returns the amount of payments of member of id.
//
// NOTE:
//
// It is a member-limited operation:
//	Only the authenticated members can access to this operation.
func Amount(year, semester int, id string) (amount int, err error) {
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MongoURI))
	if err != nil {
		return
	}
	defer client.Disconnect(ctx)

	fee := new(Fee)
	log := new(Log)

	if err = client.Database("club").Collection("fees").FindOne(ctx, bson.M{"year": year, "semester": semester}).Decode(fee); err != nil {
		return
	}

	filter := bson.D{
		bson.E{Key: "_id", Value: bson.D{
			bson.E{Key: "$in", Value: fee.Logs},
		}},
		bson.E{Key: "member_id", Value: id},
		bson.E{Key: "type", Value: payment},
	}

	cur, err := client.Database("club").Collection("logs").Find(ctx, filter)
	if err != nil {
		return
	}

	for cur.Next(ctx) {
		if err = cur.Decode(log); err != nil {
			return
		}
		amount += log.Amount
	}

	return amount, cur.Close(ctx)
}

// Payers returns the list of members who paid the fee of year and semester.
//
// NOTE:
//
// It is a privileged operation:
//	Only the club managers can access to this operation.
func (f *Fee) Payers() (members member.Members, err error) {
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MongoURI))
	if err != nil {
		return
	}
	defer client.Disconnect(ctx)

	log := new(Log)
	memb := new(member.Member)

	if err = client.Database("club").Collection("fees").FindOne(ctx, bson.M{"year": f.Year, "semester": f.Semester}).Decode(f); err != nil {
		return
	}

	filter := func() bson.D {
		arr := make(bson.A, len(f.Logs))
		for idx, logID := range f.Logs {
			arr[idx] = logID
		}
		return bson.D{
			bson.E{Key: "_id", Value: bson.D{bson.E{Key: "$in", Value: arr}}},
			bson.E{Key: "$or", Value: bson.A{
				bson.D{bson.E{Key: "type", Value: payment}},
				bson.D{bson.E{Key: "type", Value: exemption}}}}}
	}()

	cur, err := client.Database("club").Collection("logs").Find(ctx, filter)
	if err != nil {
		return
	}

	amounts := make(map[string]int)

	for cur.Next(ctx) {
		if err = cur.Decode(log); err != nil {
			return
		}
		amounts[log.MemberID] += log.Amount
	}

	if err = cur.Close(ctx); err != nil {
		return
	}

	filter = func() bson.D {
		arr := bson.A{}
		for membID, amount := range amounts {
			if f.Amount <= amount {
				arr = append(arr, membID)
			}
		}
		return bson.D{bson.E{Key: "id", Value: bson.D{bson.E{Key: "$in", Value: arr}}}}
	}()

	cur, err = client.Database("club").Collection("members").Find(ctx, filter)
	if err != nil {
		return
	}

	for cur.Next(ctx) {
		if err = cur.Decode(memb); err != nil {
			return
		}
		members = append(members, *memb)
	}

	return members, cur.Close(ctx)
}

// Deptors returns the list of members who did not pay the fee of year and semester.
//
// NOTE:
//
// It is a privileged operation:
//	Only the club managers can access to this operation.
func (f *Fee) Deptors() (deptors member.Members, depts []int, err error) {
	payers, err := f.Payers()
	if err != nil {
		return
	}

	ids := make(bson.A, len(payers)+1)
	for idx, payer := range payers {
		ids[idx] = payer.ID
	}
	ids = append(ids, member.MASTER)

	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MongoURI))
	if err != nil {
		return
	}
	defer client.Disconnect(ctx)

	memb := new(member.Member)
	cur, err := client.Database("club").Collection("members").Find(ctx, bson.D{
		bson.E{Key: "attendance", Value: bson.D{bson.E{Key: "$ne", Value: member.Graduate}}},
		bson.E{Key: "id", Value: bson.D{bson.E{Key: "$nin", Value: ids}}},
	})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			err = nil
		}
		return
	}

	for cur.Next(ctx) {
		if err = cur.Decode(memb); err != nil {
			if err == mongo.ErrNoDocuments {
				err = nil
			}
			return
		}
		deptors = append(deptors, *memb)
	}
	if err = cur.Close(ctx); err != nil {
		return
	}

	ids = make(bson.A, len(deptors))
	for idx, deptor := range deptors {
		ids[idx] = deptor.ID
	}

	log := new(Log)
	amounts := make(map[string]int)
	cur, err = client.Database("club").Collection("logs").Find(ctx, bson.D{
		bson.E{Key: "_id", Value: bson.D{bson.E{Key: "$in", Value: f.Logs}}},
		bson.E{Key: "type", Value: payment},
		bson.E{Key: "member_id", Value: bson.D{bson.E{Key: "$in", Value: ids}}},
	})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			err = nil
		}
		return
	}

	for cur.Next(ctx) {
		if err = cur.Decode(log); err != nil {
			if err == mongo.ErrNoDocuments {
				err = nil
			}
			return
		}
		amounts[log.MemberID] += log.Amount
	}

	depts = make([]int, len(deptors))
	for idx, deptor := range deptors {
		depts[idx] = f.Amount - amounts[deptor.ID]
	}

	return deptors, depts, cur.Close(ctx)
}

// Search returns the fee history of year and semester.
//
// NOTE:
//
// It is a member-limited operation:
//	Only the authenticated members can access to this operation.
func (f *Fee) Search() (carryOver int, _ []map[string]interface{}, total int, err error) {
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MongoURI))
	if err != nil {
		return
	}
	defer client.Disconnect(ctx)

	if err = client.Database("club").Collection("fees").FindOne(ctx, bson.D{
		bson.E{Key: "year", Value: f.Year},
		bson.E{Key: "semester", Value: f.Semester}}).Decode(f); err == mongo.ErrNoDocuments {
		return 0, Logs{}.Public(), 0, nil
	} else if err != nil {
		return
	}

	filter := func() bson.D {
		arr := make(bson.A, len(f.Logs))
		for idx, logID := range f.Logs {
			arr[idx] = logID
		}
		return bson.D{
			bson.E{Key: "_id", Value: bson.D{bson.E{Key: "$in", Value: arr}}},
			bson.E{Key: "$or", Value: bson.A{
				bson.D{bson.E{Key: "type", Value: payment}},
				bson.D{bson.E{Key: "type", Value: deposit}}}}}
	}()

	cur, err := client.Database("club").Collection("logs").Find(ctx, filter)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			err = nil
		}
		return
	}

	total = f.CarryOver
	log := new(Log)
	var logs Logs

	for cur.Next(ctx) {
		if err = cur.Decode(log); err != nil {
			if err == mongo.ErrNoDocuments {
				err = nil
			}
			return
		}
		logs = append(logs, *log)
		total += log.Amount
	}

	sort.Slice(logs, func(i, j int) bool { return logs[i].CreatedAt < logs[j].CreatedAt })

	return f.CarryOver, logs.Public(), total, cur.Close(ctx)
}

// Pay registers payments of members of ids for each amount of amounts.
//
// Note:
//
// It is a privileged operation:
// 	Only the club managers can access to this operation.
func Pay(year, semester int, ids []string, amounts []int) error {
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MongoURI))
	if err != nil {
		return err
	}
	defer client.Disconnect(ctx)

	fee := new(Fee)
	if err = client.Database("club").Collection("fees").FindOne(ctx, bson.M{"year": year, "semester": semester}).Decode(fee); err != nil {
		return err
	}

	logs := make(bson.A, len(ids))
	logIDs := make([]primitive.ObjectID, len(ids))
	for idx, id := range ids {
		log := NewLog(id, "회비 납부", amounts[idx], payment)
		logs[idx] = log
		logIDs = append(logIDs, log.ID)
	}

	if _, err = client.Database("club").Collection("logs").InsertMany(ctx, logs); err != nil {
		return err
	}

	_, err = client.Database("club").Collection("fees").UpdateOne(ctx, bson.M{"year": year, "semester": semester}, bson.D{
		bson.E{Key: "$push", Value: bson.D{
			bson.E{Key: "logs", Value: bson.D{
				bson.E{Key: "$each", Value: logIDs}}}}}})
	return err
}

// Deposit makes a new log with amount and append it to fee with year of YEAR, semester of SEMESTER.
//
// Note:
//
// It is a privileged operation:
// 	Only the club managers can access to this operation.
func Deposit(year, semester, amount int, description string) error {
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MongoURI))
	if err != nil {
		return err
	}
	defer client.Disconnect(ctx)

	log := NewLog("", description, amount, deposit)

	if _, err = client.Database("club").
		Collection("logs").
		InsertOne(ctx, log); err != nil {
		return err
	}

	_, err = client.Database("club").Collection("fees").UpdateOne(ctx,
		bson.D{bson.E{Key: "year", Value: year}, bson.E{Key: "semester", Value: semester}},
		bson.D{bson.E{Key: "$push", Value: bson.D{bson.E{Key: "logs", Value: log.ID}}}})
	return err
}

// Exempt exempts the member of id from the fee of year and semester.
//
// Note :
//
// It is a privileged operation:
// 	Only the club managers can access to this operation.
func (f *Fee) Exempt(id string) error {
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MongoURI))
	if err != nil {
		return err
	}
	defer client.Disconnect(ctx)

	db := client.Database("club")
	feeCollection := db.Collection("fees")
	logCollection := db.Collection("logs")

	if err = feeCollection.FindOne(ctx, bson.D{bson.E{Key: "year", Value: f.Year}, bson.E{Key: "semester", Value: f.Semester}}).Decode(f); err != nil {
		return err
	}

	log := new(Log)
	filter := bson.D{
		bson.E{Key: "_id", Value: bson.D{bson.E{Key: "$in", Value: f.Logs}}},
		bson.E{Key: "member_id", Value: id},
		bson.E{Key: "type", Value: exemption},
	}

	if err = logCollection.FindOne(ctx, filter).Decode(log); err == nil {
		return ErrAlreadyExempted
	}
	if err == mongo.ErrNoDocuments {
		log = NewLog(id, "회비 면제", f.Amount, exemption)
		if _, err = logCollection.InsertOne(ctx, log); err != nil {
			return err
		}
		_, err = feeCollection.UpdateOne(ctx,
			bson.D{bson.E{Key: "year", Value: f.Year}, bson.E{Key: "semester", Value: f.Semester}},
			bson.D{bson.E{Key: "$push", Value: bson.D{bson.E{Key: "logs", Value: log.ID}}}})
		return err
	}
	return err
}
