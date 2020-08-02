package jwt

import "github.com/dgrijalva/jwt-go"

var (
	Algorithm = "HS256"

	Secret = []byte("123456")
)

type ClaimsMapper func() jwt.Claims

func CreateTokenString(mapper ClaimsMapper) (token string, err error) {
	mapClaims := mapper
	sign := jwt.New(jwt.GetSigningMethod(Algorithm))
	sign.Claims = mapClaims()
	token, err = sign.SignedString(Secret)
	return
}
