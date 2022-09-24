package main

import (
	"fmt"
	"os"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const mosquittoServer = "mosquitto:1883"

func main() {
	fmt.Println("Initializing")
	r := new(rhino)
	r.Init()
	defer r.Close()
	r.InitTopic("bus-positions")
	r.InitTopic("bus-positions-elastic")

	// Connecting to mqtt
	fmt.Println("Connecting to mqtt")
	opts := mqtt.NewClientOptions().AddBroker(mosquittoServer)
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	// Subscribing to bus/positions mqtt topic
	fmt.Println("Subcribing to topic")
	if token := client.Subscribe("bus/positions", 0, r.positionHandler); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

	// Wait forever
	forever := make(chan bool)
	<-forever

}
