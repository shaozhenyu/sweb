package odm

import (
	"sync"

	"gopkg.in/mgo.v2"
)

type DB struct {
	dbName  string
	Session *mgo.Session
	sync.Mutex
	*Option
}

type Option struct {
}

func New(url, dbName string, opt *Option) (*DB, error) {
	Session, err := mgo.Dial(url)
	if err != nil {
		return nil, err
	}

	return &DB{
		dbName:  dbName,
		Session: Session,
		Option:  opt,
	}, nil
}
