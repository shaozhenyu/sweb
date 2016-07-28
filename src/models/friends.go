package models

type Friends struct {
	Id        int64  `json:"id" bson:"_id"`
	Name      string `json:"name" bson:"name"`
	Age       int64  `json:"age" bson:"age"`
	Sex       string `json:"sex" bson:"sex"`
	CreatedAt int64  `json:"created_at" bson:"created_at"`
	UpdatedAt int64  `json:"updated_at" bson:"updated_at"`
}
