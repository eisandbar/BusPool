package internal

import (
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const mosquittoServer = "mosquitto:1883"

func NewClient() mqtt.Client {

	opts := mqtt.NewClientOptions().AddBroker(mosquittoServer)
	opts.SetKeepAlive(2 * time.Second)
	opts.SetPingTimeout(1 * time.Second)
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	return client
}
