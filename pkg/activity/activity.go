// Package activity provides access to the club activity of the Buddy System.
package activity

import (
	"context"
	"errors"
	"time"

	"github.com/kmu-kcc/buddy-backend/config"
	"github.com/kmu-kcc/buddy-backend/pkg/member"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ErrAlreadyParticipant = errors.New("already participating")
	ErrAlreadyApplicant   = errors.New("already applied")
	ErrNotMember          = errors.New("no such member")
	ErrBeingProcessed     = errors.New("already being processed")
	ErrNotParticipant     = errors.New("not participating")
)

// Activity represents a club activity state.
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
	// Pictures     []Picture          `json:"pictures" bson:"pictures"`
}

type Activities []Activity

// New returns a new activity.
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

// Actfilter returns limited information of activity.
func (as Activities) Actfilter() (res []map[string]interface{}) {
	for _, activities := range as {
		res = append(res, map[string]interface{}{
			"start":        activities.Start,
			"end":          activities.End,
			"place":        activities.Place,
			"type":         activities.Type,
			"description":  activities.Description,
			"participants": activities.Participants,
		})
	}
	return
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

// Participants returns the participants list of the activity of activityID.
//
// Note:
//
// It is privileged operation:
//	Only the club managers can access to this operation.
func Participants(activityID primitive.ObjectID) (members member.Members, err error) {
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
		arr := make(bson.A, len(activity.Participants))
		for idx, p := range activity.Participants {
			arr[idx] = p
		}
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

	if err = cur.Close(ctx); err != nil {
		return
	}

	return members, client.Disconnect(ctx)
}

// ApplyP applies for an activity of activityID.
//
// NOTE:
//
// It is member-limited operation:
//	Only the authenticated members can access to this operation.
func ApplyP(activityID primitive.ObjectID, memberID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MongoURI))
	if err != nil {
		return err
	}

	collection := client.Database("club").Collection("activities")
	activity := new(Activity)

	if err = collection.FindOne(ctx, bson.M{"_id": activityID}).Decode(activity); err != nil {
		return err
	}

	if func() bool {
		for _, p := range activity.Participants {
			if memberID == p {
				return true
			}
		}
		return false
	}() {
		if err = client.Disconnect(ctx); err != nil {
			return err
		}
		return ErrAlreadyParticipant
	}

	if func() bool {
		for _, a := range activity.Applicants {
			if memberID == a {
				return true
			}
		}
		return false
	}() {
		if err = client.Disconnect(ctx); err != nil {
			return err
		}
		return ErrAlreadyApplicant
	}

	if _, err = collection.UpdateByID(ctx, activity.ID, bson.M{"$push": bson.M{"applicants": memberID}}); err != nil {
		return err
	}
	return client.Disconnect(ctx)
}

// Papplies returns the applicant list of the activity of activityID.
//
// NOTE:
//
// It is privileged operation:
//	Only the club managers can access to this operation.
func Papplies(activityID primitive.ObjectID) (members member.Members, err error) {
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
		Decode(activity); err != nil {
		return
	}

	filter := func() bson.D {
		arr := make(bson.A, len(activity.Applicants))
		for idx, applicant := range activity.Applicants {
			arr[idx] = applicant
		}
		return bson.D{bson.E{Key: "id", Value: bson.D{bson.E{Key: "$in", Value: arr}}}}
	}()

	cur, err := client.Database("club").Collection("members").Find(ctx, filter)
	if err != nil {
		return
	}

	for cur.Next(ctx) {
		if err = cur.Decode(member); err != nil {
			return
		}
		members = append(members, *member)
	}
	if err = cur.Close(ctx); err != nil {
		return
	}
	return members, client.Disconnect(ctx)
}

// ApproveP approve the applicants lis of the activity of activityID.
//
// NOTE:
//
// It is privileged operation:
//	Only the club managers can access to this operation.
func ApproveP(activityID primitive.ObjectID, ids []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MongoURI))
	if err != nil {
		return err
	}

	update := func() bson.D {
		arr := make(bson.A, len(ids))
		for idx, id := range ids {
			arr[idx] = id
		}
		return bson.D{
			bson.E{
				Key: "$pull",
				Value: bson.D{
					bson.E{
						Key: "applicants",
						Value: bson.D{
							bson.E{
								Key:   "$in",
								Value: arr,
							},
						},
					},
				},
			},
			bson.E{
				Key: "$push",
				Value: bson.D{
					bson.E{
						Key: "participants",
						Value: bson.D{
							bson.E{
								Key:   "$each",
								Value: arr},
						},
					},
				},
			},
		}
	}()

	if _, err := client.Database("club").
		Collection("activities").
		UpdateByID(ctx, activityID, update); err != nil {
		return err
	}
	return client.Disconnect(ctx)
}

// RejectP reject the applicanst list of the activity of activityID.
//
// NOTE:
//
// It is privileged operation:
//	Only the club managers can access to this operation.
func RejectP(activityID primitive.ObjectID, ids []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MongoURI))
	if err != nil {
		return err
	}

	update := func() bson.D {
		arr := make(bson.A, len(ids))
		for idx, id := range ids {
			arr[idx] = id
		}
		return bson.D{
			bson.E{
				Key: "$pull",
				Value: bson.D{
					bson.E{
						Key: "applicants",
						Value: bson.D{
							bson.E{
								Key:   "$in",
								Value: arr,
							},
						},
					},
				},
			},
		}
	}()

	if _, err = client.Database("club").
		Collection("activities").
		UpdateByID(ctx, activityID, update); err != nil {
		return err
	}
	return client.Disconnect(ctx)
}

// CancelP cancels the member of memberID's apply request to the activity of activityID.
//
// NOTE:
//
// It is member-limited operation:
//	Only the authenticated members can access to this operation.s
func CancelP(activityID primitive.ObjectID, memberID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MongoURI))
	if err != nil {
		return err
	}

	if _, err = client.Database("club").
		Collection("activities").
		UpdateByID(ctx, activityID, bson.M{"$pull": bson.M{"applicants": memberID}}); err != nil {
		return err
	}
	return client.Disconnect(ctx)
}

// ApplyC applies a cancelation of registration
// for an activity of activityID of a member of memberID.
//
// NOTE:
//
// It is member-limited operation:
//	Only the authenticated members can access to this operation.
func ApplyC(activityID primitive.ObjectID, memberID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MongoURI))
	if err != nil {
		return err
	}

	collection := client.Database("club").Collection("activities")
	targetActivity := new(Activity)

	if err = collection.FindOne(
		ctx,
		bson.D{
			bson.E{Key: "_id", Value: activityID},
		},
	).Decode(targetActivity); err != nil {
		return err
	}

	if !func() bool {
		for _, p := range targetActivity.Participants {
			if memberID == p {
				return true
			}
		}
		return false
	}() {
		if err = client.Disconnect(ctx); err != nil {
			return err
		}
		return ErrNotParticipant
	}

	if func() bool {
		for _, c := range targetActivity.Cancelers {
			if memberID == c {
				return true
			}
		}
		return false
	}() {
		if err = client.Disconnect(ctx); err != nil {
			return err
		}
		return ErrBeingProcessed
	}

	if _, err = collection.UpdateByID(ctx, activityID, bson.D{
		bson.E{Key: "$push", Value: bson.D{
			bson.E{
				Key:   "cancelers",
				Value: memberID,
			}}}}); err != nil {
		return err
	}
	return client.Disconnect(ctx)
}

// CancelC deletes the member of memberID
// from cancelers of the activity of activityID.
//
// NOTE:
//
// It is member-limited operation:
// 	Only the authenticated members can access to this operation.
func CancelC(activityID primitive.ObjectID, memberID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MongoURI))
	if err != nil {
		return err
	}

	collection := client.Database("club").Collection("activities")
	targetActivity := new(Activity)

	if err = collection.FindOne(
		ctx,
		bson.D{
			bson.E{Key: "_id", Value: activityID},
		},
	).Decode(targetActivity); err != nil {
		return err
	}

	if !func() bool {
		for _, p := range targetActivity.Cancelers {
			if memberID == p {
				return true
			}
		}
		return false
	}() {
		if err = client.Disconnect(ctx); err != nil {
			return err
		}
		return ErrNotMember
	}

	if _, err = collection.UpdateByID(ctx, activityID, bson.D{
		bson.E{Key: "$pull", Value: bson.D{
			bson.E{
				Key:   "cancelers",
				Value: memberID,
			}}}}); err != nil {
		return err
	}
	return client.Disconnect(ctx)
}

// Capplies returns participants of an activity of activityID
//
// NOTE:
//
// It is privileged operation:
// 	Only the club managers can access to this operation.
func Capplies(activityID primitive.ObjectID) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MongoURI))
	if err != nil {
		return nil, err
	}

	collection := client.Database("club").Collection("activities")
	targetActivity := new(Activity)

	if err = collection.FindOne(
		ctx,
		bson.D{
			bson.E{Key: "_id", Value: activityID},
		},
	).Decode(targetActivity); err != nil {
		return nil, err
	}

	return targetActivity.Participants, err
}

// ApproveC approves the member of memberID to participate to the activity of activityID.
//
// NOTE:
//
// It is privileged operation:
// 	Only the club managers can access to this operation.
func ApproveC(activityID primitive.ObjectID, memberID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MongoURI))
	if err != nil {
		return err
	}

	collection := client.Database("club").Collection("activities")
	targetActivity := new(Activity)

	if err = collection.FindOne(
		ctx,
		bson.D{
			bson.E{Key: "_id", Value: activityID},
		},
	).Decode(targetActivity); err != nil {
		return err
	}

	if !func() bool {
		for _, p := range targetActivity.Cancelers {
			if memberID == p {
				return true
			}
		}
		return false
	}() {
		if err = client.Disconnect(ctx); err != nil {
			return err
		}
		return ErrNotMember
	}

	if _, err = collection.UpdateByID(ctx, activityID, bson.D{
		bson.E{Key: "$push", Value: bson.D{
			bson.E{
				Key:   "participants",
				Value: memberID,
			}}}}); err != nil {
		return err
	}

	if _, err = collection.UpdateByID(ctx, activityID, bson.D{
		bson.E{Key: "$pull", Value: bson.D{
			bson.E{
				Key:   "cancelers",
				Value: memberID,
			}}}}); err != nil {
		return err
	}

	return client.Disconnect(ctx)
}

// RejectC rejects the member of memberID to participate to the activity of activityID.
//
// NOTE:
//
// It is privileged operation:
// 	Only the club managers can access to this operation.
func RejectC(activityID primitive.ObjectID, memberID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MongoURI))
	if err != nil {
		return err
	}

	collection := client.Database("club").Collection("activities")
	targetActivity := new(Activity)

	if err = collection.FindOne(
		ctx,
		bson.D{
			bson.E{Key: "_id", Value: activityID},
		},
	).Decode(targetActivity); err != nil {
		return err
	}

	if !func() bool {
		for _, p := range targetActivity.Cancelers {
			if memberID == p {
				return true
			}
		}
		return false
	}() {
		if err = client.Disconnect(ctx); err != nil {
			return err
		}
		return ErrNotMember
	}

	if _, err = collection.UpdateByID(ctx, activityID, bson.D{
		bson.E{Key: "$pull", Value: bson.D{
			bson.E{
				Key:   "cancelers",
				Value: memberID,
			}}}}); err != nil {
		return err
	}
	return client.Disconnect(ctx)
}
