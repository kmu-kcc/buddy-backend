// Package activity provides access to the club activity of the Buddy System.
package activity

import (
	"context"
	"time"

	"github.com/kmu-kcc/buddy-backend/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	foundingEvent = iota
	study
	etc
)

// Activity represents a club activity state.
type Activity struct {
	ID           primitive.ObjectID `json:"id" bson:"_id"`
	Start        int64              `json:"start,string" bson:"start"`
	End          int64              `json:"end,string" bson:"end"`
	Place        string             `json:"place" bson:"place"`
	Type         int                `json:"type" bson:"type"`
	Description  string             `json:"description" bson:"description"`
	Participants []string           `json:"participants" bson:"participants"`
	Private      bool               `json:"private" bson:"private"`
	// Pictures     []Picture          `json:"pictures" bson:"pictures"`
}

type Activities []Activity

// New returns a new activity.
func New(start, end int64, place, description string, typ int, participants []string, private bool) *Activity {
	return &Activity{
		ID:           primitive.NewObjectID(),
		Start:        start,
		End:          end,
		Place:        place,
		Type:         typ,
		Description:  description,
		Participants: participants,
		Private:      private,
	}
}

// Public returns the limited informations of as.
func (as Activities) Public() []map[string]interface{} {
	pubs := make([]map[string]interface{}, len(as))

	for idx, a := range as {
		pubs[idx] = make(map[string]interface{})
		pubs[idx]["start"] = a.Start
		pubs[idx]["end"] = a.End
		pubs[idx]["place"] = a.Place
		pubs[idx]["type"] = a.Type
		pubs[idx]["description"] = a.Description
		pubs[idx]["participants"] = a.Participants
		// pubs[idx]["pictures"] = a.Pictures
	}
	return pubs
}

// Create creates a new activity.
//
// NOTE:
//
// It is privileged operation:
//	Only the club managers can access to this operation.
func (a Activity) Create() (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MongoURI))
	if err != nil {
		return
	}

	if _, err = client.Database("club").
		Collection("activities").
		InsertOne(ctx, a); err != nil {
		return
	}
	return client.Disconnect(ctx)
}

// Search returns search results filtered by filter.
//
// NOTE:
//
// It is privileged operation:
//	Only the club managers can access to this operation.
func Search(query string) (activities Activities, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MongoURI))
	if err != nil {
		return
	}

	collection := client.Database("club").Collection("activities")
	activity := new(Activity)

	cur, err := collection.Find(ctx,
		bson.M{"$or": []bson.M{
			// {"start": bson.M{"$regex": query}},
			// {"end": bson.M{"$regex": query}},
			{"place": bson.M{"$regex": query}},
			{"type": bson.M{"$regex": query}},
			{"description": bson.M{"$regex": query}},
		}})

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

	if err = cur.Close(ctx); err != nil {
		return
	}

	return activities, client.Disconnect(ctx)
}

// Update updates the contents to update.
//
// NOTE:
//
// It is privileged operation:
//	Only the club managers can access to this operation.
func (a Activity) Update(update map[string]interface{}) (err error) {
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
