package public

import (
	"errors"
	"github.com/golang-jwt/jwt/v4"
)

func JwtDecode(tokenString string) (*jwt.StandardClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(JwtSignKey), nil
	})
	//内部进行超时验证
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*jwt.StandardClaims); ok {

		return claims, nil
	} else {
		return nil, errors.New("token is not jwt.StandardClaims")
	}
}

func JwtEncode(claims jwt.StandardClaims) (string, error) {
	mySigningKey := []byte(JwtSignKey)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(mySigningKey)
}
