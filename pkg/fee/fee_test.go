package fee_test

import (
	"context"
	"testing"
	"time"

	"github.com/kmu-kcc/buddy-backend/config"
	"github.com/kmu-kcc/buddy-backend/pkg/fee"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestCreate(t *testing.T) {
	fees := []*fee.Fee{
		fee.New(2022, 1, 0, 30000),
		fee.New(2022, 2, 0, 40000),
	}

	for _, fee := range fees {
		if err := fee.Create(); err != nil {
			t.Error(err)
		}
	}
}

func TestAmount(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MongoURI))
	if err != nil {
		return
	}

	log1 := fee.NewLog("abc", 20000, 0)
	log2 := fee.NewLog("abc", 30000, 0)

	testFee := new(fee.Fee)
	testFee.Year = 2022
	testFee.Semester = 2
	testFee.Amount = 30000
	testFee.Logs = []primitive.ObjectID{log1.ID, log2.ID}

	if _, err := client.Database("club").Collection("logs").InsertMany(ctx, []interface{}{log1, log2}); err != nil {
		t.Error(err)
	}
	if _, err := client.Database("club").Collection("fees").InsertOne(ctx, testFee); err != nil {
		t.Error(err)
	}

	sum, err := fee.Amount(2022, 2, "abc")
	if err != nil {
		t.Error(err)
	} else {
		t.Log(sum)
	}
}

func TestDones(t *testing.T) {
	f := fee.Fee{Year: 2021, Semester: 1}
	if members, err := f.Dones(); err != nil {
		t.Error(err)
	} else {
		t.Log(members)
	}
}

func TestYets(t *testing.T) {
	f := fee.Fee{Year: 2021, Semester: 1}
	if members, err := f.Yets(); err != nil {
		t.Error(err)
	} else {
		t.Log(members)
	}
}

// func TestAll(t *testing.T) {
// 	if logs, err := fee.All(1627794157, 1627794157); err != nil {
// 		t.Error(err)
// 	} else {
// 		t.Log(logs)
// 	}
// }

// func TestApprove(t *testing.T) {
// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
// 	defer cancel()

// 	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MongoURI))
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	collection := client.Database("club").Collection("fees")
// 	collectionLogs := client.Database("club").Collection("logs")

// 	testLog := fee.NewLog("20181681", "unapproved", 0)
// 	testFee := fee.New(2021, 4, 0)

// 	testFee.Logs = append(testFee.Logs, testLog.ID)

// 	// insert test log
// 	if _, err := collection.InsertOne(ctx, testFee); err != nil {
// 		t.Fatal()
// 	}
// 	if _, err := collectionLogs.InsertOne(ctx, testLog); err != nil {
// 		t.Fatal(err)
// 	}

// 	if err := fee.Approve([]primitive.ObjectID{testLog.ID}); err != nil {
// 		t.Fatal(err)
// 	}

// 	if err = client.Disconnect(ctx); err != nil {
// 		t.Fatal(err)
// 	}
// }

// func TestDeposit(t *testing.T) {
// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
// 	defer cancel()

// 	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MongoURI))
// 	collection := client.Database("club").Collection("fees")

// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	targetSemester := fee.New(2021, 3, 0)

// 	if _, err := collection.InsertOne(ctx, targetSemester); err != nil {
// 		t.Fatal(err)
// 	}

// 	if err := fee.Deposit(2021, 3, 100); err != nil {
// 		t.Fatal(err)
// 	}

// 	if err = client.Disconnect(ctx); err != nil {
// 		t.Fatal(err)
// 	}
// }
