package web

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/Wilder60/ArtemisV2/Calendar/internal/security"
	"github.com/gin-gonic/gin"
)

var validTokenHeader = regexp.MustCompile(`^Bearer [A-Za-z0-9-_=]+\.[A-Za-z0-9-_=]+\.?[A-Za-z0-9-_=]*$`)

type jwt interface {
	Validate(string) error
	GetClaims(string) (*security.Claims, error)
}

// Authorize is used to check the
func Authorize(sec jwt) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		unvalidatedToken, err := getAuthToken(ctx)
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		err = sec.Validate(unvalidatedToken)
		if err != nil {
			fmt.Println(err.Error())
			ctx.AbortWithStatus(http.StatusUnauthorized)
		}
		ctx.Next()
	}
}

// AddUser function will
func AddUser(sec jwt) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token, err := getAuthToken(ctx)
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
		}
		claims, err := sec.GetClaims(token)
		if err != nil {

		}
		ctx.Set("UserID", claims.UserID)
		ctx.Next()
	}
}

func getAuthToken(ctx *gin.Context) (string, error) {
	header := ctx.GetHeader("Authorization")
	valid := validTokenHeader.MatchString(header)
	if !valid {
		return "", errors.New("Invalid token header format")
	}
	token := strings.Split(header, " ")
	// We don't have to check if len(token) >= 2 because regexp validated it
	return token[1], nil
}
