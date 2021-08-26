// Package member provides access to the club member of the Buddy System.
package member

// Role represents the member role.
type Role struct {
	Master             bool `json:"-" bson:"master"`
	MemberManagement   bool `json:"member_management" bson:"member_management"`
	ActivityManagement bool `json:"activity_management" bson:"activity_management"`
	FeeManagement      bool `json:"fee_management" bson:"fee_management"`
}

// NewRole returns a new role without any authorities.
func NewRole() *Role { return &Role{} }
