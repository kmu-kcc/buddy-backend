package activity_test

import (
	"context"
	"testing"
	"time"

	"github.com/kmu-kcc/buddy-backend/config"
	"github.com/kmu-kcc/buddy-backend/pkg/activity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestInsert(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MongoURI))
	if err != nil {
		t.Error(err)
	}

	collection := client.Database("club").Collection("activities")

	if res, err := collection.InsertMany(ctx, []interface{}{
		bson.D{bson.E{Key: "type", Value: "MT"}},
		bson.D{bson.E{Key: "type", Value: "meet"}, {Key: "place", Value: "cafe"}},
		bson.D{bson.E{Key: "place", Value: "home"}},
		bson.D{bson.E{Key: "type", Value: "study"}, bson.E{Key: "description", Value: "ok"}},
	}); err != nil {
		t.Error(err)
	} else {
		t.Log(res)
	}

	if err = client.Disconnect(ctx); err != nil {
		t.Error(err)
	}
}
func TestSearch(t *testing.T) {
	filters := []map[string]interface{}{
		{
			"type": "MT",
		},
		{
			"place": "cafe",
			"type":  "meet",
		},
		{
			"place": "home",
		},
		{
			"type":        "study",
			"description": "ok",
		},
	}

	for _, filter := range filters {
		if activities, err := activity.Search(filter); err != nil {
			t.Error(err)
		} else {
			t.Log(activities)
		}
	}
}

func TestUpdate(t *testing.T) {
	objectId, err := primitive.ObjectIDFromHex("60fcac8824c06103861b13f2")
	if err != nil {
		t.Error(err)
	}
	filters := []map[string]interface{}{
		{
			"_id":  objectId,
			"type": "meet",
		},
	}

	for _, filter := range filters {
		if err := activity.Update(filter); err != nil {
			t.Error(err)
		}
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
