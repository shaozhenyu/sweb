package auth

import (
	"net/http"

	"libs/cache"
	"libs/errorcode"
	"libs/odm"
	"models"
	"utils/rand"

	"github.com/qiniu/xlog"
	"golang.org/x/crypto/bcrypt"
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

	if err := cache_.Expire("mobile_code."+args.Mobile, 60*10); err != nil {
		return errorcode.HandleError(err)
	}

	return http.StatusOK, nil
}

func RegisterHandler(log *xlog.Logger, db *odm.DB, cache_ cache.CacheStorage, args *RegisterArgs, req *http.Request) (int, interface{}) {
	bs, err := cache_.Get("mobile_code." + args.Mobile)
	if err != nil {
		log.Error(args, err)
		if cache.IsCacheNotFound(err) {
			return errorcode.HandleError(errorcode.ErrDBNotFound)
		}
		return errorcode.HandleError(err)
	}

	code := string(bs)
	if code != args.Code {
		log.Info(code, args)
		return errorcode.HandleError(errorcode.ErrDBNotFound)
	}

	password, err := bcrypt.GenerateFromPassword([]byte(args.Password), bcrypt.DefaultCost)
	if err != nil {
		return errorcode.HandleError(err)
	}

	log.Info("password : ", string(password))

	return 200, nil
}
