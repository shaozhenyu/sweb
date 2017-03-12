package auth

import (
	"net/http"
	"time"

	"libs/cache"
	"libs/errorcode"
	"libs/odm"
	"uc/models"
	"utils/rand"

	"github.com/qiniu/xlog"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2"
)

type Options struct {
	Path     string
	Domain   string
	MaxAge   int
	Secure   bool
	HttpOnly bool
}

var DefaultSessionOptions = &Options{
	Path:     "/",
	MaxAge:   86400 * 7,
	HttpOnly: true,
}

const SESSION_NAME = "SESSION_ID"

func NewCookie(name, value string, options *Options) *http.Cookie {
	cookie := &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     options.Path,
		Domain:   options.Domain,
		MaxAge:   options.MaxAge,
		Secure:   options.Secure,
		HttpOnly: options.HttpOnly,
	}

	if cookie.MaxAge > 0 {
		d := time.Duration(options.MaxAge) * time.Second
		cookie.Expires = time.Now().Add(d)
	} else if options.MaxAge < 0 {
		cookie.Expires = time.Unix(1, 0)
	}
	return cookie
}

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

func GenResult(user models.User) (token string, ret map[string]interface{}) {
	token = user.GenToken()
	return token, map[string]interface{}{"token": map[string]interface{}{"token": token, "type": "Bearer"}, "user": user}
}

func ChatLoginHander(log *xlog.Logger, db *odm.DB, args *ChatUserArgs, req *http.Request, w http.ResponseWriter) (int, interface{}) {
	name := args.Name
	password := args.Password

	m := models.User{}
	coll := db.C(&m)
	defer coll.Close()

	err := coll.Find(odm.M{"name": name, "password": password}).One(&m)
	if err != nil {
		return errorcode.HandleError(err)
	}
	return 200, "ok"
}

func ChatRegisterHandler(log *xlog.Logger, db *odm.DB, args *ChatUserArgs, req *http.Request, w http.ResponseWriter) (int, interface{}) {
	newName := args.Name
	m := models.User{}
	coll := db.C(&m)
	defer coll.Close()

	err := coll.Find(odm.M{"name": newName}).One(&m)
	if err != nil && err != mgo.ErrNotFound {
		return errorcode.HandleError(err)
	}
	if err == nil {
		return errorcode.HandleError(errorcode.ErrUserExisted)
	}
	m.Name = newName
	m.Password = args.Password
	if err = db.Insert(&m, nil); err != nil {
		return errorcode.HandleError(err)
	}
	return 200, m
}

func RegisterHandler(log *xlog.Logger, db *odm.DB, cache_ cache.CacheStorage, args *RegisterArgs, req *http.Request, w http.ResponseWriter) (int, interface{}) {
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

	mi := models.MobileIdentity{}
	icoll := db.C(mi)
	defer icoll.Close()

	err = icoll.Find(odm.M{"mobile": args.Mobile}).One(&mi)
	if err == nil {
		return errorcode.HandleError(errorcode.ErrDBDup)
	}
	if err != nil && err != mgo.ErrNotFound {
		return errorcode.HandleError(err)
	}

	user := models.User{}
	user.Password = string(password)
	coll := db.C(user)
	defer coll.Close()
	if err = db.Insert(&user, nil); err != nil {
		return errorcode.HandleError(err)
	}

	mi.Uid = user.Id
	mi.Mobile = args.Mobile
	if err = db.Insert(&mi, nil); err != nil {
		return errorcode.HandleError(err)
	}

	token, ret := GenResult(user)
	cookie := NewCookie(SESSION_NAME, token, DefaultSessionOptions)
	http.SetCookie(w, cookie)
	return 200, ret
}
