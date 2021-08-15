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
	Amount    int                `json:"amount,string" bson:"amount"`
	Type      int                `json:"type" bson:"type"`
	CreatedAt int64              `json:"created_at,string" bson:"created_at"`
}

type Logs []Log

// NewLog returns a new fee log.
func NewLog(memberID string, typ, amount int) *Log {
	return &Log{
		ID:        primitive.NewObjectID(),
		MemberID:  memberID,
		Amount:    amount,
		Type:      typ,
		CreatedAt: time.Now().Unix(),
	}
}

// Memfilter returns limited information of member.
func (ls Logs) Logfilter() (res []map[string]interface{}) {
	for _, log := range ls {
		res = append(res, map[string]interface{}{
			"id":         log.ID,
			"member_id":  log.MemberID,
			"amount":     log.Amount,
			"type":       log.Type,
			"created_at": log.CreatedAt,
		})
	}
	return
}
