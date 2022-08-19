package main

import (
	"fmt"
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
	id int
	s2.Point
}

func (b bus) Id() string {
	return fmt.Sprintf("bus %v", b.id)
}

func (b bus) Report(client mqtt.Client) {
	token := client.Publish("bus/positions", 0, false, fmt.Sprintf("id: %v, coord: %+v", b.id, b.Vector))
	token.Wait()
}
