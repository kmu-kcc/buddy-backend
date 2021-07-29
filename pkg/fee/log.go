package fee

import "go.mongodb.org/mongo-driver/bson/primitive"

// Log represents a history of fee usage.
type Log struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	MemberID  string             `json:"memberid" bson:"memberid"`
	Amount    int                `json:"amount" bson:"amount"`
	Type      string             `json:"type" bson:"type"`
	CreatedAt int64              `json:"createdat" bson:"createdat"`
	UpdatedAt int64              `json:"updatedat" bson:"updatedat"`
}

func Newlog(memberID, typ string, amount int, created, updated int64) Log {
	return Log{
		ID:        primitive.NewObjectID(),
		MemberID:  memberID,
		Type:      typ,
		Amount:    amount,
		CreatedAt: created,
		UpdatedAt: updated,
	}
}
