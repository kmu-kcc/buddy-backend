// Package fee provides access to the club fee of the Buddy System.
package fee

import (
	"context"
	"errors"
	"sort"
	"time"

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

// Create creates a new fees history.
//
// NOTE:
//
// It is privileged operation:
//	Only the club managers can access to this operation.
func (f Fee) Create() (err error) {
	var year, semester int

	if f.Semester == 1 {
		year, semester = f.Year-1, 2
	} else {
		year, semester = f.Year, 1
	}

	f.CarryOver, _, _, err = New(year, semester, 0, 0).Search()
	if err != nil {
		return
	}

	f.Logs = []primitive.ObjectID{}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MongoURI))
	if err != nil {
		return
	}

	collection := client.Database("club").Collection("fees")
	fee := new(Fee)

	if err = collection.FindOne(
		ctx,
		bson.D{
			bson.E{Key: "year", Value: f.Year},
			bson.E{Key: "semester", Value: f.Semester},
		}).
		Decode(fee); err == mongo.ErrNoDocuments {
		if _, err = collection.InsertOne(ctx, f); err != nil {
			return
		}
		return client.Disconnect(ctx)
	} else if err == nil {
		if err = client.Disconnect(ctx); err != nil {
			return
		}
		return ErrDuplicatedFee
	}
	return
}

// Amount returns the sum of payments of member of memberID.
//
// NOTE:
//
// It is member-limited operation:
//	Only the authenticated members can access to this operation.
func Amount(year, semester int, id string) (sum int, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MongoURI))
	if err != nil {
		return
	}

	fee := new(Fee)
	log := new(Log)

	if err = client.Database("club").
		Collection("fees").
		FindOne(ctx, bson.M{
			"year":     year,
			"semester": semester,
		}).Decode(fee); err != nil {
		return
	}

	filter := bson.D{
		bson.E{Key: "_id", Value: bson.D{
			bson.E{Key: "$in", Value: fee.Logs},
		}},
		bson.E{Key: "member_id", Value: id},
		bson.E{Key: "type", Value: payment},
	}

	cur, err := client.Database("club").
		Collection("logs").
		Find(ctx, filter)

	if err != nil {
		return
	}

	for cur.Next(ctx) {
		if err = cur.Decode(log); err != nil {
			return
		}
		sum += log.Amount
	}

	if err = cur.Close(ctx); err != nil {
		return
	}

	return sum, client.Disconnect(ctx)
}

// Payers returns the list of members who paid the fee of year and semester.
//
// NOTE:
//
// It is privileged operation:
//	Only the club managers can access to this operation.
func (f *Fee) Payers() (members member.Members, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MongoURI))
	if err != nil {
		return
	}

	log := new(Log)
	memb := new(member.Member)

	if err = client.Database("club").
		Collection("fees").
		FindOne(ctx, bson.M{"year": f.Year, "semester": f.Semester}).Decode(f); err != nil {
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
	if err = cur.Close(ctx); err != nil {
		return
	}

	return members, client.Disconnect(ctx)
}

// Deptors returns the list of members who did not pay the fee of year and semester.
//
// NOTE:
//
// It is privileged operation:
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
	ids = append(ids, "MASTER")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MongoURI))
	if err != nil {
		return
	}

	memb := new(member.Member)
	cur, err := client.Database("club").Collection("members").Find(ctx, bson.D{
		bson.E{Key: "attendance", Value: bson.D{bson.E{Key: "$ne", Value: member.Graduate}}},
		bson.E{Key: "id", Value: bson.D{bson.E{Key: "$nin", Value: ids}}},
	})
	if err != nil {
		return
	}

	for cur.Next(ctx) {
		if err = cur.Decode(memb); err != nil {
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
		return
	}

	for cur.Next(ctx) {
		if err = cur.Decode(log); err != nil {
			return
		}
		amounts[log.MemberID] += log.Amount
	}
	if err = cur.Close(ctx); err != nil {
		return
	}

	depts = make([]int, len(deptors))
	for idx, deptor := range deptors {
		depts[idx] = f.Amount - amounts[deptor.ID]
	}
	return
}

// Search returns the fee history of year and semester.
//
// NOTE:
//
// It is member-limited operation:
//	Only the authenticated members can access to this operation.
func (f *Fee) Search() (carryOver int, _ []map[string]interface{}, total int, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MongoURI))
	if err != nil {
		return
	}

	if err = client.Database("club").Collection("fees").FindOne(ctx,
		bson.D{
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
		return
	}

	total = f.CarryOver
	log := new(Log)
	var logs Logs

	for cur.Next(ctx) {
		if err = cur.Decode(log); err != nil {
			return
		}
		logs = append(logs, *log)
		total += log.Amount
	}
	if err = cur.Close(ctx); err != nil {
		return
	}

	sort.Slice(logs, func(i, j int) bool { return logs[i].CreatedAt < logs[j].CreatedAt })

	return f.CarryOver, logs.Public(), total, client.Disconnect(ctx)
}

// Pay registers payments of members of ids for each amount of amounts.
//
// Note:
//
// This is privileged operation:
// 	Only the club managers can access to this operation.
func Pay(year, semester int, ids []string, amounts []int) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MongoURI))
	if err != nil {
		return err
	}

	if _, err = client.Database("club").
		Collection("fees").
		Find(ctx, bson.M{"year": year, "semester": semester}); err != nil {
		return err
	}

	docs := func() bson.A {
		logs := make(bson.A, len(ids))
		for idx, id := range ids {
			logs[idx] = NewLog(id, "회비 납부", amounts[idx], payment)
		}
		return logs
	}()

	if _, err = client.Database("club").
		Collection("logs").
		InsertMany(ctx, docs); err != nil {
		return err
	}

	cur, err := client.Database("club").
		Collection("logs").
		Find(ctx,
			bson.D{
				bson.E{Key: "member_id", Value: bson.D{
					bson.E{Key: "$in", Value: ids}}}})

	if err != nil {
		return err
	}

	log := new(Log)
	logs := bson.A{}

	for i := 0; cur.Next(ctx); i++ {
		if err = cur.Decode(log); err != nil {
			return err
		}
		logs = append(logs, log.ID)
	}

	if err = cur.Close(ctx); err != nil {
		return err
	}

	if _, err = client.Database("club").
		Collection("fees").
		UpdateOne(ctx,
			bson.M{"year": year, "semester": semester}, bson.D{
				bson.E{Key: "$push", Value: bson.D{
					bson.E{Key: "logs", Value: bson.D{
						bson.E{Key: "$each", Value: logs}}}}}}); err != nil {
		return err
	}
	return client.Disconnect(ctx)
}

// Deposit makes a new log with amount and append it to fee with Year  of year, Semester of semester
//
// Note:
//
// This is privileged operation:
// 	Only the club managers can access to this operation
func Deposit(year, semester, amount int, description string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MongoURI))
	if err != nil {
		return err
	}

	log := NewLog("", description, amount, deposit)

	if _, err = client.Database("club").
		Collection("fees").
		UpdateOne(ctx,
			bson.D{
				bson.E{Key: "year", Value: year},
				bson.E{Key: "semester", Value: semester},
			},
			bson.D{
				bson.E{Key: "$push", Value: bson.D{
					bson.E{Key: "logs", Value: log.ID},
				}},
			}); err != nil {
		return err
	}

	if _, err = client.Database("club").
		Collection("logs").
		InsertOne(ctx, log); err != nil {
		return err
	}
	return client.Disconnect(ctx)
}

// Exempt exempts the member of id from the fee of year and semester.
//
// Note :
//
// This is a privileged operation:
// 	Only the club managers can access to this operation
func (f *Fee) Exempt(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MongoURI))
	if err != nil {
		return err
	}

	db := client.Database("club")
	feeCollection := db.Collection("fees")
	logCollection := db.Collection("logs")

	if err = feeCollection.FindOne(
		ctx,
		bson.D{
			bson.E{Key: "year", Value: f.Year},
			bson.E{Key: "semester", Value: f.Semester},
		}).
		Decode(f); err != nil {
		return err
	}

	log := new(Log)
	filter := bson.D{
		bson.E{Key: "_id", Value: bson.D{
			bson.E{Key: "$in", Value: f.Logs},
		}},
		bson.E{Key: "member_id", Value: id},
		bson.E{Key: "type", Value: exemption},
	}

	if err = logCollection.FindOne(ctx, filter).Decode(log); err == nil {
		if err = client.Disconnect(ctx); err != nil {
			return err
		}
		return ErrAlreadyExempted
	} else if err == mongo.ErrNoDocuments {
		log = NewLog(id, "회비 면제", f.Amount, exemption)
		if _, err = feeCollection.UpdateOne(
			ctx,
			bson.D{
				bson.E{Key: "year", Value: f.Year},
				bson.E{Key: "semester", Value: f.Semester}},
			bson.D{
				bson.E{Key: "$push", Value: bson.D{
					bson.E{Key: "logs", Value: log.ID},
				}},
			}); err != nil {
			return err
		}
	} else {
		return err
	}

	if _, err = logCollection.InsertOne(ctx, log); err != nil {
		return err
	}

	return client.Disconnect(ctx)
}
