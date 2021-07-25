// Package activity provides access to the club activity of the Buddy System.
package activity

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

// Activity represents a club acitivity state.
type Activity struct {
	ID           primitive.ObjectID `json:"id" bson:"_id"`
	Start        int64              `json:"start,string" bson:"start"`
	End          int64              `json:"end,string" bson:"end"`
	Place        string             `json:"place" bson:"place"`
	Type         string             `json:"type" bson:"type"`
	Description  string             `json:"description" bson:"description"`
	Participants []string           `json:"participants" bson:"participants"`
	Applicants   []string           `json:"applicants" bson:"applicants"`
	Cancelers    []string           `json:"cancelers" bson:"cancelers"`
	Private      bool               `json:"private" bson:"private"`
	// Pictures []Picture `json:"pictures" bson:"pictures"`
}

// New returns a new club activity.
func New(start, end int64, place, typ, description string, participants []string, private bool) *Activity {
	return &Activity{
		ID:           primitive.NewObjectID(),
		Start:        start,
		End:          end,
		Place:        place,
		Type:         typ,
		Description:  description,
		Participants: participants,
		Applicants:   []string{},
		Cancelers:    []string{},
		Private:      private,
	}
}

// Search returns search results filtered by filter.
//
// NOTE:
//
// It is member-limited operation:
//	Only the authenticated members can access to this operation.
func Search(filter map[string]interface{}) (activities []Activity, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MongoURI))
	if err != nil {
		return
	}

	collection := client.Database("club").Collection("activities")
	activity := new(Activity)

	cur, err := collection.Find(ctx, filter)
	if err == mongo.ErrNoDocuments {
		return activities, client.Disconnect(ctx)
	} else if err != nil {
		return
	}

	for cur.Next(ctx) {
		if err = cur.Decode(activity); err != nil {
			return
		}
		activities = append(activities, *activity)
	}
	return activities, client.Disconnect(ctx)
}

// Update updates the contents to update.
//
// NOTE:
//
// It is privileged operation:
//	Only the club managers can access to this operation.
func Update(update map[string]interface{}) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MongoURI))
	if err != nil {
		return
	}

	if _, err = client.Database("club").
		Collection("activities").
		UpdateByID(ctx, update["_id"], bson.M{"$set": update}); err != nil {
		return
	}
	return client.Disconnect(ctx)
}

// Delete deletes a club activity using activityID.
//
// NOTE:
//
// It is privileged operation:
//	Only the club managers can access to this operation.
func Delete(activityID primitive.ObjectID) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MongoURI))
	if err != nil {
		return
	}

	if _, err = client.Database("club").
		Collection("activities").
		DeleteOne(ctx, bson.M{"_id": activityID}); err != nil {
		return
	}

	return client.Disconnect(ctx)
}

// Participants returns the participants list of the activity of activityID.
//
// Note:
//
// It is privileged operation:
//	Only the club managers can access to this operation.
func Participants(activityID primitive.ObjectID) (members []member.Member, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MongoURI))
	if err != nil {
		return
	}

	activity := new(Activity)
	member := new(member.Member)

	if err = client.Database("club").
		Collection("activities").
		FindOne(ctx, bson.M{"_id": activityID}).
		Decode(activity); err == mongo.ErrNoDocuments {
		return
	} else if err != nil {
		return
	}

	filter := func() bson.M {
		arr := func() bson.A {
			arr := make(bson.A, len(activity.Participants))
			for idx, p := range activity.Participants {
				arr[idx] = p
			}
			return arr
		}()
		return bson.M{"id": bson.M{"$in": arr}}
	}()

	// do transaction with members collection
	cur, err := client.Database("club").
		Collection("members").
		Find(ctx, filter)

	if err == mongo.ErrNoDocuments {
		return members, client.Disconnect(ctx)
	} else if err != nil {
		return
	}

	for cur.Next(ctx) {
		if err = cur.Decode(member); err != nil {
			return
		}
		members = append(members, *member)
	}
	return members, client.Disconnect(ctx)
}
