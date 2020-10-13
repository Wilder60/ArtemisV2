package security

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// Authorize will check the authorization for a given request, this will check if they just have a valid token
// not if they are an adminstrator, that will be handed by the AuthorizeAdmin function
func Authorize() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Authorization Header should be Bearer: Token as is standard
		tokenHeader := ctx.GetHeader("Authorization")
		splitHeader := strings.Split(tokenHeader, " ")
		if len(splitHeader) != 2 {
			ctx.AbortWithStatus(http.StatusBadRequest)
		}

		err := Validate(splitHeader[1])
		if err == ErrInvalidToken {
			ctx.AbortWithStatus(http.StatusBadRequest)
			return
		} else if err != nil {
			ctx.AbortWithStatus(http.StatusInternalServerError)
		}
		ctx.Next()
	}
}

func Admin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenHeader := ctx.GetHeader("Authorization")
		splitHeader := strings.Split(tokenHeader, " ")
		if len(splitHeader) != 2 {
			ctx.AbortWithError(http.StatusBadRequest, errors.New("Invalid authorization header format"))
		}

	}

}

func getTokenString() {

}
