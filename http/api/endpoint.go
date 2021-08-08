// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package api

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/mainflux/mainflux/http"
)

func sendMessageEndpoint(svc http.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(publishReq)
		return svc.PubSub(ctx, req.token, req.msg, req.responseTopic)
	}
}
