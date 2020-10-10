package security

import (
	"errors"
	"time"

	"github.com/Wilder60/KeyRing/configs"

	"github.com/dgrijalva/jwt-go"
)

var ErrInvalidToken = errors.New("Token is not valid")

type Claims struct {
	UserID string
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
		UserID: username,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  now.Add(5 * time.Minute).Unix(),
			ExpiresAt: now.Unix(),
			Issuer:    "Artemis",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	cfg := configs.Get()
	signedToken, err := token.SignedString(cfg.Security.SecretKey)
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

// Validate will take a string version of a token and construct the claims for the
func Validate(token string) error {
	claims := &Claims{}

	tkn, err := parseTokenString(token, claims)
	if err != nil {
		return err
	}

	if !tkn.Valid {
		err = ErrInvalidToken
	}
	return err
}

// GetUser is
// This function should only be called after the token is validated so their should be no
// reason to validate the request or maybe idk
func GetUser(token string) (string, error) {
	claims := &Claims{}
	_, err := parseTokenString(token, claims)
	return claims.UserID, err
}

func parseTokenString(token string, claims *Claims) (*jwt.Token, error) {
	return jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		cfg := configs.Get()
		return cfg.Security.SecretKey, nil
	})
}
