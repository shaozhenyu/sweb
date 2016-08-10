package models

import (
	"strconv"
	"time"

	"libs/odm"

	"utils/apikey"
	"utils/token"
)

type User struct {
	Id          int64  `json:"id" bson:"_id"`
	Gender      string `json:"gender" bson:"gender"`
	Name        string `json:"name" bson:"name"`
	Mobile      string `json:"mobile" bson:"-"`
	Birthday    string `json:"birthday" bson:"birthday"` //1990-03-23
	CreatedAt   int64  `json:"created_at" bson:"created_at"`
	UpdatedAt   int64  `json:"updated_at" bson:"updated_at"`
	APIKey      []byte `json:"-" bson:"api_key" readonly:"api_key"`
	Password    string `json:"-" bson:"password" readonly:"password"`
	AllowMethod string `json:"-" bson:"-"`
}

func (u User) AllowedMethod() string {
	return u.AllowMethod
}

type MobileIdentity struct {
	Id          int64  `json:"id" bson:"_id"`
	Uid         int64  `json:"uid" bson:"uid"`
	Mobile      string `json:"mobile" bson:"mobile"`
	CreatedAt   int64  `json:"created_at" bson:"created_at"`
	UpdatedAt   int64  `json:"updated_at" bson:"updated_at"`
	AllowMethod string `json:"-" bson:"-"`
}

const (
	tokenExpireInSec = 3600 * 24 * 30 * 12 * 3 * time.Second
)

func (u *User) GenToken() string {
	return token.GenrateToken(
		map[string]interface{}{"uid": strconv.Itoa(int(u.Id))},
		tokenExpireInSec,
		u.APIKey,
	)
}

func (u *User) BeforeInsert(db *odm.DB) error {
	u.APIKey = apikey.Gen(16)
	return nil
}
