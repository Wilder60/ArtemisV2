package middleware

import (
	"context"

	"go.uber.org/fx"

	"github.com/Wilder60/ArtemisV2/Calendar/internal/security"

	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
)

type GRPC struct {
	Security *security.Security
}

func ProvideGRPC(sec *security.Security) *GRPC {
	return &GRPC{
		Security: sec,
	}
}

// GRPCAuthFunc
func (g *GRPC) GRPCAuthFunc(ctx context.Context) (context.Context, error) {
	token, err := grpc_auth.AuthFromMD(ctx, "bearer")
	if err != nil {
		return nil, err
	}

	err = g.Security.Validate(token)
	if err != nil {
		return nil, err
	}

	claims, err := g.Security.GetClaims(token)
	if err != nil {
		return nil, err
	}

	newCtx := context.WithValue(ctx, contextKeyClaims, claims)
	newCtx = context.WithValue(newCtx, contextKeyAuthToken, token)

	return newCtx, nil
}

var GRPCMiddlewareModule = fx.Option(
	fx.Provide(ProvideGRPC),
)
