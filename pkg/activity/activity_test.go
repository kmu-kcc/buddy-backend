package activity_test

import (
	"context"
	"testing"
	"time"

	"github.com/kmu-kcc/buddy-backend/config"
	"github.com/kmu-kcc/buddy-backend/pkg/activity"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestCreate(t *testing.T) {
	acts := []*activity.Activity{
		activity.New(1, 1, "cafe", "study", "good", []string{}, true),
		activity.New(2, 2, "school", "founding festival", "wow", []string{}, true),
	}

	for _, act := range acts {
		if err := act.Create(); err != nil {
			t.Error(err)
		}
	}
}

func TestSearch(t *testing.T) {
	if activities, err := activity.Search("go"); err != nil {
		t.Error(err)
	} else {
		t.Log(activities)
	}
}

func TestUpdate(t *testing.T) {
	objectId, err := primitive.ObjectIDFromHex("6113ed60c7913f56af94f532")
	if err != nil {
		t.Error(err)
	}

	act := activity.Activity{ID: objectId}
	if err := act.Update(map[string]interface{}{
		"_id":  objectId,
		"type": "meet"}); err != nil {
		t.Error(err)
	}
}

func TestDelete(t *testing.T) {
	objectId, err := primitive.ObjectIDFromHex("60fcac8824c06103861b13f2")
	if err != nil {
		t.Error(err)
	}

	if err := activity.Delete(objectId); err != nil {
		t.Error(err)
	}
}

func TestParticipants(t *testing.T) {
	objectId, err := primitive.ObjectIDFromHex("60fce3b6b7a36438a91f807d")
	if err != nil {
		t.Error(err)
	}

	if members, err := activity.Participants(objectId); err != nil {
		t.Error(err)
	} else {
		t.Log(members)
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

func TestApplyC(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MongoURI))
	if err != nil {
		t.Error(err)
	}

	testCollection := client.Database("club").Collection("activities")

	// testcase 1. whether applyC works -> Done
	targetActivity := activity.New(1, 1, "Place", "MT", "Testing", []string{"ApplyC"}, false)

	// testcase 2. ErrNotParticipant -> Done
	// testcase 3. ErrBeingProcessed -> Done

	// Initialize Collection
	if err := testCollection.Drop(ctx); err != nil {
		t.Error(err)
	}

	// insert test activity
	if _, err := testCollection.InsertOne(ctx, targetActivity); err != nil {
		t.Log("INSERT_ERR")
		t.Error(err)
	}

	if err = activity.ApplyC(targetActivity.ID, "ApplyC"); err != nil {
		t.Log("APPLYC_ERR")
		t.Error(err)
	}

	// Reset Collection
	if err := testCollection.Drop(ctx); err != nil {
		t.Error(err)
	}

	if err := client.Disconnect(ctx); err != nil {
		t.Error(err)
	}
}

func TestCancelC(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MongoURI))
	if err != nil {
		t.Error(err)
	}

	testCollection := client.Database("club").Collection("activities")
	targetActivity := activity.New(1, 1, "Place", "MT", "Testing", []string{"ApplyC"}, false)

	// insert test activity
	if _, err := testCollection.InsertOne(ctx, targetActivity); err != nil {
		t.Log("INSERT_ERR")
		t.Error(err)
	}

	if err = activity.ApplyC(targetActivity.ID, "ApplyC"); err != nil {
		t.Log("APPLYC_ERR")
		t.Log(targetActivity.ID)
		t.Error(err)
	}

	// testcase 1. whether CancelC works
	// testcase 2. ErrNotMember
	if err = activity.CancelC(targetActivity.ID, "ApplyC"); err != nil {
		t.Log("CANCELC_ERR")
		t.Error(err)
	}

	// Reset Collection
	if err := testCollection.Drop(ctx); err != nil {
		t.Error(err)
	}

	if err := client.Disconnect(ctx); err != nil {
		t.Error(err)
	}
}

func TestCapplies(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MongoURI))
	if err != nil {
		t.Error(err)
	}

	testCollection := client.Database("club").Collection("activities")
	targetActivity := activity.New(1, 1, "Place", "MT", "Testing", []string{"Succeed!", "Yeahhhhh"}, false)

	// insert test activity
	if _, err := testCollection.InsertOne(ctx, targetActivity); err != nil {
		t.Log("INSERT_ERR")
		t.Error(err)
	}

	if _, err := activity.Capplies(targetActivity.ID); err != nil {
		t.Error(err)
	}

	// Reset Collection
	if err := testCollection.Drop(ctx); err != nil {
		t.Error(err)
	}

	if err := client.Disconnect(ctx); err != nil {
		t.Error(err)
	}
}

func TestApproveC(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MongoURI))
	if err != nil {
		t.Error(err)
	}

	testCollection := client.Database("club").Collection("activities")
	targetActivity := activity.New(1, 1, "Place", "MT", "Testing", []string{"Succeed!"}, false)

	// insert test activity
	if _, err := testCollection.InsertOne(ctx, targetActivity); err != nil {
		t.Log("INSERT_ERR")
		t.Error(err)
	}

	// set cancelers
	if err = activity.ApplyC(targetActivity.ID, "Succeed!"); err != nil {
		t.Log("APPLYC_ERR")
		t.Error(err)
	}

	if err = activity.ApproveC(targetActivity.ID, "Succeed!"); err != nil {
		t.Error(err)
	}

	// Reset Collection
	if err := testCollection.Drop(ctx); err != nil {
		t.Error(err)
	}

	if err := client.Disconnect(ctx); err != nil {
		t.Error(err)
	}
}

func TestRejectC(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MongoURI))
	if err != nil {
		t.Error(err)
	}

	testCollection := client.Database("club").Collection("activities")
	targetActivity := activity.New(1, 1, "Place", "MT", "Testing", []string{"Succeed!"}, false)

	// insert test activity
	if _, err := testCollection.InsertOne(ctx, targetActivity); err != nil {
		t.Log("INSERT_ERR")
		t.Error(err)
	}

	// set cancelers
	if err = activity.ApplyC(targetActivity.ID, "Succeed!"); err != nil {
		t.Log("APPLYC_ERR")
		t.Error(err)
	}

	if err = activity.RejectC(targetActivity.ID, "Succeed!"); err != nil {
		t.Error(err)
	}

	// Reset Collection
	if err := testCollection.Drop(ctx); err != nil {
		t.Error(err)
	}

	if err := client.Disconnect(ctx); err != nil {
		t.Error(err)
	}
}
