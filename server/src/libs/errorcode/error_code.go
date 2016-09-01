package errorcode

import (
	"encoding/json"
	"fmt"
	"net/http"

	"libs/odm"

	"github.com/asaskevich/govalidator"
	"gopkg.in/mgo.v2"
)

var (
	errors = map[int]*ErrCodeMessage{}
)

func New(httpStatus, apiCode int, message string) *ErrCodeMessage {
	return (&ErrCodeMessage{httpStatus, apiCode, message}).Register()
}

type ErrCodeMessage struct {
	HttpStatus int    `json:"-"`
	Code       int    `json:"code"`
	Message    string `json:"message"`
}

func (e *ErrCodeMessage) Error() string {
	return fmt.Sprintf(`{"code":%d,"message":"%s"}`, e.Code, e.Message)
}

func (e *ErrCodeMessage) Register() *ErrCodeMessage {
	if _, ok := errors[e.Code]; ok {
		panic(fmt.Sprintf("error code:%d,message:%s has been taken", e.Code, e.Message))
	}
	errors[e.Code] = e
	return e
}

var (
	ErrUnkown     = New(599, 90000, "UNKOWN")
	ErrDBNotFound = New(404, 404, "Not Found")
	ErrDBDup      = New(400, 2000, "dup err")
)

func ErrBadRequest(value interface{}) *ErrCodeMessage {
	if _, ok := value.(string); ok {
		return &ErrCodeMessage{400, 400, "请求参数错误:" + value.(string)}
	}
	bs, _ := json.Marshal(value)
	return &ErrCodeMessage{400, 400, "请求参数错误:" + string(bs)}
}

// common error
var (
	ErrBadRequestBody      = New(400, 400, "请求格式错误")
	ErrNeedLogin           = New(401, 401, "请先登录")
	ErrNoPermission        = New(403, 403, "无权限")
	ErrTooLargeRequestBody = New(http.StatusRequestEntityTooLarge, http.StatusRequestEntityTooLarge, "请求体太大")
)

// auth error
var (
	ErrApiAuthInvalidHeader = New(401, 20000, "Authorization not in request headers")
	ErrApiAuthInvalidToken  = New(401, 20001, `invalid token`)
	ErrApiAuthInvalidType   = New(401, 20002, "invalid auth type")
	ErrApiAuthExpiredToken  = New(498, 20003, "expired token")
	ErrApiAuthUserNotFound  = New(403, 20004, "user info not found")
	ErrApiAuthNeedSSL       = New(403, 20005, "please access api with ssl")
	ErrPasswordNotMatch     = New(400, 20006, "password not match")
	ErrApiAuth              = New(403, 4030, "无权限")
)

var (
	ErrSSODigestNotMatch = New(400, 21000, "sso digest not match")
)

func HandleError(err error) (int, error) {
	if mgo.IsDup(err) {
		return 400, ErrDBDup
	}
	if mgo.ErrNotFound == err {
		return 404, ErrDBNotFound
	}
	switch err.(type) {
	case *ErrCodeMessage:
		{
			e := err.(*ErrCodeMessage)
			return e.HttpStatus, e
		}
	case odm.ODMRefError:
		{
			return http.StatusBadRequest, &ErrCodeMessage{http.StatusBadRequest, 4040, err.Error()}
		}
	// case odm.AuthError:
	// 	{
	// 		return ErrNeedLogin.HttpStatus, ErrNeedLogin
	// 	}
	case govalidator.Error:
	case govalidator.Errors:
		{
			return http.StatusBadRequest, &ErrCodeMessage{http.StatusBadRequest, 4000, err.Error()}
		}
	}
	return 599, ErrUnkown
}

// SendMsg Error
var (
	ErrSendMsgTooQuick = New(405, 405, `send msg too quickly`)
)
