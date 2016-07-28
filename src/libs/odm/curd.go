package odm

import (
	"encoding/json"
	"io"
	"reflect"
	"time"

	"github.com/fatih/structs"
	"github.com/qiniu/log"
	"gopkg.in/mgo.v2"
)

type ODMError struct {
	Err string
}

var (
	ErrNotRegisted = &ODMError{
		Err: "not registed",
	}
)

func (o *ODMError) Error() string {
	return o.Err
}

type ODMRefError struct {
	CollName  string `json:"resource_name"`
	FieldName string `json:"field_name"`
}

func (o ODMRefError) Error() string {
	bs, _ := json.Marshal(o)
	return string(bs)
}

func (db *DB) v(collName string) (interface{}, error) {

	if t, ok := db.Coll[collName]; ok {
		return reflect.New(t.Type).Interface(), nil
	}
	return nil, ErrNotRegisted
}

func setUnixTime(obj interface{}, unix int64, fields ...string) {
	elem := reflect.ValueOf(obj).Elem()
	for _, fname := range fields {
		if f := elem.FieldByName(fname); f.IsValid() && f.Type().Kind() == reflect.Int64 {
			f.SetInt(unix)
		}
	}
}

func (db *DB) Find(selector interface{}, v interface{}) error {
	coll := db.C(v)
	defer coll.Close()

	if err := coll.Find(selector).One(v); err != nil {
		return err
	}
	return nil
}

func (db *DB) Find2(selector interface{}, collName string) (interface{}, error) {
	v, err := db.v(collName)
	if err != nil {
		return nil, err
	}

	err = db.Find(selector, v)
	if err != nil {
		return nil, err
	}

	return v, nil
}

func (db *DB) Insert(v interface{}, collName string) error {
	coll := db.C(v)
	defer coll.Close()

	setUnixTime(v, time.Now().UnixNano(), "UpdatedAt", "CreatedAt")

	err := coll.Insert(v)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) Insert2(collName string, reader io.Reader) (interface{}, error) {
	v, err := db.v(collName)
	if err != nil {
		return nil, err
	}

	decoder := json.NewDecoder(reader)
	if err := decoder.Decode(v); err != nil {
		return nil, err
	}

	err = db.Insert(v, collName)
	return v, err
}

func (db *DB) handleReadOnly(v interface{}, changed map[string]interface{}) error {
	//collName_ := db.CollName(v)
	//m := db.Coll[collName_]

	delete(changed, "_id")
	delete(changed, "created_at")
	delete(changed, "updated_at")

	return nil
}

func (db *DB) Update(selector interface{}, v interface{}, collName string, params map[string]string) error {
	coll := db.C(v)
	defer coll.Close()

	structs.DefaultTagName = "bson"
	m := structs.Map(v)

	db.handleReadOnly(collName, m)

	m["updated_at"] = time.Now().UnixNano()

	chg := mgo.Change{
		Update:    M{"$set": m},
		ReturnNew: true,
	}

	_, err := coll.Find(selector).Apply(chg, v)

	return err
}

func (db *DB) Update2(selector interface{}, collName string, r io.Reader, params map[string]string) (interface{}, error) {
	v, err := db.v(collName)
	if err != nil {
		return nil, err
	}

	decoder := json.NewDecoder(r)
	if err := decoder.Decode(v); err != nil {
		return nil, err
	}

	err = db.Update(selector, v, collName, params)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func (db *DB) Remove(selector interface{}, v interface{}) error {
	coll := db.C(v)
	defer coll.Close()

	chg := mgo.Change{Remove: true}

	if _, err := coll.Find(selector).Apply(chg, v); err != nil {
		return err
	}
	return nil
}

func (db *DB) Remove2(selector interface{}, collName string) (interface{}, error) {
	v, err := db.v(collName)
	if err != nil {
		return nil, err
	}

	err = db.Remove(selector, v)
	return v, nil
}