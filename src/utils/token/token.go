package token

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var (
	ErrTokenExpire = fmt.Errorf("token expired")
)

func GenrateToken(data map[string]interface{}, expire time.Duration, key []byte) string {
	token := jwt.New(jwt.GetSigningMethod("HS256"))
	if data != nil {
		token.Claims = data
	} else {
		token.Claims = make(map[string]interface{})
	}

	token.Claims["exp"] = time.Now().Add(expire).Unix()
	signed, _ := token.SignedString(key)
	return signed

}

func ParseTokenWithFunc(str string, keyFunc func(token *jwt.Token) (interface{}, error)) (*jwt.Token, error) {
	token, err := jwt.Parse(str, keyFunc)
	if err == nil {
		return token, nil
	}

	if validtionErr, ok := err.(*jwt.ValidationError); ok {
		if validtionErr.Errors == jwt.ValidationErrorExpired {
			return token, err
		}
	}
	return token, err
}
