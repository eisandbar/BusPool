package publisher

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	. "github.com/eisandbar/BusPool/lion/typing"
)

type Publisher interface {
	Publish(Bus, Instruction)
}

func NewMQTTPublisher() MQTTPublisher {
	return MQTTPublisher{client: NewClient()}
}

type MQTTPublisher struct {
	client mqtt.Client
}

func (pub MQTTPublisher) Publish(bus Bus, inst Instruction) {
	body, err := json.Marshal(inst)
	if err != nil {
		log.Println("Error marshalling instruction")
	}
	topic := fmt.Sprintf("bus/requests/%d", bus.Id)
	token := pub.client.Publish(topic, 0, false, body)
	token.Wait()
}

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
