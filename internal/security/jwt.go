package security

import (
	"errors"
	"time"

	"go.uber.org/fx"

	"github.com/Wilder60/KeyRing/configs"
	"github.com/dgrijalva/jwt-go"
)

// This can be changed to package private now that the middleware has been
// moved into the security module
var ErrInvalidToken = errors.New("Token is not valid")
var audience = "keyring"

var config *configs.Config

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
// Claims is an extended version of the struct to be added to the Standarded
// Claims added above
type Claims struct {
	UserID string
	jwt.StandardClaims
}

// Security is the struct
type Security struct {
	config *configs.Config
}

func CreateDefaultSecurity(cfg *configs.Config) *Security {
	return &Security{
		config: cfg,
	}
}

// CreateToken will create the jwt for the given username that is passed into the function
// NOTE: This is just for debugging purpose and should be removed shortly
// But most likely this will not be removed becasue... reasons
func (sec *Security) CreateToken(username string) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID: username,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  now.Unix(),
			ExpiresAt: now.Add(5 * time.Minute).Unix(),
			Issuer:    "Artemis",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(config.Security.SecretKey))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

// Validate will take a string version of a token and construct the claims for the
func (sec *Security) Validate(token string) error {
	claims := &Claims{}

	tkn, err := sec.parseTokenString(token, claims)
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
func (sec *Security) GetUserFromToken(token string) string {
	claims := &Claims{}
	sec.parseTokenString(token, claims)
	return claims.UserID
}

func (sec *Security) parseTokenString(token string, claims *Claims) (*jwt.Token, error) {
	return jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Security.SecretKey), nil
	})
}

var SecurityModule = fx.Option(
	fx.Provide(CreateDefaultSecurity),
)
