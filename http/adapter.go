// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

// Package http contains the domain concept definitions needed to support
// Mainflux http adapter service functionality.
package http

import (
	"context"
	"errors"
	"time"

	"github.com/mainflux/mainflux"
	"github.com/mainflux/mainflux/pkg/messaging"
)

// Service specifies coap service API.
type Service interface {
	// Publish and Subscribe Messssage
	PubSub(ctx context.Context, token string, msg messaging.Message, subtopicResponse string) ([]byte, error)
}

var _ Service = (*adapterService)(nil)

type adapterService struct {
	pubsub          messaging.PubSub
	things          mainflux.ThingsServiceClient
	responseTimeout time.Duration
}

// New instantiates the HTTP adapter implementation.
func New(pupsub messaging.PubSub, things mainflux.ThingsServiceClient, responseTimeout time.Duration) Service {
	return &adapterService{
		pubsub:          pupsub,
		things:          things,
		responseTimeout: responseTimeout,
	}
}

func (as *adapterService) PubSub(ctx context.Context, token string, msg messaging.Message, responseTopic string) ([]byte, error) {
	ar := &mainflux.AccessByKeyReq{
		Token:  token,
		ChanID: msg.Channel,
	}
	thid, err := as.things.CanAccessByKey(ctx, ar)
	if err != nil {
		return nil, err
	}

	var msgChan chan messaging.Message
	if responseTopic != "" {
		msgChan = make(chan messaging.Message)

		as.pubsub.Subscribe(responseTopic, func(msg messaging.Message) error {
			msgChan <- msg
			return nil
		})
		defer as.pubsub.Unsubscribe(responseTopic)
	}

	msg.Publisher = thid.GetValue()
	err = as.pubsub.Publish(msg.Channel, msg)
	if err != nil {
		return nil, err
	}

	if responseTopic != "" {
		select {
		case msg := <-msgChan:
			return msg.Payload, nil
		case <-time.After(as.responseTimeout):
			return nil, errors.New("wait response timeout")
		}
	}

	return nil, nil
}
