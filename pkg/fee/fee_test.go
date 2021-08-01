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
	if err := fee.Create(2021, 2, 40000); err != nil {
		t.Error(err)
	}
}

func TestSubmit(t *testing.T) {
	if err := fee.Submit("abc", 2021, 2, 20000); err != nil {
		t.Error(err)
	}
}

func TestAmount(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MongoURI))
	if err != nil {
		return
	}

	log1 := new(fee.Log)
	log1.ID = primitive.NewObjectID()
	log1.MemberID = "abc"
	log1.Type = "approved"
	log1.Amount = 20000
	log1.CreatedAt = time.Now().Unix()
	log1.UpdatedAt = time.Now().Unix()

	log2 := new(fee.Log)
	log2.ID = primitive.NewObjectID()
	log2.MemberID = "abc"
	log2.Type = "approved"
	log2.Amount = 20000
	log2.CreatedAt = time.Now().Unix()
	log2.UpdatedAt = time.Now().Unix()

	testFee := fee.New(2021, 1, 40000)
	testFee.Logs = []primitive.ObjectID{log1.ID, log2.ID}

	if _, err := client.Database("club").Collection("logs").InsertMany(ctx, []interface{}{log1, log2}); err != nil {
		t.Error(err)
	}
	if _, err := client.Database("club").Collection("fees").InsertOne(ctx, testFee); err != nil {
		t.Error(err)
	}

	sum, err := fee.Amount(2021, 1, "abc")
	if err != nil {
		t.Error(err)
	} else {
		t.Log(sum)
	}
}

func TestDones(t *testing.T) {
	if members, err := fee.Dones(2021, 1); err != nil {
		t.Error(err)
	} else {
		t.Log(members)
	}
}

func TestYets(t *testing.T) {
	if members, err := fee.Yets(2021, 1); err != nil {
		t.Error(err)
	} else {
		t.Log(members)
	}
}

func TestAll(t *testing.T) {
	if logs, err := fee.All(2021, 1); err != nil {
		t.Error(err)
	} else {
		t.Log(logs)
	}
}
