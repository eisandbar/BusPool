package bus

import (
	"fmt"
	"sync"

	"github.com/eisandbar/BusPool/lion/types"
)

type BusStore interface {
	FindBus(point types.GeoPoint) int
}

type Bus struct {
	Id       int
	Location types.GeoPoint
}

type MemoryBusStore struct {
	bus []Bus
	sync.RWMutex
}

// Init bus slice mostly for testing
func (bs *MemoryBusStore) Init(bus []Bus) {
	bs.bus = bus
}

// Find bus goes over the list of buses and finds the nearest one
func (bs *MemoryBusStore) FindBus(point types.GeoPoint) int {
	bs.RLock()
	defer bs.RUnlock()
	id, dist := bs.bus[0].Id, point.Distance(bs.bus[0].Location.LatLng).Abs().Degrees()
	for _, bus := range bs.bus {
		fmt.Println(dist, point.Distance(bus.Location.LatLng).Abs().Degrees())
		if point.Distance(bus.Location.LatLng).Abs().Degrees() < dist {
			id = bus.Id
			dist = point.Distance(bus.Location.LatLng).Abs().Degrees()
		}
	}
	return id
}
