package fee

import (
	"context"
	"time"

	"github.com/kmu-kcc/buddy-backend/config"
	"github.com/kmu-kcc/buddy-backend/pkg/member"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Log struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	MemberID  string             `json:"memberID" bson:"id"`
	Amount    int                `json:"amount,string" bson:"amount"`
	Type      string             `json:"type" bson:"type"`
	CreatedAt int64              `json:"creatat,string" bson:"creatat"`
	UpdatedAt int64              `json:"updateat,string" bson:"creatat"`
}

type Fee struct {
	Year     int                  `json:"year,string" bson:"year"`
	Semester int                  `json:"semstter,string" bson:"semester"`
	Amount   int                  `json:"amount,string" bson:"amount"`
	Logs     []primitive.ObjectID `json:"logs" bson:"logs"`
}

func New(a string, b int, c string) *Log {
	return &Log{
		ID:        primitive.NewObjectID(),
		MemberID:  a,
		Amount:    b,
		Type:      c,
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}
}

// Donse returns the list of members who submitted the fee in specific year and semester
//
// NOTE:
//
// It is privileged operation:
//	Only the club managers can access to this operation.
func Dones(year, semester int) (members []member.Member, err error) {
	ctx, cancle := context.WithTimeout(context.Background(), time.Second)
	defer cancle()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MongoURI))
	if err != nil {
		return
	}

	collection := client.Database("club").Collection("fees")

	fee := new(Fee)
	log := new(Log)
	memberi := new(member.Member)

	if err = collection.FindOne(ctx, bson.M{"year": year, "semester": semester}).Decode(fee); err == mongo.ErrNoDocuments {
		return
	} else if err != nil {
		return
	}

	filter := func() bson.D {
		arr := make(bson.A, len(fee.Logs))
		for idx, pid := range fee.Logs {
			arr[idx] = pid
		}
		return bson.D{bson.E{Key: "pid", Value: bson.D{bson.E{Key: "$in", Value: arr}}}}
	}()
	logs := make([]Log, len(fee.Logs))
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

	memberg := make([]member.Member, 100)
	cur1, err := client.Database("club").Collection("members").Find(ctx, bson.M{})
	if err != nil {
		return
	}
	for cur.Next(ctx) {
		if err = cur1.Decode(memberi); err != nil {
			return
		}
		memberg = append(memberg, *memberi)
	}

	mam := make(map[string]int)
	for _, meb := range memberg {
		mam[meb.ID] = 0
	}
	for _, log := range logs {
		if log.Type == "approved" {
			mam[log.MemberID] += log.Amount
		}
	}
	for _, meb := range memberg {
		if mam[meb.ID] >= fee.Amount {
			members = append(members, meb)
		}
	}
	return members, client.Disconnect(ctx)
}

// Yets returns the list of members who have not yet submitted the fee in specific year and semester
//
// NOTE:
//
// It is privileged operation:
//	Only the club managers can access to this operation.
func Yets(year, semester int) (members []member.Member, err error) {
	ctx, cancle := context.WithTimeout(context.Background(), time.Second)
	defer cancle()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MongoURI))
	if err != nil {
		return
	}

	collection := client.Database("club").Collection("fees")

	fee := new(Fee)
	log := new(Log)
	memberi := new(member.Member)

	if err = collection.FindOne(ctx, bson.M{
		"year":     year,
		"semester": semester,
	}).Decode(fee); err == mongo.ErrNoDocuments {
		return
	} else if err != nil {
		return
	}

	filter := func() bson.D {
		arr := make(bson.A, len(fee.Logs))
		for idx, pid := range fee.Logs {
			arr[idx] = pid
		}
		return bson.D{bson.E{Key: "pid", Value: bson.D{bson.E{Key: "$in", Value: arr}}}}
	}()
	logs := make([]Log, len(fee.Logs))
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

	memberg := make([]member.Member, 100)
	cur1, err := client.Database("club").Collection("members").Find(ctx, bson.M{})
	if err != nil {
		return
	}
	for cur.Next(ctx) {
		if err = cur1.Decode(memberi); err != nil {
			return
		}
		memberg = append(memberg, *memberi)
	}

	mam := make(map[string]int)
	for _, meb := range memberg {
		mam[meb.ID] = 0
	}
	for _, log := range logs {
		mam[log.MemberID] += log.Amount
	}
	for _, meb := range memberg {
		if mam[meb.ID] < fee.Amount {
			members = append(members, meb)
		}
	}

	return members, client.Disconnect(ctx)
}

// Enquiry returns the log of club spendings
//
// NOTE:
//
// It is member-limited operation:
//	Only the authenticated members can access to this operation.
func Enquiry(year, semester int) (logs []Log, err error) {
	ctx, cancle := context.WithTimeout(context.Background(), time.Second)
	defer cancle()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MongoURI))
	if err != nil {
		return
	}

	fee := new(Fee)
	log := new(Log)

	filter := func() bson.D {
		arr := make(bson.A, len(fee.Logs))
		for idx, pid := range fee.Logs {
			arr[idx] = pid
		}
		return bson.D{bson.E{Key: "pid", Value: bson.D{bson.E{Key: "$in", Value: arr}}}}
	}()
	logss := make([]Log, len(fee.Logs))
	cur, err := client.Database("club").Collection("logs").Find(ctx, filter)
	if err != nil {
		return
	}
	for cur.Next(ctx) {
		if err = cur.Decode(log); err != nil {
			return
		}
		logss = append(logss, *log)
	}

	for _, log := range logss {
		if log.Type == "direct" {
			logs = append(logs, Log{})
		}
	}

	return logs, client.Disconnect(ctx)

}

// Qussetion1. on func Enquiry do I need to add "approved" as a condition???
// Im asking this because wouldn't that makes clubs members able to see how much other person paied for the fee????
