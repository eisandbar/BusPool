package main

import (
	"fmt"
	"time"

	"github.com/eisandbar/BusPool/bus/internal"
)

func main() {

	fmt.Println("Starting")
	for i := 0; i < fleetSize; i++ {
		go internal.InitOne(internal.NewClient(), newBus(i), time.Tick(time.Second))
	}

	fmt.Println("Fleet initialized")

	fmt.Println("Waiting forever")
	forever := make(chan bool)
	<-forever
}
