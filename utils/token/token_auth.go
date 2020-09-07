// @Author: Perry
// @Date  : 2020/5/11
// @Desc  : 

package token

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"time"
)

var (
	tokenExpired     = errors.New(`token is expired`)
	tokenNotValidYet = errors.New("token not active yet")
	tokenMalformed   = errors.New("that's not even a token")
	tokenInvalid     = errors.New("token is illegal")
	signKey          = "650faf143358284f941db1b5ebf65deb"
)

type CustomClaims struct {
	Name string   `json:"name"` //对端名称
	IP   string   `json:"ip"`   //对端申请token时的IP地址
	Sub  []string `json:"sub"`  //允许访问的资源
	jwt.StandardClaims
}

type LoginResult struct {
	Token string `json:"token"`
}
type authJWT struct {
	SigningKey []byte
}

func NewJWT() *authJWT {
	return &authJWT{
		SigningKey: []byte(signKey),
	}
}

func (j *authJWT) CreateToken(claims CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SigningKey)
}

func (j *authJWT) ParseToken(tokenString string) (*CustomClaims, error) {
	claims := new(CustomClaims)
	authToken, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, tokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, tokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, tokenNotValidYet
			} else {
				return nil, tokenInvalid
			}
		}
	}
	if authToken != nil {
		if claims, ok := authToken.Claims.(*CustomClaims); ok && authToken.Valid {
			return claims, nil
		}
	}
	return nil, tokenInvalid
}

func (j *authJWT) RefreshToken(tokenString string) (string, error) {
	claims := new(CustomClaims)
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		claims.StandardClaims.ExpiresAt = time.Now().Add(1 * time.Hour).Unix()
		return j.CreateToken(*claims)
	}
	return "", tokenInvalid
}
