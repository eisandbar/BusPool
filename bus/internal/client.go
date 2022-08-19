package internal

import (
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func NewClient() mqtt.Client {

	opts := mqtt.NewClientOptions().AddBroker("localhost:1883")
	opts.SetKeepAlive(2 * time.Second)
	opts.SetPingTimeout(1 * time.Second)
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	return client
}
