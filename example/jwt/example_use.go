package main

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"net/http"
)

func createToken(key []byte) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	tokenString, err := token.SignedString(key)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func main() {

	token, err := createToken([]byte("szy"))
	if err != nil {
		fmt.Println("createToken error : ", err)
		return
	}

}
