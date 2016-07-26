package install

import (
	"net/http"

	"libs/odm"

	"github.com/qiniu/xlog"
)

func GetResource(log *xlog.Logger, db *odm.DB, id int64, collName string, req *http.Request) (int, interface{}) {
	log.Info("GetResource")

	v, err := db.Find2(odm.M{"id": id}, collName)
	if err != nil {
		return 400, err
	}

	return http.StatusOK, v
}
