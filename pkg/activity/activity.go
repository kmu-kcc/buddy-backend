// Package activity provides access to the club activity of the Buddy System.
package activity

import (
	"context"
	"strings"

	"github.com/kmu-kcc/buddy-backend/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	FoundingEvent = iota
	Study
	Etc
)

// Activity represents a club activity state.
type Activity struct {
	ID           primitive.ObjectID `json:"id" bson:"_id"`
	Title        string             `json:"title" bson:"title"`
	Start        int64              `json:"start,string" bson:"start"`
	End          int64              `json:"end,string" bson:"end"`
	Place        string             `json:"place" bson:"place"`
	Type         int                `json:"type" bson:"type"`
	Description  string             `json:"description" bson:"description"`
	Participants []string           `json:"participants" bson:"participants"`
	Private      bool               `json:"private" bson:"private"`
	Files        Files              `json:"files" bson:"files"`
}

type Activities []Activity

// New returns a new activity.
func New(title string, start, end int64, place, description string, typ int, participants []string, private bool) *Activity {
	return &Activity{
		ID:           primitive.NewObjectID(),
		Title:        title,
		Start:        start,
		End:          end,
		Place:        place,
		Type:         typ,
		Description:  description,
		Participants: participants,
		Private:      private,
		Files:        Files{},
	}
}

// Public returns the limited informations of a.
func (a Activity) Public() map[string]interface{} {
	pub := make(map[string]interface{})

	pub["title"] = a.Title
	pub["start"] = a.Start
	pub["end"] = a.End
	pub["place"] = a.Place
	pub["type"] = a.Type
	pub["description"] = a.Description
	pub["participants"] = a.Participants
	pub["files"] = a.Files

	return pub
}

// Public returns the limited informations of as.
func (as Activities) Public() []map[string]interface{} {
	pubs := make([]map[string]interface{}, len(as))

	for idx, activity := range as {
		pubs[idx] = activity.Public()
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
	ctx := context.Background()
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
func Search(query string, private bool) (activities Activities, err error) {
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MongoURI))
	if err != nil {
		return
	}

	var filter bson.D
	activity := new(Activity)

	switch strings.TrimSpace(query) {
	case "창립제":
		filter = bson.D{bson.E{Key: "type", Value: FoundingEvent}}
	case "스터디":
		fallthrough
	case "study":
		filter = bson.D{bson.E{Key: "type", Value: Study}}
	default:
		filter = bson.D{
			bson.E{Key: "$or", Value: bson.A{
				bson.D{
					bson.E{Key: "title", Value: bson.D{
						bson.E{Key: "$regex", Value: query}}}},
				bson.D{
					bson.E{Key: "place", Value: bson.D{
						bson.E{Key: "$regex", Value: query}}}},
				bson.D{
					bson.E{Key: "description", Value: bson.D{
						bson.E{Key: "$regex", Value: query}}}}}}}
	}

	cur, err := client.Database("club").Collection("activities").Find(ctx, filter)

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
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MongoURI))
	if err != nil {
		return
	}

	if _, err = client.Database("club").
		Collection("activities").
		UpdateByID(ctx, a.ID, bson.M{"$set": update}); err != nil {
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
func Delete(id primitive.ObjectID) (err error) {
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MongoURI))
	if err != nil {
		return
	}

	if _, err = client.Database("club").
		Collection("activities").
		DeleteOne(ctx, bson.M{"_id": id}); err != nil {
		return
	}
	return client.Disconnect(ctx)
}

// Upload saves file of filename into a.
func (a Activity) Upload(filename string) error {
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MongoURI))
	if err != nil {
		return err
	}

	if _, err = client.Database("club").
		Collection("activities").
		UpdateByID(
			ctx,
			a.ID,
			bson.D{
				bson.E{Key: "$push", Value: bson.D{
					bson.E{Key: "files", Value: NewFile(filename)}}}}); err != nil {
		return err
	}
	return client.Disconnect(ctx)
}

// DeleteFile deletes file of filename from a.
func (a Activity) DeleteFile(filename string) error {
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MongoURI))
	if err != nil {
		return err
	}

	if _, err = client.Database("club").
		Collection("activities").
		UpdateByID(
			ctx,
			a.ID,
			bson.D{
				bson.E{Key: "$pull", Value: bson.D{
					bson.E{Key: "files", Value: bson.D{
						bson.E{Key: "$in", Value: bson.A{filename}}}}}}}); err != nil {
		return err
	}

	if err := NewFile(filename).Delete(); err != nil {
		return err
	}
	return client.Disconnect(ctx)
}
