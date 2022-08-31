package bus

import (
	"fmt"
	"sync"

	"github.com/eisandbar/BusPool/lion/types"
)

type BusStore interface {
	FindBus(point types.GeoPoint) Bus
	Store(Bus)
}

type Bus struct {
	Id       int
	Location types.GeoPoint
}

type MemoryBusStore struct {
	bus map[int]Bus
	sync.RWMutex
}

func NewMemoryBusStore() *MemoryBusStore {
	bs := MemoryBusStore{}
	bs.bus = make(map[int]Bus)
	return &bs
}

// Store bus data
func (bs *MemoryBusStore) Store(bus Bus) {
	bs.Lock()
	defer bs.Unlock()
	bs.bus[bus.Id] = bus
}

// Find bus goes over the list of buses and finds the nearest one
func (bs *MemoryBusStore) FindBus(point types.GeoPoint) Bus {
	bs.RLock()
	defer bs.RUnlock()
	res, dist := bs.bus[0], point.Distance(bs.bus[0].Location.LatLng).Abs()
	for _, bus := range bs.bus {
		fmt.Println(dist, point.Distance(bus.Location.LatLng).Abs())
		if point.Distance(bus.Location.LatLng).Abs() < dist {
			res = bus
			dist = point.Distance(bus.Location.LatLng).Abs()
		}
	}
	return res
}
