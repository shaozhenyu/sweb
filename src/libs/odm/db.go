package odm

import (
	"reflect"
	"strings"
	"sync"

	"github.com/qiniu/log"
	"github.com/qor/inflection"
	"gopkg.in/mgo.v2"
)

type M map[string]interface{}

type DB struct {
	dbName  string
	Session *mgo.Session
	Coll    map[string]*Model
	sync.Mutex
	*Option
}

type Collection struct {
	*mgo.Collection
}

type Option struct {
}

func New(url, dbName string, opt *Option) (*DB, error) {
	Session, err := mgo.Dial(url)
	if err != nil {
		return nil, err
	}

	coll := map[string]*Model{}

	return &DB{
		dbName:  dbName,
		Session: Session,
		Coll:    coll,
		Option:  opt,
	}, nil
}

func (d *DB) NewGroup(v ...interface{}) error {
	log.Info("newGroup")

	d.Lock()
	defer d.Unlock()

	coll := map[string]*Model{}

	for _, v_ := range v {
		collname := collName(v_)
		if reflect.TypeOf(v_).Kind() != reflect.Struct {
			return &ODMError{Err: "should be struct , not ptr " + collname}
		}

		coll[collname] = newModel(v_)
		if _, has := d.Coll[collname]; !has {
			d.Coll[collname] = coll[collname]
		}
	}
	return nil
}

func collName(v interface{}) string {
	type_ := reflect.TypeOf(v)
	if type_.Kind() == reflect.Ptr {
		return strings.ToLower(inflection.Plural(type_.Elem().Name()))
	} else if type_.Kind() == reflect.Struct {
		return strings.ToLower(inflection.Plural(type_.Name()))
	}
	panic("should be struct ptr")
}

func (db *DB) CollName(v interface{}) string {
	return collName(v)
}

func (db *DB) C(v interface{}) *Collection {
	collName := db.CollName(v)
	return db.C2(collName)
}

func (db *DB) C2(collName string) *Collection {
	if _, ok := db.Coll[collName]; ok {
		return &Collection{db.Session.Copy().DB(db.dbName).C(collName)}
	}
	panic("not regiested struct " + collName)
}

type Model struct {
	Val interface{}
	reflect.Type
}

func newModel(val interface{}) *Model {
	ret := &Model{Val: val, Type: reflect.TypeOf(val)}
	return ret
}
