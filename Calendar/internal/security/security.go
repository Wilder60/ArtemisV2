package security

import (
	"errors"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/Wilder60/ArtemisV2/Calendar/config"
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
		// SigningMethodHS512
		if tkn.Method != jwt.SigningMethodHS256 {
			return nil, errors.New("Invalid Signing Method")
		}
		return []byte(s.config.Security.SecretKey), nil
	})
	if err != nil {
		return err
	}

	if !token.Valid {
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

// CreateToken will create the jwt for the given username that is passed into the function
// NOTE: This is just for debugging purpose and should be removed shortly
// But most likely this will not be removed becasue... reasons
func (sec *Security) CreateToken(username string) (string, error) {
	claims := Claims{
		UserID: username,
		StandardClaims: jwt.StandardClaims{
			Issuer: "Artemis",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(sec.config.Security.SecretKey))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

var SecurityModule = fx.Option(
	fx.Provide(CreateSecurity),
)
