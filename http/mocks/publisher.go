// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package mocks

import (
	"github.com/mainflux/mainflux/pkg/messaging"
)

type mockPublisher struct{}

// NewPublisher returns mock message publisher.
func NewPublisher() messaging.Publisher {
	return mockPublisher{}
}

// NewPubSub returns mock message publisher.
func NewPubSub() messaging.PubSub {
	return mockPublisher{}
}

func (pub mockPublisher) Publish(topic string, msg messaging.Message) error {
	return nil
}

func (pub mockPublisher) Subscribe(topic string, handler messaging.MessageHandler) error {
	return nil
}

func (pub mockPublisher) Unsubscribe(topic string) error {
	return nil
}
