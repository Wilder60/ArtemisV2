package middleware

import (
	"context"
	"errors"

	"github.com/Wilder60/ArtemisV2/Calendar/internal/security"
)

type contextKey string

var (
	contextKeyAuthToken = contextKey("auth-token")
	contextKeyClaims    = contextKey("claims")
	errValueCast        = errors.New("could not cast context value")
)

func (c contextKey) String() string {
	return string(c)
}

func GetAuthToken(ctx context.Context) (string, error) {
	token, ok := ctx.Value(contextKeyAuthToken).(string)
	if !ok {
		return "", errValueCast
	}
	return token, nil
}

func GetClaims(ctx context.Context) (*security.Claims, error) {
	claims, ok := ctx.Value(contextKeyClaims).(*security.Claims)
	if !ok {
		return nil, errValueCast
	}
	return claims, nil
}
