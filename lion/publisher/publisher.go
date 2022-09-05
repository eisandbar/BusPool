package publisher

import (
	"fmt"

	"github.com/eisandbar/BusPool/lion/bus"
)

type Publisher interface {
	Publish(bus.Bus, string)
}

type EmptyPublisher struct {
}

func (pub EmptyPublisher) Publish(bus bus.Bus, path string) {
	fmt.Println(bus, path)
}
