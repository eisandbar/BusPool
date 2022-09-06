package bus

import (
	"errors"
	"sync"

	. "github.com/eisandbar/BusPool/lion/typing"
	"github.com/golang/geo/s1"
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

	found := false

	res, dist := Bus{}, s1.Angle(2)
	for _, bus := range bs.bus {
		if client.Distance(bus.Location).Abs() <= dist && bus.Occupancy < bus.Capacity {
			found = true
			res = bus
			dist = client.Distance(bus.Location).Abs()
		}
	}

	if !found {
		return Bus{}, errors.New("No buses available")
	}

	res.Occupancy++
	bs.bus[res.Id] = res
	return res, nil
}
