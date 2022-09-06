package main

import (
	"fmt"
	"time"

	"github.com/eisandbar/BusPool/bus/internal"
)

func main() {
	coords := internal.LoadData()

	fmt.Println("Starting")
	for i := 0; i < fleetSize; i++ {
		go internal.InitOne(internal.NewClient(), newBus(i, coords), time.Tick(time.Second))
	}

	fmt.Println("Fleet initialized")

	fmt.Println("Waiting forever")
	forever := make(chan bool)
	<-forever
}
