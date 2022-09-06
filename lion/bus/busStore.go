package bus

import (
	"errors"
	"sync"

	. "github.com/eisandbar/BusPool/lion/typing"
	"github.com/golang/geo/s2"
)

type BusStore interface {
	FindBus(s2.LatLng) (Bus, error)
	Store(Bus)
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
func (bs *MemoryBusStore) FindBus(client s2.LatLng) (Bus, error) {
	bs.RLock()
	defer bs.RUnlock()
	if len(bs.bus) == 0 {
		return Bus{}, errors.New("No buses available")
	}
	res, dist := bs.bus[0], client.Distance(bs.bus[0].Location).Abs()
	for _, bus := range bs.bus {
		if client.Distance(bus.Location).Abs() < dist {
			res = bus
			dist = client.Distance(bus.Location).Abs()
		}
	}
	return res, nil
}
