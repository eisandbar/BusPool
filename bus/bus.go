package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/eisandbar/BusPool/bus/path"
	. "github.com/eisandbar/BusPool/bus/typing"
	"github.com/golang/geo/s2"
)

const fleetSize = 200    // number of buses to initialize
const errAngle = 0.00015 // error allowance in coordinates

func newBus(id int, coords [][]float64) *bus {
	rand.Seed(time.Now().UnixMilli())
	n := rand.Intn(len(coords))
	lat := coords[n][0]
	lng := coords[n][1]
	location := s2.LatLngFromDegrees(lat, lng)
	return &bus{
		Bus: Bus{
			Id:        id,
			Capacity:  6,
			Occupancy: 0,
			Location:  location,
			Path:      []s2.LatLng{location},
		},
		Path: path.DumbPathFinder{},
	}
}

type Request struct {
	Client s2.LatLng
	Dest   s2.LatLng
}

type bus struct {
	sync.RWMutex
	Bus  Bus
	Path path.PathFinder
}

func (b *bus) Report(client mqtt.Client) {
	b.RLock()
	b.Bus.Time = time.Now()
	body, err := json.Marshal(b.Bus)
	b.RUnlock()

	if err != nil {
		panic("Bad bus struct")
	}
	token := client.Publish("bus/positions", 0, false, body)
	token.Wait()
}

func (b *bus) Subscribe(client mqtt.Client) {
	b.RLock()
	topic := fmt.Sprintf("bus/requests/%d", b.Bus.Id)
	b.RUnlock()

	if token := client.Subscribe(topic, 0, b.requestHandler); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}
}

func (b *bus) Move() {
	b.Lock()
	defer b.Unlock()

	// Nowhere we need to go
	if len(b.Bus.Path) == 0 {
		return
	}

	// Move to next point
	b.Bus.Location = b.Bus.Path[0]

	// Pickup
	for i := len(b.Bus.Clients) - 1; i >= 0; i-- {
		if b.Bus.Location.Distance(b.Bus.Clients[i]) < errAngle {
			b.Bus.Clients = append(b.Bus.Clients[0:i], b.Bus.Clients[i+1:]...)
		}
	}

	// Drop-off
	for i := len(b.Bus.Destinations) - 1; i >= 0; i-- {
		if b.Bus.Location.Distance(b.Bus.Destinations[i]) < errAngle {
			b.Bus.Destinations = append(b.Bus.Destinations[0:i], b.Bus.Destinations[i+1:]...)
			b.Bus.Occupancy--
		}
	}

	// Update path
	b.Bus.Path = b.Bus.Path[1:]
	if len(b.Bus.Path) == 0 {
		go b.reroute()
	}
}

func (b *bus) requestHandler(client mqtt.Client, msg mqtt.Message) {
	b.Lock()
	defer b.Unlock()
	log.Println("Received instructions for bus %d", b.Bus.Id)
	var req Request
	err := json.Unmarshal(msg.Payload(), &req)
	if err != nil {
		log.Println("Failed to unmarshal request for bus", b.Bus.Id)
	}
	b.Bus.Clients = append(b.Bus.Clients, req.Client)
	b.Bus.Destinations = append(b.Bus.Destinations, req.Dest)
	b.Bus.Occupancy++
	go b.reroute()
}

func (b *bus) reroute() {
	b.Lock()
	defer b.Unlock()
	points := make([]s2.LatLng, 0)
	if len(b.Bus.Clients) > 0 {
		points = append([]s2.LatLng{b.Bus.Location}, b.Bus.Clients...)
	} else if len(b.Bus.Destinations) > 0 {
		points = append([]s2.LatLng{b.Bus.Location}, b.Bus.Destinations...)
	}
	path, err := b.Path.GetPath(points)
	if err != nil {
		log.Println("Error finding path", err)
		return
	}
	b.Bus.Path = path
}
