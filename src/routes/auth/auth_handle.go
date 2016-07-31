package auth

import (
	"net/http"

	"libs/cache"
	"libs/errorcode"
	"libs/odm"
	"models"
	"utils/rand"

	"github.com/qiniu/xlog"
	"gopkg.in/mgo.v2"
)

func GetVerifyCode(log *xlog.Logger, db *odm.DB, cache_ cache.CacheStorage, args *GetVerifyCodeArgs, req *http.Request) (int, interface{}) {

	mi := models.MobileIdentity{}
	icoll := db.C(mi)
	defer icoll.Close()

	err := icoll.Find(odm.M{"mobile": args.Mobile}).One(&mi)
	if err == nil {
		return errorcode.HandleError(errorcode.ErrDBDup)
	}
	if err != nil && err != mgo.ErrNotFound {
		return errorcode.HandleError(err)
	}

	//TODO update mobile

	code := rand.NumString(4)
	log.Info(args.Mobile, code)

	if err := cache_.Set("mobile_code."+args.Mobile, []byte(code)); err != nil {
		return errorcode.HandleError(err)
	}

	if err := cache_.Expire("mobile_code."+args.Mobile, 600*10); err != nil {
		return errorcode.HandleError(err)
	}

	return http.StatusOK, nil
}
