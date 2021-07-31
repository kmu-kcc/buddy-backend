package fee

import (
	"context"
	"sort"
	"time"

	"github.com/kmu-kcc/buddy-backend/config"
	"github.com/kmu-kcc/buddy-backend/pkg/member"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Fee struct {
	Year     int                  `json:"year,string" bson:"year"`
	Semester int                  `json:"semester,string" bson:"semester"`
	Amount   int                  `json:"amount,string" bson:"amount"`
	Logs     []primitive.ObjectID `json:"logs" bson:"logs"`
}

func New(year, semester, amount int) *Fee {
	return &Fee{
		Year:     year,
		Semester: semester,
		Amount:   amount,
		Logs:     []primitive.ObjectID{},
	}
}

// Dones returns the list of members who submitted the fee in specific year and semester
//
// NOTE:
//
// It is privileged operation:
//	Only the club managers can access to this operation.
func Dones(year, semester int) (members member.Members, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MongoURI))
	if err != nil {
		return
	}

	fee := new(Fee)
	log := new(Log)
	memb := new(member.Member)

	if err = client.Database("club").
		Collection("fees").
		FindOne(ctx, bson.M{"year": year, "semester": semester}).Decode(fee); err != nil {
		return
	}

	filter := func() bson.D {
		arr := make(bson.A, len(fee.Logs))
		for idx, logID := range fee.Logs {
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
			if fee.Amount <= amount {
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

// Yets returns the list of members who have not yet submitted the fee in specific year and semester
//
// NOTE:
//
// It is privileged operation:
//	Only the club managers can access to this operation.
func Yets(year, semester int) (members member.Members, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MongoURI))
	if err != nil {
		return
	}

	fee := new(Fee)
	log := new(Log)
	memb := new(member.Member)

	if err = client.Database("club").
		Collection("fees").
		FindOne(ctx, bson.M{"year": year, "semester": semester}).Decode(fee); err != nil {
		return
	}

	filter := func() bson.D {
		arr := make(bson.A, len(fee.Logs))
		for idx, logID := range fee.Logs {
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
			if amount < fee.Amount {
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

// All returns the all club fee logs.
//
// NOTE:
//
// It is member-limited operation:
//	Only the authenticated members can access to this operation.
func All(year, semester int) (logs Logs, err error) {
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
		FindOne(ctx, bson.D{
			bson.E{Key: "year", Value: year},
			bson.E{Key: "semester", Value: semester},
		}).Decode(fee); err != nil {
		return
	}

	filter := func() bson.D {
		arr := make(bson.A, len(fee.Logs))
		for idx, logID := range fee.Logs {
			arr[idx] = logID
		}
		return bson.D{
			bson.E{Key: "_id", Value: bson.D{bson.E{Key: "$in", Value: arr}}},
			bson.E{Key: "$or", Value: bson.A{
				bson.D{bson.E{Key: "type", Value: "approved"}},
				bson.D{bson.E{Key: "type", Value: "direct"}}}}}
	}()

	cur, err := client.Database("club").Collection("logs").Find(ctx, filter)
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

	sort.Slice(logs, func(i, j int) bool { return logs[i].UpdatedAt < logs[j].UpdatedAt })

	return logs, client.Disconnect(ctx)
}
