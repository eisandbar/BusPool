package bus_test

import (
	"testing"

	"github.com/eisandbar/BusPool/lion/bus"
	"github.com/eisandbar/BusPool/lion/types"
	"github.com/golang/geo/s2"
	"github.com/stretchr/testify/assert"
)

func TestFindBus(t *testing.T) {
	bs := bus.NewMemoryBusStore()
	buses := []bus.Bus{
		{Id: 1, Location: types.GeoPoint{LatLng: s2.LatLngFromDegrees(19, 13)}},
		{Id: 2, Location: types.GeoPoint{LatLng: s2.LatLngFromDegrees(21, 15)}},
		{Id: 3, Location: types.GeoPoint{LatLng: s2.LatLngFromDegrees(22, 22)}},
	}
	for _, bus := range buses {
		bs.Store(bus)
	}

	testData := []struct {
		point types.GeoPoint
		id    int
	}{
		{point: types.GeoPoint{LatLng: s2.LatLngFromDegrees(22, 16)}, id: 2},
		{point: types.GeoPoint{LatLng: s2.LatLngFromDegrees(22, 13)}, id: 2},
		{point: types.GeoPoint{LatLng: s2.LatLngFromDegrees(15, 13)}, id: 1},
		{point: types.GeoPoint{LatLng: s2.LatLngFromDegrees(21, 21)}, id: 3},
	}

	for _, tt := range testData {
		assert.Equal(t, tt.id, bs.FindBus(tt.point).Id)
	}
}
