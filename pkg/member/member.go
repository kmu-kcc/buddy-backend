// Package member provides access to the club member of the Buddy System.
package member

import "time"

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
