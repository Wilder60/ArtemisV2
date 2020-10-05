package security

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Claims struct {
	Username string
	jwt.StandardClaims
}

/*

type StandardClaims struct {
	Audience  string `json:"aud,omitempty"`
	ExpiresAt int64  `json:"exp,omitempty"`
	Id        string `json:"jti,omitempty"`
	IssuedAt  int64  `json:"iat,omitempty"`
	Issuer    string `json:"iss,omitempty"`
	NotBefore int64  `json:"nbf,omitempty"`
	Subject   string `json:"sub,omitempty"`
}

*/

// CreateToken will create the jwt for the given username that is passed into the function
// NOTE: This is just for debugging purpose and should be removed shortly
// But most likely this will not be removed becasue... reasons
func CreateToken(username string) (string, error) {
	now := time.Now()
	claims := Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  now.Add(5 * time.Minute).Unix(),
			ExpiresAt: now.Unix(),
			Issuer:    "Artemis",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString("123ewfasdgarga")
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

// Validate will take a string version of a token and construct the claims for the
func Validate(token string) error {
	claims := &Claims{}

	tkn, err := jwt.ParseWithClaims(token, claims,
		func(token *jwt.Token) (interface{}, error) {
			// TODO pull this from enviorment file rather then hardcoded string
			return "123ewfasdgarga", nil
		})

	if err != nil && !tkn.Valid {
		err = errors.New("Token is not valid")
	}
	return err
}
