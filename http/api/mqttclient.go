package api

import (
	"errors"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var qos byte = 2

func NewClient(address string, username, password string, timeout time.Duration) (mqtt.Client, error) {
	opts := mqtt.NewClientOptions().
		SetUsername(username).
		SetPassword(password).
		AddBroker(address)
	client := mqtt.NewClient(opts)
	token := client.Connect()
	if token.Error() != nil {
		return nil, token.Error()
	}

	ok := token.WaitTimeout(timeout)
	if ok && token.Error() != nil {
		return nil, token.Error()
	}
	if !ok {
		return nil, errors.New("failed to connect to MQTT broker")
	}

	return client, nil
}

func Requester(client mqtt.Client, topicSub, topicPub string, timeout time.Duration, payload []byte) ([]byte, error) {
	messageChan := make(chan []byte)
	token := client.Subscribe(topicSub, qos, func(client mqtt.Client, message mqtt.Message) {
		messageChan <- message.Payload()
	})
	if token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}

	defer func() {
		token = client.Unsubscribe(topicSub)
		token.Wait()
	}()

	client.Publish(topicPub, qos, false, payload)
	if token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}

	select {
	case msg := <-messageChan:
		return msg, nil
	case <-time.After(timeout):
		return nil, errors.New("recevie response timeout")
	}
}
