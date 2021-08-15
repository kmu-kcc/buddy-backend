// Package fee provides access to the club fee of the Buddy System.
package fee

import (
	"context"
	"errors"

	// "sort"
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

var ErrDuplicatedFee = errors.New("duplicated fee")

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
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MongoURI))
	if err != nil {
		return
	}

	collection := client.Database("club").Collection("fees")
	fee := new(Fee)

	if err = collection.FindOne(ctx, bson.D{
		bson.E{Key: "year", Value: f.Year},
		bson.E{Key: "semester", Value: f.Semester},
	}).Decode(fee); err != mongo.ErrNoDocuments {
		if err = client.Disconnect(ctx); err != nil {
			return
		}
		return ErrDuplicatedFee
	}

	if _, err = collection.InsertOne(ctx, f); err != nil {
		return
	}

	return client.Disconnect(ctx)
}

// Amount finds log by year and semester, and returns the sum of all amounts using memberID and type.
//
// NOTE:
//
// It is member-limited operation:
//	Only the authenticated members can access to this operation.
func Amount(year, semester int, memberID string) (sum int, err error) {
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

	cur, err := client.Database("club").
		Collection("logs").
		Find(ctx, bson.M{
			"member_id": memberID,
			"type":      "approved",
		})
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

// Dones returns the list of members who submitted the fee in specific year and semester.
//
// NOTE:
//
// It is privileged operation:
//	Only the club managers can access to this operation.
func (f *Fee) Dones() (members member.Members, err error) {
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
			bson.E{Key: "type", Value: "approved"},
		}
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
		return bson.D{bson.E{Key: "id", Value: bson.D{bson.E{Key: "$in", Value: arr}}},
			bson.E{Key: "attendance", Value: bson.D{bson.E{Key: "$ne", Value: 2}}}}
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

// Yets returns the list of members who have not yet submitted the fee in specific year and semester.
//
// NOTE:
//
// It is privileged operation:
//	Only the club managers can access to this operation.
func (f *Fee) Yets() (members member.Members, err error) {
	dones, err := f.Dones()
	if err != nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MongoURI))
	if err != nil {
		return
	}

	filter := func() bson.D {
		ids := make(bson.A, len(dones))
		for idx, memb := range dones {
			ids[idx] = memb.ID
		}
		return bson.D{
			bson.E{Key: "id", Value: bson.D{bson.E{Key: "$nin", Value: ids}}},
			bson.E{Key: "attendance", Value: bson.D{bson.E{Key: "$ne", Value: 2}}}}
	}()

	memb := new(member.Member)

	cur, err := client.Database("club").Collection("members").Find(ctx, filter)
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

// All returns the all club fee logs.
//
// NOTE:
//
// It is member-limited operation:
//	Only the authenticated members can access to this operation.
func All(startdate, enddate int) (logs Logs, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MongoURI))
	if err != nil {
		return
	}

	log := new(Log)

	cur, err := client.Database("club").Collection("logs").Find(ctx, bson.D{bson.E{Key: "$and",
		Value: bson.A{
			bson.D{bson.E{Key: "updated_at", Value: bson.D{bson.E{Key: "$lte", Value: enddate}}}},
			bson.D{bson.E{Key: "updated_at", Value: bson.D{bson.E{Key: "$gte", Value: startdate}}}}}}})
	if err != nil {
		return
	}

	for cur.Next(ctx) {
		if err = cur.Decode(log); err != nil {
			return
		}
		logs = append(logs, *log)
	}
	if err = cur.Close(ctx); err != nil {
		return
	}

	// sort.Slice(logs, func(i, j int) bool { return logs[i].UpdatedAt < logs[j].UpdatedAt })

	return logs, client.Disconnect(ctx)
}

// Pay makes payment
//
// Note:
//
// This is privileged operation:
// 	Only the club managers can access to this operation.

func Pay(year, sem int, ids []string, amounts []int) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MongoURI))
	if err != nil {
		return err
	}
	// bson.D for insertMany
	insertDoc := func() []interface{} {
		ans := make([]interface{}, len(amounts))
		for i := 0; i < len(amounts) && i < len(ids); i++ {
			tmp := bson.D{
				bson.E{Key: "member_id", Value: ids[i]},
				bson.E{Key: "amount", Value: amounts[i]},
			}
			ans[i] = tmp
		}
		return ans
	}()

	// fmt.Print(insertDoc)
	// logs insertMany
	if _, err := client.Database("club").Collection("logs").InsertMany(ctx, insertDoc); err != nil {
		// fmt.Print("Log Insertion")
		return err
	}

	targetID := make(bson.A, len(ids))
	for idx, id := range ids {
		targetID[idx] = id
	}
	// Fee of year,sem .logs +
	// fmt.Print("feee Insertion")
	if _, err = client.Database("club").Collection("fees").UpdateMany(ctx, bson.M{"year": year, "semester": sem},
		bson.D{
			bson.E{Key: "$push", Value: bson.D{
				bson.E{Key: "logs", Value: bson.D{
					bson.E{Key: "_id", Value: bson.D{
						bson.E{Key: "$in", Value: bson.D{
							bson.E{Key: "member_id", Value: targetID}}}}}}}}}}); err != nil {
		return err
	}
	return client.Disconnect(ctx)
}

// Deposit makes a new log with amount and append it to fee with Year  of year, Semester of semester
//
// Note :
//
// This is privileged operation:
// 	Only the club managers can access to this operation
func Deposit(year, semester, amount int) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MongoURI))
	if err != nil {
		return err
	}

	depo := NewLog("", amount, 1)

	if _, err = client.Database("club").
		Collection("fees").
		UpdateOne(ctx,
			bson.D{
				bson.E{Key: "year", Value: year},
				bson.E{Key: "semester", Value: semester},
			},
			bson.D{
				bson.E{Key: "$push", Value: bson.D{
					bson.E{Key: "logs", Value: depo.ID},
				}},
			}); err != nil {
		return err
	}

	if _, err = client.Database("club").
		Collection("logs").
		InsertOne(ctx, depo); err != nil {
		return err
	}
	return client.Disconnect(ctx)
}
