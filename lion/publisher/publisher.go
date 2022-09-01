package publisher

import (
	"github.com/eisandbar/BusPool/lion/bus"
)

type Publisher interface {
	Publish(bus.Bus, string)
}
