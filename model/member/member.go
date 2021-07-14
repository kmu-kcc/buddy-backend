// Package member provides CRUD operations of the club member of the Buddy System.
package member

// Member represents a club member state.
type Member struct {
	ID         string `json:"id"`           // student ID*
	Password   string `json:"-"`            // password - soon deprecated
	Name       string `json:"name"`         // Name
	Department string `json:"department"`   // department - magic number needed
	Grade      uint   `json:"grade,string"` // grade
	Phone      string `json:"phone"`        // phone number
	Email      string `json:"email"`        // e-mail address
	Enrollment string `json:"enrollment"`   // enrollment state (attending/absent/graduated) - magic number needed
	Approved   bool   `json:"approved"`     // approved or not
}

// New registers a new club member.
// If the member already exists (approved or not), nothing changes.
// Else it registers an unapproved member.
func New(id string, name string, department string, grade uint, phone string, email string, enrollment string) (*Member, error) {
	// TODO
	// MongoDB query
	// e-mail varification (OAuth)
	return &Member{
		ID:         id,
		Password:   id,
		Name:       name,
		Department: department,
		Grade:      grade,
		Phone:      phone,
		Email:      email,
		Enrollment: enrollment,
		Approved:   false,
	}, nil
}
