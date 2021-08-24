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
	ErrIdentityMismatch = errors.New("계정 정보를 확인해주세요")
	ErrAlreadyMember    = errors.New("이미 등록된 사용자입니다")
	ErrUnderReview      = errors.New("승인 검토 중입니다")
	ErrOnDelete         = errors.New("이미 탈퇴 신청하셨습니다")
	ErrAlreadyActive    = errors.New("already active")
	ErrAlreadyInactive  = errors.New("already inactive")
	ErrPermissionDenied = errors.New("권한이 없습니다")
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
	CreatedAt  int64  `json:"created_at,string" bson:"created_at"` // when created - Unix timestamp
	UpdatedAt  int64  `json:"updated_at,string" bson:"updated_at"` // last updated - Unix timestamp
	Role       Role   `json:"role" bson:"role"`                    // role of member
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
		CreatedAt:  now,
		UpdatedAt:  now,
		Role:       Role{},
	}
}

// Public returns the limited informations of m.
func (m Member) Public() map[string]interface{} {
	pub := make(map[string]interface{})

	pub["id"] = m.ID
	pub["name"] = m.Name
	pub["department"] = m.Department
	pub["email"] = m.Email
	pub["grade"] = m.Grade
	pub["role"] = m.Role

	return pub
}

// Public returns the limited informations of ms.
func (ms Members) Public() []map[string]interface{} {
	pubs := make([]map[string]interface{}, len(ms))

	for idx, member := range ms {
		pubs[idx] = member.Public()
	}

	return pubs
}

// SingIn checks whether m is a club member.
//
// NOTE:
//
// It is a member-limited operation:
//	Only the authenticated members can access to this operation.
func (m Member) SingIn() error {
	ctx := context.Background()
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
		return ErrIdentityMismatch
	} else if err != nil {
		return err
	}
	if err = client.Disconnect(ctx); err != nil {
		return err
	}
	if member.ID != "MASTER" && !member.Approved {
		return ErrUnderReview
	}
	if m.Password != member.Password {
		return ErrIdentityMismatch
	}
	return nil
}

// SignUp applies a membership of m.
// If m already exists (approved or not), nothing changes.
// Else it registers an unapproved member.
func (m Member) SignUp() error {
	ctx := context.Background()
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
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MongoURI))
	if err != nil {
		return
	}

	cur, err := client.Database("club").
		Collection("members").
		Find(
			ctx,
			bson.D{
				bson.E{Key: "id", Value: bson.D{bson.E{Key: "$not", Value: "MASTER"}}},
				bson.E{Key: "approved", Value: false}})

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
	ctx := context.Background()
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
	ctx := context.Background()
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
	ctx := context.Background()
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
	ctx := context.Background()
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

// My returns the personal information of m.
func (m *Member) My() (map[string]interface{}, error) {
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MongoURI))
	if err != nil {
		return nil, err
	}

	member := new(Member)

	if err = client.Database("club").
		Collection("members").
		FindOne(
			ctx,
			bson.D{
				bson.E{Key: "id", Value: m.ID}}).
		Decode(member); err == mongo.ErrNoDocuments {
		if err = client.Disconnect(ctx); err != nil {
			return nil, err
		}
		return nil, mongo.ErrNoDocuments
	}
	if err != nil {
		return nil, err
	}

	if m.Password != member.Password {
		if err = client.Disconnect(ctx); err != nil {
			return nil, err
		}
		return nil, ErrIdentityMismatch
	}

	data := member.Public()

	data["password"] = member.Password
	data["phone"] = member.Phone
	data["attendance"] = member.Attendance
	data["approved"] = member.Approved
	data["on_delete"] = member.OnDelete

	return data, client.Disconnect(ctx)
}

// Search returns the search result with query.
//
// NOTE:
//
// It is a member-limited operation:
//	Only the authenticated members can access to this operation.
func Search(query string) (members Members, err error) {
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MongoURI))
	if err != nil {
		return
	}

	filter := bson.D{
		bson.E{Key: "approved", Value: true},
		bson.E{Key: "$or", Value: bson.A{
			bson.D{bson.E{Key: "id", Value: bson.D{bson.E{Key: "$regex", Value: query}}}},
			bson.D{bson.E{Key: "name", Value: bson.D{bson.E{Key: "$regex", Value: query}}}},
			bson.D{bson.E{Key: "department", Value: bson.D{bson.E{Key: "$regex", Value: query}}}},
		}}}

	cur, err := client.Database("club").Collection("members").Find(ctx, filter)

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
	ctx := context.Background()
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

// Active returns the activation status for member signup.
func Active() (bool, error) {
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MongoURI))
	if err != nil {
		return false, err
	}

	active := new(struct {
		Active bool `bson:"active"`
	})

	if err = client.Database("club").
		Collection("signup").
		FindOne(ctx, bson.D{}).
		Decode(active); err != nil {
		return false, err
	}

	return active.Active, client.Disconnect(ctx)
}

// Activate updates the activation status for member signup.
//
// NOTE:
//
// It is a privileged operation:
//	Only the club managers can access to this operation.
func Activate(activate bool) (bool, error) {
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MongoURI))
	if err != nil {
		return false, err
	}

	active := struct {
		Active bool `bson:"active"`
	}{Active: activate}

	if err = client.Database("club").
		Collection("signup").
		FindOneAndUpdate(
			ctx,
			bson.D{},
			bson.D{bson.E{Key: "$set", Value: active}}).
		Decode(&active); err != nil {
		return false, err
	}

	if err = client.Disconnect(ctx); err != nil {
		return activate, err
	}

	if active.Active == activate {
		if activate {
			return active.Active, ErrAlreadyActive
		} else {
			return active.Active, ErrAlreadyInactive
		}
	}

	return activate, nil
}

// Graduates returns all graduate members.
//
// NOTE:
//
// It is a privileged operation:
//	Only the club managers can access to this operation.
func Graduates() (members Members, err error) {
	ctx := context.Background()
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

// UpdateRole updates the role of member of id.
//
// NOTE:
//
// It is a privileged operation:
//	Only the club managers can access to this operation.
func UpdateRole(id string, role Role) error {
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MongoURI))
	if err != nil {
		return err
	}

	if _, err = client.Database("club").
		Collection("members").
		UpdateOne(
			ctx,
			bson.D{bson.E{Key: "id", Value: id}},
			bson.D{bson.E{Key: "$set", Value: bson.D{
				bson.E{Key: "role", Value: role}}}}); err != nil {
		return err
	}

	return client.Disconnect(ctx)
}

// String implements fmt.Stringer.
func (m Member) String() string {
	return fmt.Sprintf(
		"Member {\n  %-12s%s\n  %-12s%s\n  %-12s%s\n  %-12s%s\n  %-12s%s\n  %-12s%s\n  %-12s%d\n  %-12s%d\n  %-12s%t\n  %-12s%t\n  %-12s%d\n  %-12s%d\n}",
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
		"CreatedAt:",
		m.CreatedAt,
		"UpdatedAt:",
		m.UpdatedAt,
	)
}
