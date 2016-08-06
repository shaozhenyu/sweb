package main

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

func main() {

	//use HS256 to create a token
	token := jwt.New(jwt.GetSigningMethod("HS256"))
	//token.Claims = make(map[string]interface{})

	expireTime := 5 * time.Second
	//set expire time
	token.Claims["exp"] = time.Now().Add(expireTime).Unix()

	key := []byte("szy")

	//generate token
	tokenString, err := token.SignedString(key)
	if err != nil {
		fmt.Println("err : ", err)
		return
	}
	fmt.Println("new token : ", tokenString)
	fmt.Println("expire time : ", time.Now().Add(expireTime))

	time.Sleep(10 * time.Second)
	//parse key
	token1, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte("szy"), nil
	})

	if token1.Valid {
		fmt.Println("validate success")
	} else if ve, ok := err.(*jwt.ValidationError); ok {
		if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			fmt.Println("time error")
		} else {
			fmt.Println("could not handle this token : ", err)
		}
	} else {
		fmt.Println("could not handle this token : ", err)
	}

}
