package security

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// Authorize will check the authorization for a given request, this will check if they just have a valid token
// not if they are an adminstrator, that will be handed by the AuthorizeAdmin function
func (sec *Security) Authorize() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Authorization Header should be Bearer: Token as is standard
		tokenHeader := ctx.GetHeader("Authorization")
		splitHeader := strings.Split(tokenHeader, " ")
		if len(splitHeader) != 2 {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		err := sec.Validate(splitHeader[1])
		if err == ErrInvalidToken {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		} else if err != nil {
			sec.logger.Error("Failed to parse")
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		ctx.Next()
	}
}

//Admin will be used to check if a given token is a valid admin token
func (sec *Security) Admin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenHeader := ctx.GetHeader("Authorization")
		splitHeader := strings.Split(tokenHeader, " ")
		if len(splitHeader) != 2 {
			ctx.AbortWithStatus(http.StatusUnauthorized)
		}

		ctx.Next()
	}
}
