package models

type MobileIdentity struct {
	Id        int64  `json:"id" bson:"_id"`
	Uid       int64  `json:"uid" bson:"uid"`
	Mobile    string `json:"mobile" bson:"mobile"`
	CreatedAt int64  `json:"created_at" bson:"created_at"`
	UpdatedAt int64  `json:"updated_at" bson:"updated_at"`
}
