package bind

import (
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"libs/errorcode"
	libBind "libs/install/bind"
	"libs/odm"
	"uc/models"
	tokenPkg "utils/token"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-martini/martini"
	"github.com/qiniu/xlog"
)

func BindAuthUser() martini.Handler {
	return func(context martini.Context, req *http.Request) {
		var err error
		defer func() {
			if err != nil {
				log := context.Get(reflect.TypeOf((*xlog.Logger)(nil))).Interface().(*xlog.Logger)
				log.Error(err)
				context.Map(err)
				context.Invoke(libBind.ErrorHandle)
				return
			}
		}()

		db := context.Get(reflect.TypeOf((*odm.DB)(nil))).Interface().(*odm.DB)
		u, err1 := GetUserFromReq(db, req)
		if err1 != nil {
			return
		}
		context.MapTo(u, (*odm.IUser)(nil))
	}
}

func GetUserFromReq(db *odm.DB, req *http.Request) (*models.User, error) {
	token, err := getTokenFromHeader(req)
	if err != nil {
		return nil, err
	}
	return getUserFromToken(db, token)
}

func getTokenFromHeader(req *http.Request) (token string, err error) {
	c, err := req.Cookie("SESSION_ID")
	if err != nil {
		token, err = getTokenFromAuthorization(req)
		if err != nil {
			return "", err
		}
	} else {
		token = c.Value
	}

	return token, nil
}

func getTokenFromAuthorization(req *http.Request) (token string, err error) {
	header := strings.TrimSpace(req.Header.Get("Authorization"))
	if header == "" {
		return "", errorcode.ErrApiAuthInvalidHeader
	}

	spice := strings.Split(header, " ")
	if len(spice) != 2 {
		return "", errorcode.ErrApiAuthInvalidToken
	}

	if spice[0] != "Bearer" {
		return "", errorcode.ErrApiAuthInvalidType
	}

	if token = strings.TrimSpace(spice[1]); token == "" {
		err = errorcode.ErrApiAuthInvalidToken
	}
	return
}

func getUserFromToken(db *odm.DB, token string) (*models.User, error) {
	u := models.User{}

	var err error
	_, err = tokenPkg.ParseTokenWithFunc(token, func(jwt1 *jwt.Token) (interface{}, error) {
		uidStr, ok := jwt1.Claims["uid"].(string)
		uid := 0
		uid, err = strconv.Atoi(uidStr)
		if err != nil {
			err = errorcode.ErrApiAuthInvalidToken
			return nil, err
		}

		if !ok {
			err = errorcode.ErrApiAuthInvalidToken
			return nil, err
		}

		coll := db.C(u)
		defer coll.Close()

		if err1 := coll.FindId(uid).One(&u); err1 != nil {
			err = errorcode.ErrApiAuthUserNotFound
			return nil, err
		}
		return u.APIKey, nil
	})

	if err != nil {
		if err2, ok := err.(*jwt.ValidationError); ok {
			if err2.Errors == jwt.ValidationErrorExpired {
				return nil, errorcode.ErrApiAuthExpiredToken
			}
		}
		return nil, errorcode.ErrApiAuthInvalidToken
	}
	return &u, nil
}
