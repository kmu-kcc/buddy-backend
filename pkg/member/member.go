// Package member provides access to the club member of the Buddy System.
package member

import (
	"context"
	"errors"
	"fmt"
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
	ErrPasswordMismatch = errors.New("password mismatch")
	ErrAlreadyMember    = errors.New("already a member")
	ErrUnderReview      = errors.New("under review")
	ErrOnDelete         = errors.New("already on delete")
	ErrNotOnDelete      = errors.New("not on delete")
	ErrOnGraduate       = errors.New("already on graduate")
	ErrNotOnGraduate    = errors.New("not on graduate")
	ErrGraduate         = errors.New("already graduate")
)

// Member represents a club member state.
type Member struct {
	ID         string `json:"id" bson:"id"`                        // student ID
	Password   string `json:"password" bson:"password"`            // password
	Name       string `json:"name" bson:"name"`                    // Name
	Department string `json:"department" bson:"department"`        // department
	Phone      string `json:"phone" bson:"phone"`                  // phone number
	Email      string `json:"email" bson:"email"`                  // e-mail address
	Grade      int    `json:"grade" bson:"grade"`                  // grade
	Attendance int    `json:"attendance" bson:"attendance"`        // attendance status (attending/absent/graduate)
	Approved   bool   `json:"approved" bson:"approved"`            // approved or not
	OnDelete   bool   `json:"on_delete" bson:"on_delete"`          // on exit process or not
	OnGraduate bool   `json:"on_graduate" bson:"on_graduate"`      // on graduation process or not
	CreatedAt  int64  `json:"created_at,string" bson:"created_at"` // when created - Unix timestamp
	UpdatedAt  int64  `json:"updated_at,string" bson:"updated_at"` // last updated - Unix timestamp
}

type Members []Member

// New returns a new club member.
func New(id, name, department, phone, email string, grade, attendance int) *Member {
	now := time.Now().Unix()
	return &Member{
		ID:         id,
		Password:   id,
		Name:       name,
		Department: department,
		Phone:      phone,
		Email:      email,
		Grade:      grade,
		Attendance: attendance,
		Approved:   false,
		OnDelete:   false,
		OnGraduate: false,
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
		return ErrPasswordMismatch
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
func SignUps() (members Members, err error) {
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
			bson.D{
				bson.E{Key: "$set", Value: bson.D{
					bson.E{Key: "approved", Value: true},
					bson.E{Key: "updated_at", Value: time.Now().Unix()}}}}); err != nil {
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
func (m *Member) Exit() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MongoURI))
	if err != nil {
		return err
	}

	if err = client.Database("club").
		Collection("members").
		FindOneAndUpdate(
			ctx,
			bson.D{bson.E{Key: "id", Value: m.ID}},
			bson.D{
				bson.E{Key: "$set", Value: bson.D{
					bson.E{Key: "on_delete", Value: true},
					bson.E{Key: "updated_at", Value: time.Now().Unix()}}}}).
		Decode(m); err != nil {
		return err
	}

	if err = client.Disconnect(ctx); err != nil {
		return err
	}
	if m.OnDelete {
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
func Exits() (members Members, err error) {
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
func (m *Member) CancelExit() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MongoURI))
	if err != nil {
		return err
	}

	if err = client.Database("club").
		Collection("members").
		FindOneAndUpdate(
			ctx,
			bson.D{bson.E{Key: "id", Value: m.ID}},
			bson.D{
				bson.E{Key: "$set", Value: bson.D{
					bson.E{Key: "on_delete", Value: false},
					bson.E{Key: "updated_at", Value: time.Now().Unix()}}}}).
		Decode(m); err != nil {
		return err
	}

	if err = client.Disconnect(ctx); err != nil {
		return err
	}
	if !m.OnDelete {
		return ErrNotOnDelete
	}
	return nil
}

// Search returns the search result filtered by filter.
//
// NOTE:
//
// It is a member-limited operation:
//	Only the authenticated members can access to this operation.
func Search(filter map[string]interface{}) (members Members, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MongoURI))
	if err != nil {
		return
	}

	filter["approved"] = true

	cur, err := client.Database("club").
		Collection("members").
		Find(ctx, filter)

	if err == mongo.ErrNoDocuments {
		return members, client.Disconnect(ctx)
	} else if err != nil {
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

// Update updates the state of m to update.
//
// NOTE:
//
// It is a member-limited operation:
//	Only the authenticated members can access to this operation.
func (m Member) Update(update map[string]interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MongoURI))
	if err != nil {
		return err
	}

	update["updated_at"] = time.Now().Unix()

	if _, err = client.Database("club").
		Collection("members").
		UpdateOne(
			ctx,
			bson.D{bson.E{Key: "id", Value: m.ID}},
			bson.D{bson.E{Key: "$set", Value: update}}); err != nil {
		return err
	}
	return client.Disconnect(ctx)
}

// ApplyGraduate registers m to be a graduate.
//
// NOTE:
//
// It is a member-limited operation:
//	Only the authenticated members can access to this operation.
func (m *Member) ApplyGraduate() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MongoURI))
	if err != nil {
		return err
	}

	collection := client.Database("club").Collection("members")

	if err = collection.FindOne(ctx, bson.D{bson.E{Key: "id", Value: m.ID}}).
		Decode(m); err != nil {
		return err
	}

	if m.Attendance == Graduate {
		if err = client.Disconnect(ctx); err != nil {
			return err
		}
		return ErrGraduate
	}
	if m.OnGraduate {
		if err = client.Disconnect(ctx); err != nil {
			return err
		}
		return ErrOnGraduate
	}

	if _, err = collection.UpdateOne(
		ctx,
		bson.D{bson.E{Key: "id", Value: m.ID}},
		bson.D{bson.E{Key: "$set", Value: bson.D{
			bson.E{Key: "on_graduate", Value: true},
			bson.E{Key: "updated_at", Value: time.Now().Unix()}}}}); err != nil {
		return err
	}
	return client.Disconnect(ctx)
}

// CancelGraduate cancels the graduation request of m.
//
// NOTE:
//
// It is a member-limited operation:
//	Only the authenticated members can access to this operation.
func (m *Member) CancelGraduate() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MongoURI))
	if err != nil {
		return err
	}

	if err = client.Database("club").
		Collection("members").
		FindOneAndUpdate(
			ctx,
			bson.D{bson.E{Key: "id", Value: m.ID}},
			bson.D{bson.E{Key: "$set", Value: bson.D{
				bson.E{Key: "on_graduate", Value: false},
				bson.E{Key: "updated_at", Value: time.Now().Unix()}}}}).
		Decode(m); err != nil {
		return err
	}

	if err = client.Disconnect(ctx); err != nil {
		return err
	}
	if m.Attendance == Graduate {
		return ErrGraduate
	}
	if !m.OnGraduate {
		return ErrNotOnGraduate
	}
	return nil
}

// GraduateApplies returns the graduate request list.
//
// NOTE:
//
// It is a privileged operation:
//	Only the club managers can access to this operation.
func GraduateApplies() (members Members, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MongoURI))
	if err != nil {
		return
	}

	cur, err := client.Database("club").
		Collection("members").
		Find(ctx, bson.D{bson.E{Key: "on_graduate", Value: true}})
	if err != nil {
		return
	}

	memb := new(Member)

	for cur.Next(ctx) {
		if err = cur.Decode(memb); err != nil {
			return
		}
		members = append(members, *memb)
	}

	if err = cur.Close(ctx); err != nil {
		return
	}
	return members, client.Disconnect(ctx)
}

// ApproveGraduate updates m to be a graduate.
//
// NOTE:
//
// It is a privileged operation:
//	Only the club managers can access to this operation.
func ApproveGraduate(ids []string) error {
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
			bson.D{bson.E{Key: "$set", Value: bson.D{
				bson.E{Key: "attendance", Value: Graduate},
				bson.E{Key: "on_graduate", Value: false},
				bson.E{Key: "updated_at", Value: time.Now().Unix()}}}}); err != nil {
		return err
	}
	return client.Disconnect(ctx)
}

// Graduates returns all graduate members.
//
// NOTE:
//
// It is a member-limited operation:
//	Only the authenticated members can access to this operation.
func Graduates() (members Members, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MongoURI))
	if err != nil {
		return
	}

	cur, err := client.Database("club").
		Collection("members").
		Find(ctx, bson.D{bson.E{Key: "attendance", Value: Graduate}})
	if err != nil {
		return
	}

	memb := new(Member)

	for cur.Next(ctx) {
		if err = cur.Decode(memb); err != nil {
			return
		}
		members = append(members, *memb)
	}

	if err = cur.Close(ctx); err != nil {
		return
	}
	return members, client.Disconnect(ctx)
}

// String implements fmt.Stringer.
func (m Member) String() string {
	return fmt.Sprintf(
		"Member {\n  %-12s%s\n  %-12s%s\n  %-12s%s\n  %-12s%s\n  %-12s%s\n  %-12s%s\n  %-12s%d\n  %-12s%d\n  %-12s%t\n  %-12s%t\n  %-12s%t\n  %-12s%d\n  %-12s%d\n}",
		"ID:",
		m.ID,
		"Password:",
		m.Password,
		"Name:",
		m.Name,
		"Department:",
		m.Department,
		"Phone:",
		m.Phone,
		"Email:",
		m.Email,
		"Grade:",
		m.Grade,
		"Attendance:",
		m.Attendance,
		"Approved:",
		m.Approved,
		"OnDelete:",
		m.OnDelete,
		"OnGraduate:",
		m.OnGraduate,
		"CreatedAt:",
		m.CreatedAt,
		"UpdatedAt:",
		m.UpdatedAt,
	)
}
