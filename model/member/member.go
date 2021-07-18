// Package member provides access to the club member of the Buddy System.
package member

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Member represents a club member state.
type Member struct {
	ID         string `json:"id"`         // student ID
	password   string `json:"-"`          // password - soon deprecated
	Name       string `json:"name"`       // Name
	Department string `json:"department"` // department - magic number needed
	Grade      string `json:"grade"`      // grade
	Phone      string `json:"phone"`      // phone number
	Email      string `json:"email"`      // e-mail address
	Enrollment string `json:"enrollment"` // enrollment state (attending/absent/graduated) - magic number needed
	Verified   bool   `json:"verified"`   // e-mail verified or not
	Approved   bool   `json:"approved"`   // approved or not
	Privileged bool   `json:"privileged"` // is administrator or not - soon deprecated
	OnDelete   bool   `json:"on_delete"`  // on exit process or not
	CreatedAt  int64  `json:"created_at"` // when created - Unix timestamp
	UpdatedAt  int64  `json:"updated_at"` // last updated - Unix timestamp
}

// New returns a new club member.
func New(id string, name string, department string, grade string, phone string, email string, enrollment string) *Member {
	return &Member{
		ID:         id,
		password:   id,
		Name:       name,
		Department: department,
		Grade:      grade,
		Phone:      phone,
		Email:      email,
		Enrollment: enrollment,
		Verified:   false,
		Approved:   false,
		Privileged: false,
		OnDelete:   false,
		CreatedAt:  time.Now().Unix(),
		UpdatedAt:  time.Now().Unix(),
	}
}

// toBSON converts m to an ordered BSON document.
func (m Member) toBSON() bson.D {
	return bson.D{
		bson.E{Key: "id", Value: m.ID},
		bson.E{Key: "password", Value: m.password},
		bson.E{Key: "name", Value: m.Name},
		bson.E{Key: "department", Value: m.Department},
		bson.E{Key: "grade", Value: m.Grade},
		bson.E{Key: "phone", Value: m.Phone},
		bson.E{Key: "email", Value: m.Email},
		bson.E{Key: "enrollment", Value: m.Enrollment},
		bson.E{Key: "verified", Value: m.Verified},
		bson.E{Key: "approved", Value: m.Approved},
		bson.E{Key: "privileged", Value: m.Privileged},
		bson.E{Key: "on_delete", Value: m.OnDelete},
		bson.E{Key: "created_at", Value: m.CreatedAt},
		bson.E{Key: "updated_at", Value: m.UpdatedAt},
	}
}

// SignUp applies a membership of m.
// If m already exists (approved or not), nothing changes.
// Else it registers an unapproved member.
func (m Member) SignUp() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// TODO
	// update URI
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		return err
	}

	collection := client.Database("club").Collection("members")
	member := new(Member)

	if err = collection.FindOne(ctx, bson.D{bson.E{Key: "id", Value: m.ID}}).Decode(member); err == mongo.ErrNoDocuments {
		if _, err = collection.InsertOne(ctx, m.toBSON()); err != nil {
			return err
		}
		return client.Disconnect(ctx)
	} else if err != nil {
		return err
	} else if err = client.Disconnect(ctx); err != nil {
		return err
	} else if member.Approved {
		return errors.New("already a member")
	} else {
		return errors.New("under review")
	}
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

	// TODO
	// update URI
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		return err
	}

	collection := client.Database("club").Collection("members")

	for _, id := range ids {
		if _, err = collection.UpdateOne(ctx, bson.D{bson.E{Key: "id", Value: id}}, bson.D{bson.E{Key: "approved", Value: true}}); err != nil {
			return err
		}
	}

	return client.Disconnect(ctx)
}

// Reject deletes the signup requests of ids.
//
// NOTE:
//
// It is a privileged operation:
//	Only the club managers can access to this operation.
func Reject(ids []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// TODO
	// update URI
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		return err
	}

	collection := client.Database("club").Collection("members")

	for _, id := range ids {
		if _, err = collection.DeleteOne(ctx, bson.D{bson.E{Key: "id", Value: id}}); err != nil {
			return err
		}
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

	// TODO
	// update URI
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
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

	if member.OnDelete {
		if err = client.Disconnect(ctx); err != nil {
			return err
		}
		return errors.New("already on delete")
	} else {
		if err = client.Disconnect(ctx); err != nil {
			return err
		}
		return errors.New("exit request success")
	}
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

	// TODO
	// update URI
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		return err
	}

	collection := client.Database("club").Collection("members")

	for _, id := range ids {
		if _, err = collection.DeleteOne(ctx, bson.D{bson.E{Key: "id", Value: id}}); err != nil {
			return err
		}
	}

	return client.Disconnect(ctx)
}

// Members returns the all member state.
//
// NOTE:
//
// It is a member-limited operation:
//	Only the authenticated members can access to this operation.
func Members() (members []Member, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// TODO
	// update URI
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		return
	}

	member := new(Member)

	cur, err := client.Database("club").
		Collection("members").
		Find(ctx, bson.D{})
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

// Verify sets m to be verified.
func (m Member) Verify() error {
	// TODO
	// e-mail verification (OAuth)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// TODO
	// update URI
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		return err
	}

	if _, err = client.Database("club").
		Collection("member").
		UpdateOne(ctx,
			bson.D{bson.E{Key: "id", Value: m.ID}},
			bson.D{bson.E{Key: "verified", Value: true}}); err != nil {
		return err
	}

	return client.Disconnect(ctx)
}
