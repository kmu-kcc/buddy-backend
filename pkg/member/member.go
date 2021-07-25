package member

type Member struct {
	ID         string `json:"id" bson:"id"`
	Password   string `json:"password" bson:"password"`
	Name       string `json:"name" bson:"name"`
	Department string `json:"department" bson:"department"`
	Grade      string `json:"grade" bson:"grade"`
	Phone      string `json:"phone" bson:"phone"`
	Email      string `json:"email" bson:"email"`
	Attendance int    `json:"attendance" bson:"attendance"`
	Approved   bool   `json:"approved" bson:"approved"`
	OnDelete   bool   `json:"on_delete" bson:"on_delete"`
	CreatedAt  int64  `json:"created_at" bson:"created_at"`
	UpdatedAt  int64  `json:"updated_at" bson:"updated_at"`
}
