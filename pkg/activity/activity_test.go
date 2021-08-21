package activity_test

import (
	"testing"

	"github.com/kmu-kcc/buddy-backend/pkg/activity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestCreate(t *testing.T) {
	acts := []*activity.Activity{
		activity.New(1, 1, "cafe", "study", 0, []string{}, true),
		activity.New(2, 2, "school", "founding festival", 1, []string{}, true),
	}

	for _, act := range acts {
		if err := act.Create(); err != nil {
			t.Error(err)
		}
	}
}

func TestSearch(t *testing.T) {
	if activities, err := activity.Search("te", false); err != nil {
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
