// Package fee provides access to the club fee of the Buddy System.
package fee

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Log represents a fees history.
type Log struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	MemberID  string             `json:"member_id" bson:"member_id"`
	Amount    int                `json:"amount" bson:"amount"`
	Type      int                `json:"type" bson:"type"`
	CreatedAt int64              `json:"created_at,string" bson:"created_at"`
}

type Logs []Log

// NewLog returns a new fee log.
func NewLog(memberID string, amount, typ int) *Log {
	return &Log{
		ID:        primitive.NewObjectID(),
		MemberID:  memberID,
		Amount:    amount,
		Type:      typ,
		CreatedAt: time.Now().Unix(),
	}
}

// Public returns the limited information of ls.
func (ls Logs) Public() []map[string]interface{} {
	pubs := make([]map[string]interface{}, len(ls))

	for idx, l := range ls {
		pubs[idx] = make(map[string]interface{})
		pubs[idx]["member_id"] = l.MemberID
		pubs[idx]["amount"] = l.Amount
		pubs[idx]["type"] = l.Type
		pubs[idx]["created_at"] = l.CreatedAt
	}

	return pubs
}
