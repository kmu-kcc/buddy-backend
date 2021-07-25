package activity_test

import (
	"context"
	"testing"
	"time"

	"github.com/kmu-kcc/buddy-backend/pkg/activity"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestNew(t *testing.T) {
	ns := make([]string, 10, 20)
	a := activity.New(1, 2, "", "", "", ns, true)

	ctx, cancle := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancle()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		t.Error(err)
	}

	collection := client.Database("club").Collection("activities")

	if res, err := collection.InsertOne(
		ctx,
		a,
	); err != nil {
		t.Error(err)
	} else {
		t.Log(res)
	}

	if err = client.Disconnect(ctx); err != nil {
		t.Error(err)
	}
}

func TestApplyP(t *testing.T) {
	res, err := primitive.ObjectIDFromHex("60fd5ad1e26bd52bc5b0bf47")
	if err != nil {
		t.Error(err)
	}
	if err = activity.ApplyP(res, "20172228"); err != nil {
		t.Error(err)
	}
}

func TestPapplies(t *testing.T) {
	res, err := primitive.ObjectIDFromHex("60fd54362b5226020a8c945b")
	if err != nil {
		t.Error(err)
	}

	if li, err := activity.Papplies(res); err != nil {
		t.Error(err)
	} else {
		t.Log(li)
	}
}

func TestApproveP(t *testing.T) {
	res, err := primitive.ObjectIDFromHex("60fd54362b5226020a8c945b")
	if err != nil {
		t.Error(err)
	}
	if err = activity.ApproveP(res, []string{"20172229", "20172228"}); err != nil {
		t.Error(err)
	}
}

func TestRejectP(t *testing.T) {
	res, err := primitive.ObjectIDFromHex("60fd5ad1e26bd52bc5b0bf47")
	if err != nil {
		t.Error(err)
	}

	if err = activity.RejectP(res, []string{"20172229", "20172228"}); err != nil {
		t.Error(err)
	}
}

func TestCancelP(t *testing.T) {
	res, err := primitive.ObjectIDFromHex("60fd54362b5226020a8c945b")
	if err != nil {
		t.Error(err)
	}

	if err = activity.CancelP(res, "20172229"); err != nil {
		t.Error(err)
	}
}
