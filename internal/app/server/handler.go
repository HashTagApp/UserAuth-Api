package server

import (
	"context"
	"github.com/HashTagApp/UserAuth-Api/internal/app/services"
	"github.com/HashTagApp/UserAuth-Api/internal/domain/errors"
	httpRequest "github.com/HashTagApp/UserAuth-Api/transport/http/request"
	httpResponse "github.com/HashTagApp/UserAuth-Api/transport/http/response"
	"net/http"
)

import httpKitTransport "github.com/go-kit/kit/transport/http"

const XTokenRefresh = `X-Token-Refresh`

var opts = []httpKitTransport.ServerOption{
	// httpKitTransport.ServerErrorLogger(new(errorHandlers.ErrorLogger)),
	httpKitTransport.ServerErrorEncoder(errors.ErrorEncoder),
	httpKitTransport.ServerBefore(httpRequest.HttpToContext()),
}

func Ping() (handler http.Handler) {
	var PingingService services.PingService
	{
		PingingService = services.PingSvc{}
		//PingingService = instrumenting.PingStruct{Next: PingingService}
	}
	return httpKitTransport.NewServer(
		func(ctx context.Context, request interface{}) (response interface{}, err error) {
			return PingingService.Pinging(ctx)
		},
		httpRequest.DecodePingRequest,
		httpResponse.EncodePing,
		opts...,
	)
}
