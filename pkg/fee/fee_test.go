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

func TestApprove(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MongoURI))
	if err != nil {
		t.Fatal(err)
	}

	collection := client.Database("club").Collection("fee")

	testLog := fee.NewLog("20181681", "unapproved", 0)
	// testLog := new(fee.Log)
	// testLog.MemberID = "20181681"
	// testLog.Type = "unapproved"

	// insert test log
	if _, err := collection.InsertOne(ctx, testLog); err != nil {
		t.Fatal(err)
	}

	if err := fee.Approve([]primitive.ObjectID{testLog.ID}); err != nil {
		t.Fatal(err)
	}

	if err = client.Disconnect(ctx); err != nil {
		t.Fatal(err)
	}
}

func TestReject(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MongoURI))
	collection := client.Database("club").Collection("fee")

	if err != nil {
		t.Fatal(err)
	}

	testLog := fee.NewLog("20181681", "unapproved", 0)
	// testLog := new(fee.Log)
	// testLog.MemberID = "20181681"
	// testLog.Type = "unapproved"

	// insert test log
	if _, err := collection.InsertOne(ctx, testLog); err != nil {
		t.Fatal(err)
	}

	if err := fee.Reject([]primitive.ObjectID{testLog.ID}); err != nil {
		t.Fatal(err)
	}
	if err = client.Disconnect(ctx); err != nil {
		t.Fatal(err)
	}
}

func TestDeposit(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MongoURI))
	collection := client.Database("club").Collection("fee")

	if err != nil {
		t.Fatal(err)
	}

	targetSemester := fee.New(2021, 4, 0, []primitive.ObjectID{})

	if _, err := collection.InsertOne(ctx, targetSemester); err != nil {
		t.Fatal(err)
	}

	if err := fee.Deposit(2021, 4, 100); err != nil {
		t.Fatal(err)
	}

	if err = client.Disconnect(ctx); err != nil {
		t.Fatal(err)
	}
}
