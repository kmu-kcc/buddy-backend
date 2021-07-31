// Package fee provides access to the fee of club of the Buddy System.
package fee

import (
	"context"
	"time"

	"github.com/kmu-kcc/buddy-backend/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Fee represents a club fee.
type Fee struct {
	Year     int                  `json:"year,string" bson:"year"`
	Semester int                  `json:"semester,string" bson:"semester"`
	Amount   int                  `json:"amount,string" bson:"amount"`
	Logs     []primitive.ObjectID `json:"logs" bson:"logs"`
}

// New returns a new fee
func New(year, semester, amount int, logs []primitive.ObjectID) Fee {
	return Fee{
		Year:     year,
		Semester: semester,
		Amount:   amount,
		Logs:     logs,
	}
}

// Approve approves the requests of ids & changes type from unapproved to approved
// Note :
// This is privileged operation:
// 	Only the club managers can access to this operation
func Approve(ids []primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MongoURI))

	if err != nil {
		return err
	}

	// Update logs to approved & renew UpdatedAt with current time.
	filter := func() bson.D {
		arr := make(bson.A, len(ids))
		for idx, id := range ids {
			arr[idx] = id
		}
		return bson.D{bson.E{Key: "_id", Value: bson.D{bson.E{Key: "$in", Value: arr}}}}
	}()

	if _, err = client.Database("club").
		Collection("fee").
		UpdateMany(
			ctx,
			filter,
			bson.D{
				bson.E{Key: "$set", Value: bson.D{
					bson.E{Key: "type", Value: "approved"},
					bson.E{Key: "updatedat", Value: time.Now().Unix()}}}}); err != nil {
		return err
	}

	return client.Disconnect(ctx)
}

// Reject rejects the requests of ids & remove request of id from logs
// Note :
// This is privileged operation:
// 	Only the club managers can access to this operation
func Reject(ids []primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MongoURI))
	if err != nil {
		return err
	}

	for _, id := range ids {
		if _, err := client.Database("club").Collection("fee").UpdateByID(ctx, id, bson.M{"$pull": bson.M{"logs": id}}); err != nil {
			return err
		}
	}

	return client.Disconnect(ctx)
}

// Deposit makes a new log with amount and append it to fee with Year  of year, Semester of semester
// Note :
// This is privileged operation:
// 	Only the club managers can access to this operation
func Deposit(year, semester, amount int) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MongoURI))
	if err != nil {
		return err
	}

	deposit := NewLog("", "direct", amount)
	// deposit := new(Log)
	// deposit.Type = "direct"
	// deposit.Amount = amount
	// deposit.CreatedAt = time.Now().Unix()
	// deposit.UpdatedAt = time.Now().Unix()

	// client.Database("club").Collection("fee").FindOneAndUpdate(ctx,
	// 	bson.D{
	// 		bson.E{Key: "year", Value: year},
	// 		bson.E{Key: "semester", Value: semester},
	// 	},
	// 	bson.D{
	// 		bson.E{Key: "$push", Value: bson.D{
	// 			bson.E{Key: "logs", Value: deposit},
	// 		}},
	// 	})
	if _, err := client.Database("club").Collection("fee").UpdateOne(ctx,
		bson.D{
			bson.E{Key: "year", Value: year},
			bson.E{Key: "semester", Value: semester},
		},
		bson.D{
			bson.E{Key: "$push", Value: bson.D{
				bson.E{Key: "logs", Value: deposit.ID},
			}},
		}); err != nil {
		return err
	}
	return client.Disconnect(ctx)
}
