// Package member provides the CRUD operations for the club member of the Buddy System.
package member

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Member represents a club member state.
type Member struct {
	ID         string `json:"id"`           // student ID
	password   string `json:"-"`            // password - soon deprecated
	Name       string `json:"name"`         // Name
	Department string `json:"department"`   // department - magic number needed
	Grade      uint   `json:"grade,string"` // grade
	Phone      string `json:"phone"`        // phone number
	Email      string `json:"email"`        // e-mail address
	Enrollment string `json:"enrollment"`   // enrollment state (attending/absent/graduated) - magic number needed
	Verified   bool   `json:"verified"`     // e-mail verified or not
	Approved   bool   `json:"approved"`     // approved or not
	Privileged bool   `json:"privileged"`   // is administrator or not - soon deprecated
	OnDelete   bool   `json:"on_delete"`    // on withdrawal process or not
}

// New returns a new club member.
func New(id string, name string, department string, grade uint, phone string, email string, enrollment string) *Member {
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
	}
}

// ToBSON converts m to an ordered BSON document.
func (m *Member) ToBSON() bson.D {
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
	}
}

// SignUp applies a membership of m.
// If m already exists (approved or not), nothing changes.
// Else it registers an unapproved member.
func (m *Member) SignUp() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// TODO
	// update URI
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		return err
	}

	members := client.Database("club").Collection("members")
	if _, err = members.Find(ctx, bson.D{bson.E{Key: "id", Value: m.ID}}); err == mongo.ErrNoDocuments {
		if _, err = members.InsertOne(ctx, m.ToBSON()); err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	if err = client.Disconnect(ctx); err != nil {
		return err
	}
	return nil
}

// Verify sets m to be verified.
func (m *Member) Verify() error {
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
		Collection("members").
		UpdateOne(ctx,
			bson.D{bson.E{Key: "id", Value: m.ID}},
			bson.D{bson.E{Key: "verified", Value: true}}); err != nil {
		return err
	}

	if err = client.Disconnect(ctx); err != nil {
		return err
	}

	return nil
}
