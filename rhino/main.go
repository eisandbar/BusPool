package main

import (
	"fmt"
	"os"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func main() {
	r := new(rhino)
	r.Init()
	defer r.Close()
	r.InitTopic("bus-positions")

	// Connecting to mqtt
	opts := mqtt.NewClientOptions().AddBroker("localhost:1883")
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	// Subscribing to bus/positions mqtt topic
	if token := client.Subscribe("bus/positions", 0, r.positionHandler); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

	// Wait forever
	forever := make(chan bool)
	<-forever

}
