package bus

import (
	"errors"
	"sync"

	"github.com/eisandbar/BusPool/lion/types"
	"github.com/golang/geo/s2"
)

type BusStore interface {
	FindBus(point types.GeoPoint) (Bus, error)
	Store(Bus)
}

type Bus struct {
	Id       int
	Location types.GeoPoint
	Points   []s2.LatLng
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
func (bs *MemoryBusStore) FindBus(point types.GeoPoint) (Bus, error) {
	bs.RLock()
	defer bs.RUnlock()
	if len(bs.bus) == 0 {
		return Bus{}, errors.New("No buses available")
	}
	res, dist := bs.bus[0], point.Distance(bs.bus[0].Location.LatLng).Abs()
	for _, bus := range bs.bus {
		if point.Distance(bus.Location.LatLng).Abs() < dist {
			res = bus
			dist = point.Distance(bus.Location.LatLng).Abs()
		}
	}
	return res, nil
}
