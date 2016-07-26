package odm

import (
	"encoding/json"
	"io"
	"net/http"
	"reflect"
	"time"

	"github.com/qiniu/log"
)

type ODMError struct {
	Err string
}

var (
	bodyLimit = 8 << 20
)

var (
	ErrNotRegisted = &ODMError{
		Err: "not registed",
	}
)

func (o *ODMError) Error() string {
	return o.Err
}

func (db *DB) Close() {
	db.Session.Close()
}

func (db *DB) v(collName string) (interface{}, error) {

	if t, ok := db.Coll[collName]; ok {
		log.Info("collname : ", collName)
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

func (db *DB) Find(selector interface{}, v interface{}) error {
	coll := db.C(v)
	//defer db.close()

	if err := coll.Find(selector).One(v); err != nil {
		log.Info("eeeeeeeeeee")
		return err
	}
	return nil
}

func (db *DB) Insert2(collName string, req *http.Request) (interface{}, error) {
	v, err := db.v(collName)
	if err != nil {
		return nil, err
	}

	limitReader := io.LimitReader(req.Body, int64(bodyLimit))

	decoder := json.NewDecoder(limitReader)
	if err := decoder.Decode(v); err != nil {
		return nil, err
	}

	err = db.Insert(v, collName)
	if err != nil {
		return nil, err
	}

	return v, nil
}

func (db *DB) Insert(v interface{}, collName string) error {
	coll := db.Session.DB(db.dbName).C(collName)
	//defer db.Close()

	setUnixTime(v, time.Now().UnixNano(), "UpdatedAt", "CreatedAt")

	err := coll.Insert(v)
	if err != nil {
		return err
	}
	return nil
}
