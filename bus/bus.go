package main

import (
	"encoding/json"
	"math/rand"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/golang/geo/s2"
)

const fleetSize = 20

func newBus(id int) bus {
	X, Y, Z := rand.Float64(), rand.Float64(), rand.Float64()
	return bus{id, s2.PointFromCoords(X, Y, Z)}
}

type bus struct {
	Id int
	s2.Point
}

func (b bus) Report(client mqtt.Client) {
	body, err := json.Marshal(b)
	if err != nil {
		panic("Bad bus struct")
	}
	token := client.Publish("bus/positions", 0, false, body)
	token.Wait()
}
