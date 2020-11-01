package security

import (
	"errors"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/Wilder60/ShadowKeep/config"
	"github.com/dgrijalva/jwt-go"
)

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

// The Security struct is the module that will contain
type Security struct {
	logger *zap.Logger
	config *config.Config
}

// CreateSecurity is the Provder function for fx
func CreateSecurity(cfg *config.Config, log *zap.Logger) *Security {
	return &Security{
		logger: log,
		config: cfg,
	}
}

// Validate will take a JWT token string from the Authorization header
// and parse it into the
func (s *Security) Validate(tokenString string) error {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(tkn *jwt.Token) (interface{}, error) {
		if tkn.Method != jwt.SigningMethodHS512 {
			return nil, errors.New("Invalid Signing Method")
		}
		return []byte(s.config.Security.SecretKey), nil
	})

	if token.Valid {
		err = errors.New("Token is not valid")
	}
	return err
}

// GetClaims will parse the token and return the customs claims defined here
// You don't need to check the
func (s *Security) GetClaims(token string) (*Claims, error) {
	claims := &Claims{}
	_, err := jwt.ParseWithClaims(token, claims, func(tkn *jwt.Token) (interface{}, error) {
		return []byte(s.config.Security.SecretKey), nil
	})
	return claims, err
}

var SecurityModule = fx.Option(
	fx.Provide(CreateSecurity),
)
