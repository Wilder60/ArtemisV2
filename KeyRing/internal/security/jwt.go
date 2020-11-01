package security

import (
	"errors"

	"github.com/Wilder60/KeyRing/configs"
	"github.com/dgrijalva/jwt-go"
	"go.uber.org/fx"
	"go.uber.org/zap"
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
	logger *zap.Logger
}

func CreateDefaultSecurity(cfg *configs.Config, log *zap.Logger) *Security {
	return &Security{
		config: cfg,
		logger: log,
	}
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
