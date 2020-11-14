package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"go.uber.org/fx"

	"github.com/Wilder60/ArtemisV2/Calendar/internal/security"
	"github.com/gin-gonic/gin"
)

var validTokenHeader = regexp.MustCompile(`^Bearer [A-Za-z0-9-_=]+\.[A-Za-z0-9-_=]+\.?[A-Za-z0-9-_=]*$`)

type HTTP struct {
	Security jwt
}

type jwt interface {
	Validate(string) error
	GetClaims(string) (*security.Claims, error)
}

func ProvideHTTP(sec *security.Security) *HTTP {
	return &HTTP{
		Security: sec,
	}
}

// Authorize is used to check the
func (h *HTTP) Authorize() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		unvalidatedToken, err := h.getAuthToken(ctx)
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		err = h.Security.Validate(unvalidatedToken)
		if err != nil {
			fmt.Println(err.Error())
			ctx.AbortWithStatus(http.StatusUnauthorized)
		}
		ctx.Next()
	}
}

// AddUser function will
func (h *HTTP) AddUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token, err := h.getAuthToken(ctx)
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
		}
		claims, err := h.Security.GetClaims(token)
		if err != nil {

		}
		ctx.Set("UserID", claims.UserID)
		ctx.Next()
	}
}

func (h *HTTP) getAuthToken(ctx *gin.Context) (string, error) {
	header := ctx.GetHeader("Authorization")
	valid := validTokenHeader.MatchString(header)
	if !valid {
		return "", errors.New("Invalid token header format")
	}
	token := strings.Split(header, " ")
	// We don't have to check if len(token) >= 2 because regexp validated it
	return token[1], nil
}

var HTTPMiddlewareModule = fx.Option(
	fx.Provide(ProvideHTTP),
)
