package install

import (
	"io"
	"net/http"

	"libs/errorcode"
	"libs/odm"
	"libs/spec"

	"github.com/qiniu/xlog"
)

var (
	bodyLimit = 8 << 20
)

func GetResource(log *xlog.Logger, db *odm.DB, id int64, collName string, req *http.Request) (int, interface{}) {

	v, err := db.Find2(odm.M{"_id": id}, collName)
	if err != nil {
		log.Error(err)
		return errorcode.HandleError(err)
	}

	return http.StatusOK, v
}

func PostResource(log *xlog.Logger, db *odm.DB, collName string, req *http.Request) (int, interface{}) {
	limitReader := io.LimitReader(req.Body, int64(bodyLimit))

	v, err := db.Insert2(collName, limitReader)
	if err != nil {
		log.Error(err)
		return errorcode.HandleError(err)
	}
	return http.StatusCreated, v
}

func ListResource(log *xlog.Logger, db *odm.DB, collName string, params map[string]string, req *http.Request) (int, interface{}) {

	listq := spec.NewListQueryWithReq(db.Coll[collName], req)
	v, err := db.ListWithSpec(collName, listq)
	if err != nil {
		log.Error(err)
		return errorcode.HandleError(err)
	}

	return http.StatusOK, v
}

func PutResource(log *xlog.Logger, db *odm.DB, id int64, collName string, req *http.Request, params map[string]string) (int, interface{}) {
	limitReader := io.LimitReader(req.Body, int64(bodyLimit))

	v, err := db.Update2(odm.M{"_id": id}, collName, limitReader, params)
	if err != nil {
		return errorcode.HandleError(err)
	}
	return http.StatusOK, v
}

func DeleteResource(log *xlog.Logger, db *odm.DB, id int64, collName string, req *http.Request) (int, interface{}) {
	v, err := db.Remove2(odm.M{"_id": id}, collName)
	if err != nil {
		return errorcode.HandleError(err)
	}
	return http.StatusOK, v
}
