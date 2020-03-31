package utils

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

func GetToken()(tokenStr string,err error){
	claims := &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(30 * time.Second).Unix(),
		Issuer:    "test",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	if tokenStr, err = token.SignedString([]byte("jey.sign")); err != nil {
		return "",err
	}
	return
}

func ValidateToken(tokenStr string) error {
	_, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte("jey.sign"), nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
				return err
			} else {
				return err
			}
		}
		return err
	}
	return nil
}