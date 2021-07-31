// Package fee provides access to the club fee of the Buddy System.
package fee

import (
	"context"
	"errors"
	"time"

	"github.com/kmu-kcc/buddy-backend/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ErrDuplicatedFee = errors.New("duplicated fee")
)

// Fee represents a club fee state.
type Fee struct {
	Year     int                  `json:"year,string" bson:"year"`
	Semester int                  `json:"semester,string" bson:"semester"`
	Amount   int                  `json:"amount,string" bson:"amount"`
	Logs     []primitive.ObjectID `json:"logs" bson:"logs"`
}

// New returns a new club fee.
func New(year, semester, amount int) *Fee {
	return &Fee{
		Year:     year,
		Semester: semester,
		Amount:   amount,
		Logs:     []primitive.ObjectID{},
	}
}

// Create creats a new fees history.
//
// NOTE:
//
// It is privileged operation:
//	Only the club managers can access to this operation.
func Create(year, semester, amount int) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MongoURI))
	if err != nil {
		return
	}

	collection := client.Database("club").Collection("fees")
	fee := new(Fee)

	if err = collection.FindOne(ctx, bson.D{
		bson.E{Key: "year", Value: year},
		bson.E{Key: "semester", Value: semester},
	}).Decode(fee); err != mongo.ErrNoDocuments {
		if err = client.Disconnect(ctx); err != nil {
			return
		}
		return ErrDuplicatedFee
	}

	if _, err = collection.InsertOne(ctx, *New(year, semester, amount)); err != nil {
		return
	}

	return client.Disconnect(ctx)
}

// Submit creates fees payment application log.
//
// NOTE:
//
// It is member-limited operation:
//	Only the authenticated members can access to this operation.
func Submit(memberID string, year, semester, amount int) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MongoURI))
	if err != nil {
		return
	}

	feeCollection := client.Database("club").Collection("fees")
	logCollection := client.Database("club").Collection("logs")

	fee := new(Log)
	log := NewLog(memberID, "unapproved", amount)

	if _, err = logCollection.InsertOne(ctx, *log); err != nil {
		return
	}

	if err = feeCollection.FindOneAndUpdate(
		ctx,
		bson.D{
			bson.E{Key: "year", Value: year},
			bson.E{Key: "semester", Value: semester},
		}, bson.D{
			bson.E{Key: "$push", Value: bson.D{
				bson.E{Key: "logs", Value: log.ID},
			}},
		}).Decode(fee); err != nil {
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
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
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

	return sum, client.Disconnect(ctx)
}
