package request

import (
	"context"
	"fmt"
	"github.com/HashTagApp/hashlibry"
	"github.com/go-kit/kit/transport/http"
	"github.com/google/uuid"
	http2 "net/http"
	"strings"
)

type contextkey string

const JWTTokenContextKey contextkey = "JWTToken"

const (
	bearer       string = "bearer"
	bearerFormat string = "Bearer %s"
)

func HttpToContext() http.RequestFunc {
	return func(ctx context.Context, req *http2.Request) context.Context {
		ctx = hashlibry.WithUUID(uuid.New())

		//
		appVersion := req.Header.Get("App-Version")
		ctx = hashlibry.WithValue(ctx, "App-Version", appVersion)

		token, ok := extractTokenFromAuthHeader(req.Header.Get("Authorization"))
		if !ok {
			return ctx
		}

		return hashlibry.WithValue(ctx, JWTTokenContextKey, token)
	}
}

func ContextToHttp() http.RequestFunc {
	return func(ctx context.Context, req *http2.Request) context.Context {
		token, ok := ctx.Value(JWTTokenContextKey).(string)
		if ok {
			req.Header.Add("Authorization", generateAuthHeaderFromToken(token))
		}
		return ctx
	}
}

func extractTokenFromAuthHeader(val string) (token string, ok bool) {
	authHeaderParts := strings.Split(val, " ")
	if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != bearer {
		return "", false
	}

	return authHeaderParts[1], true
}

func generateAuthHeaderFromToken(token string) string {
	return fmt.Sprintf(bearerFormat, token)
}
