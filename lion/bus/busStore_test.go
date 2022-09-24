package bus_test

import (
	"testing"

	"github.com/eisandbar/BusPool/lion/bus"
	. "github.com/eisandbar/BusPool/lion/typing"
	"github.com/golang/geo/s2"
	"github.com/stretchr/testify/assert"
)

func TestFindBus(t *testing.T) {
	bs := bus.NewMemoryBusStore()

	_, err := bs.FindBus(s2.LatLng{})
	assert.Error(t, err)

	buses := []Bus{
		{Id: 1, Location: s2.LatLngFromDegrees(19, 13), Capacity: 6},
		{Id: 2, Location: s2.LatLngFromDegrees(21, 15), Capacity: 6},
		{Id: 3, Location: s2.LatLngFromDegrees(22, 22), Capacity: 6},
	}
	for _, bus := range buses {
		bs.Store(bus)
	}

	testData := []struct {
		point s2.LatLng
		id    int
	}{
		{point: s2.LatLngFromDegrees(22, 16), id: 2},
		{point: s2.LatLngFromDegrees(22, 13), id: 2},
		{point: s2.LatLngFromDegrees(15, 13), id: 1},
		{point: s2.LatLngFromDegrees(21, 21), id: 3},
	}

	for _, tt := range testData {
		bus, err := bs.FindBus(tt.point)
		assert.NoError(t, err)
		assert.Equal(t, tt.id, bus.Id)
	}
}
