// Package fee provides access to the club fee of the Buddy System.
package fee

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Log represents a fees history.
type Log struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	MemberID    string             `json:"member_id" bson:"member_id"`
	Description string             `json:"description" bson:"description"`
	Amount      int                `json:"amount" bson:"amount"`
	Type        int                `json:"type" bson:"type"`
	CreatedAt   int64              `json:"created_at,string" bson:"created_at"`
}

type Logs []Log

// NewLog returns a new fee log.
func NewLog(memberID, description string, amount, typ int) *Log {
	return &Log{
		ID:          primitive.NewObjectID(),
		MemberID:    memberID,
		Description: description,
		Amount:      amount,
		Type:        typ,
		CreatedAt:   time.Now().Unix(),
	}
}

// Public returns the limited information of l.
func (l Log) Public() map[string]interface{} {
	pub := make(map[string]interface{})

	pub["description"] = l.Description
	pub["amount"] = l.Amount
	pub["type"] = l.Type
	pub["created_at"] = l.CreatedAt

	return pub
}

// Public returns the limited information of ls.
func (ls Logs) Public() []map[string]interface{} {
	pubs := make([]map[string]interface{}, len(ls))

	for idx, log := range ls {
		pubs[idx] = log.Public()
	}

	return pubs
}
