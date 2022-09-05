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
	"github.com/golang/geo/s2"
)

const fleetSize = 20     // number of buses to initialize
const errAngle = 0.00015 // error allowance in coordinates

func newBus(id int) *bus {
	lat := 52.35 + rand.Float64()*0.3
	lng := 13.1 + rand.Float64()*0.6
	location := s2.LatLngFromDegrees(lat, lng)
	return &bus{
		Bus{
			Id:        id,
			Capacity:  6,
			Occupancy: 0,
			Location:  location,
			Path:      []s2.LatLng{location},
		},
		path.DumbPathFinder{},
	}
}

type Bus struct {
	Id           int
	Time         time.Time
	Capacity     int
	Occupancy    int         // Number of passengers including those not yet picked up
	Location     s2.LatLng   // Current location
	Clients      []s2.LatLng // Clients that still need to be picked up
	Destinations []s2.LatLng // Client drop-off locations
	Path         []s2.LatLng
	sync.RWMutex
}

type Request struct {
	Client s2.LatLng
	Dest   s2.LatLng
}

type bus struct {
	Bus
	path.PathFinder
}

func (b *bus) Report(client mqtt.Client) {
	b.RLock()
	b.Time = time.Now()
	body, err := json.Marshal(b)
	b.RUnlock()

	if err != nil {
		panic("Bad bus struct")
	}
	token := client.Publish("bus/positions", 0, false, body)
	token.Wait()
}

func (b *bus) Subscribe(client mqtt.Client) {
	b.RLock()
	topic := fmt.Sprintf("bus/requests/%d", b.Id)
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
	if len(b.Path) == 0 {
		return
	}

	// Move to next point
	b.Location = b.Path[0]

	// Pickup
	for i := len(b.Clients) - 1; i >= 0; i-- {
		if b.Location.Distance(b.Clients[i]) < errAngle {
			b.Clients = append(b.Clients[0:i], b.Clients[i+1:]...)
		}
	}

	// Drop-off
	for i := len(b.Destinations) - 1; i >= 0; i-- {
		if b.Location.Distance(b.Destinations[i]) < errAngle {
			b.Destinations = append(b.Destinations[0:i], b.Destinations[i+1:]...)
			b.Occupancy--
		}
	}

	// Update path
	b.Path = b.Path[1:]
	if len(b.Path) == 0 {
		go b.reroute()
	}
}

func (b *bus) requestHandler(client mqtt.Client, msg mqtt.Message) {
	b.Lock()
	defer b.Unlock()
	var req Request
	err := json.Unmarshal(msg.Payload(), &req)
	if err != nil {
		log.Println("Failed to unmarshal request for bus", b.Id)
	}
	b.Clients = append(b.Clients, req.Client)
	b.Destinations = append(b.Destinations, req.Dest)
	b.Occupancy++
	go b.reroute()
}

func (b *bus) reroute() {
	b.Lock()
	defer b.Unlock()
	points := make([]s2.LatLng, 0)
	if len(b.Clients) > 0 {
		points = append([]s2.LatLng{b.Location}, b.Clients...)
	} else if len(b.Destinations) > 0 {
		points = append([]s2.LatLng{b.Location}, b.Destinations...)
	}
	path, err := b.GetPath(points)
	if err != nil {
		log.Println("Error finding path", err)
		return
	}
	b.Path = path
}
