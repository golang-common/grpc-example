// @Author: Perry
// @Date  : 2020/5/11
// @Desc  : 

package token

import (
	"github.com/dgrijalva/jwt-go"
	"testing"
	"time"
)

func getClaims() CustomClaims {
	var a []string
	a = append(a, "1")
	ip := net_utils.GetAllAdrress()
	claims := CustomClaims{
		"admin",
		ip,
		a,
		jwt.StandardClaims{
			NotBefore: int64(time.Now().Truncate(1 * time.Hour).Unix()),
			ExpiresAt: int64(time.Now().Add(1 * time.Hour).Unix()),
			Issuer:    "sealion",
		},
	}
	return claims
}
func TestJWT_CreateToken(t *testing.T) {
	j := NewJWT()
	claims := getClaims()
	if token, err := j.CreateToken(claims); err == nil {
		t.Log(token)
	} else {
		t.Error(err)
	}
}

func TestJWT_ParseToken(t *testing.T) {
	j := NewJWT()
	token, _ := j.CreateToken(getClaims())
	if c, err := j.ParseToken(token); err == nil {
		t.Log(c)
		t.Log(time.Unix(c.NotBefore, 0).Format("2006-01-02 15:04:05"))
	} else {
		t.Error(err)
	}
	//newToken, err := j.RefreshToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoiYWRtaW4iLCJpcCI6IjEyNy4wLjAuMSIsImV4cCI6MTUzODU1MjkyOSwiaXNzIjoiZGFpcGVuZ3l1YW4iLCJuYmYiOjE1Mzg1NDY0MDB9.SjHtinCYmqEoIvjS1e1Fq7s5FRiWCaP_YI9ZuXgn7L0")
	//fmt.Println(newToken)
	//fmt.Println(err)
}

