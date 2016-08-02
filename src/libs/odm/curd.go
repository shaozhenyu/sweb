package odm

import (
	"encoding/json"
	"io"
	"reflect"
	"time"

	"libs/spec"

	"github.com/fatih/structs"
	//"github.com/qiniu/log"
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

func (db *DB) vv(collName string) (interface{}, error) {
	if t, ok := db.Coll[collName]; ok {
		return reflect.New(reflect.SliceOf(t.Type)).Interface(), nil
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

func (db *DB) SetID(obj interface{}, id interface{}) {
	n := db.CollName(obj)
	m := db.Coll[n]

	elem := reflect.ValueOf(obj).Elem()
	if f := elem.FieldByName(m.IDFieldName); f.IsValid() {
		if i, ok := id.(int64); ok {
			f.SetInt(i)
		}
	}
}

func (db *DB) Insert(v interface{}, collName string) error {
	coll := db.C(v)
	defer coll.Close()

	if db.incr != nil {
		db.SetID(v, db.incr.Incr(coll.FullName))
	}
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

func (db *DB) ListWithSpec(collName string, s *spec.ListQuery) (interface{}, error) {
	return db.List2(collName, s.Cond, s.Sorts, s.Paging.From, s.Paging.Size)
}

func (db *DB) List2(collName string, cond map[string]interface{}, sorts []string, from, size int) (interface{}, error) {
	v, err := db.v(collName)
	if err != nil {
		return nil, err
	}

	vv, err := db.vv(collName)
	if err != nil {
		return nil, err
	}

	coll := db.C(v)
	defer coll.Close()

	q := coll.Find(cond).Sort(sorts...)

	if from < 0 {
		from = 0
	}

	if size < 1 {
		size = 20
	} else if size > 100 {
		size = 100
	}

	total, err := q.Count()
	if err != nil {
		return nil, err
	}

	err = q.Skip(from).Limit(size).All(vv)
	if err != nil {
		return nil, err
	}

	ret := spec.ListResult{
		Result: vv,
		Paging: spec.Paging{
			From:  from,
			Size:  reflect.ValueOf(vv).Elem().Len(),
			Total: total,
		},
	}

	return ret, nil
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
