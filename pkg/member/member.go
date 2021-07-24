// Package member provides access to the club member of the Buddy System.
package member

import (
	"context"
	"errors"
	"time"

	"github.com/kmu-kcc/buddy-backend/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	Attending = iota
	Absent
	Graduate
)

var (
	ErrWrongPassword = errors.New("wrong password")
	ErrAlreadyMember = errors.New("already a member")
	ErrUnderReview   = errors.New("under review")
	ErrOnDelete      = errors.New("already on delete")
	ErrNotOnDelete   = errors.New("not on delete")
)

// Member represents a club member state.
type Member struct {
	ID         string `json:"id" bson:"id"`                 // student ID
	Password   string `json:"password" bson:"password"`     // password
	Name       string `json:"name" bson:"name"`             // Name
	Department string `json:"department" bson:"department"` // department
	Grade      string `json:"grade" bson:"grade"`           // grade
	Phone      string `json:"phone" bson:"phone"`           // phone number
	Email      string `json:"email" bson:"email"`           // e-mail address
	Attendance int    `json:"attendance" bson:"attendance"` // attendance status (attending/absent/graduate)
	Approved   bool   `json:"approved" bson:"approved"`     // approved or not
	OnDelete   bool   `json:"on_delete" bson:"on_delete"`   // on exit process or not
	CreatedAt  int64  `json:"created_at" bson:"created_at"` // when created - Unix timestamp
	UpdatedAt  int64  `json:"updated_at" bson:"updated_at"` // last updated - Unix timestamp
}

// New returns a new club member.
func New(id string, name string, department string, grade string, phone string, email string, attendance int) *Member {
	now := time.Now().Unix()
	return &Member{
		ID:         id,
		Password:   id,
		Name:       name,
		Department: department,
		Grade:      grade,
		Phone:      phone,
		Email:      email,
		Attendance: attendance,
		Approved:   false,
		OnDelete:   false,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}

// SingIn checks whether m is a club member.
//
// NOTE:
//
// It is a member-limited operation:
//	Only the authenticated members can access to this operation.
func (m Member) SingIn() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MongoURI))
	if err != nil {
		return err
	}

	member := new(Member)

	if err = client.Database("club").
		Collection("members").
		FindOne(
			ctx,
			bson.D{bson.E{Key: "id", Value: m.ID}}).
		Decode(member); err == mongo.ErrNoDocuments {
		if err = client.Disconnect(ctx); err != nil {
			return err
		}
		return mongo.ErrNoDocuments
	} else if err != nil {
		return err
	}
	if err = client.Disconnect(ctx); err != nil {
		return err
	}
	if !member.Approved {
		return ErrUnderReview
	}
	if m.Password != member.Password {
		return ErrWrongPassword
	}
	return nil
}

// SignUp applies a membership of m.
// If m already exists (approved or not), nothing changes.
// Else it registers an unapproved member.
func (m Member) SignUp() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MongoURI))
	if err != nil {
		return err
	}

	collection := client.Database("club").Collection("members")
	member := new(Member)

	if err = collection.FindOne(
		ctx,
		bson.D{bson.E{Key: "id", Value: m.ID}}).
		Decode(member); err == mongo.ErrNoDocuments {
		if _, err = collection.InsertOne(ctx, m); err != nil {
			return err
		}
		return client.Disconnect(ctx)
	} else if err != nil {
		return err
	}
	if err = client.Disconnect(ctx); err != nil {
		return err
	}
	if member.Approved {
		return ErrAlreadyMember
	}
	return ErrUnderReview
}

// SignUps returns the signup request list.
//
// NOTE:
//
// It is a privileged operation:
//	Only the club managers can access to this operation.
func SignUps() (members []Member, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MongoURI))
	if err != nil {
		return
	}

	cur, err := client.Database("club").
		Collection("members").
		Find(
			ctx,
			bson.D{bson.E{Key: "approved", Value: false}})

	if err != nil {
		return
	}

	member := new(Member)

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

// Approve approves the signup requests of ids.
//
// NOTE:
//
// It is a privileged operation:
//	Only the club managers can access to this operation.
func Approve(ids []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MongoURI))
	if err != nil {
		return err
	}

	filter := func() bson.D {
		arr := make(bson.A, len(ids))
		for idx, id := range ids {
			arr[idx] = id
		}
		return bson.D{bson.E{Key: "id", Value: bson.D{bson.E{Key: "$in", Value: arr}}}}
	}()

	if _, err = client.Database("club").
		Collection("members").
		UpdateMany(
			ctx,
			filter,
			bson.D{bson.E{Key: "approved", Value: true}}); err != nil {
		return err
	}

	return client.Disconnect(ctx)
}

// Delete deletes the members of ids.
//
// NOTE:
//
// It is a privileged operation:
//	Only the club managers can access to this operation.
func Delete(ids []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MongoURI))
	if err != nil {
		return err
	}

	filter := func() bson.D {
		arr := make(bson.A, len(ids))
		for idx, id := range ids {
			arr[idx] = id
		}
		return bson.D{bson.E{Key: "id", Value: bson.D{bson.E{Key: "$in", Value: arr}}}}
	}()

	if _, err = client.Database("club").
		Collection("members").
		DeleteMany(ctx, filter); err != nil {
		return err
	}

	return client.Disconnect(ctx)
}

// Exit applies an exit of m.
//
// NOTE:
//
// It is a member-limited operation:
//	Only the authenticated members can access to this operation.
func (m Member) Exit() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MongoURI))
	if err != nil {
		return err
	}

	member := new(Member)

	if err = client.Database("club").
		Collection("members").
		FindOneAndUpdate(
			ctx,
			bson.D{bson.E{Key: "id", Value: m.ID}},
			bson.D{bson.E{Key: "on_delete", Value: true}}).
		Decode(member); err != nil {
		return err
	}

	if err = client.Disconnect(ctx); err != nil {
		return err
	}
	if member.OnDelete {
		return ErrOnDelete
	}
	return nil
}

// Exits returns the exit request list.
//
// NOTE:
//
// It is a privileged operation:
//	Only the club managers can access to this operation.
func Exits() (members []Member, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MongoURI))
	if err != nil {
		return
	}

	cur, err := client.Database("club").
		Collection("members").
		Find(
			ctx,
			bson.D{bson.E{Key: "on_delete", Value: true}})

	if err != nil {
		return
	}

	member := new(Member)

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

// CancelExit cancels the exit request of m.
//
// NOTE:
//
// It is a member-limited operation:
//	Only the authenticated members can access to this operation.
func (m Member) CancelExit() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MongoURI))
	if err != nil {
		return err
	}

	member := new(Member)

	if err = client.Database("club").
		Collection("members").
		FindOneAndUpdate(
			ctx,
			bson.D{bson.E{Key: "id", Value: m.ID}},
			bson.D{bson.E{Key: "on_delete", Value: false}}).
		Decode(member); err != nil {
		return err
	}

	if err = client.Disconnect(ctx); err != nil {
		return err
	}
	if !member.OnDelete {
		return ErrNotOnDelete
	}
	return nil
}

// Members returns the all club member state.
//
// NOTE:
//
// It is a member-limited operation:
//	Only the authenticated members can access to this operation.
func Members() (members []Member, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MongoURI))
	if err != nil {
		return
	}

	cur, err := client.Database("club").
		Collection("members").
		Find(
			ctx,
			bson.D{bson.E{Key: "approved", Value: true}})

	if err != nil {
		return
	}

	member := new(Member)

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
