package main

import (
	"testing"
	"time"

	. "github.com/eisandbar/BusPool/bus/typing"
	"github.com/golang/geo/s2"
	"github.com/stretchr/testify/assert"
)

func TestMove(t *testing.T) {
	start := s2.LatLngFromDegrees(0, 0)
	path := []s2.LatLng{
		s2.LatLngFromDegrees(0, 0.01),
		s2.LatLngFromDegrees(0, 0.02),
		s2.LatLngFromDegrees(0, 0.03),
	}
	client := []s2.LatLng{
		s2.LatLngFromDegrees(0, 0.031),
	}
	dest := []s2.LatLng{
		s2.LatLngFromDegrees(0, 0.041),
	}
	testBus := bus{
		Bus: Bus{
			Location:     start,
			Path:         path,
			Clients:      client,
			Destinations: dest,
			Occupancy:    1,
		},
		Path: mockPathFinder{},
	}
	// Start
	assert.Equal(t, start, testBus.Bus.Location)
	assert.Equal(t, client, testBus.Bus.Clients)
	assert.Equal(t, 1, testBus.Bus.Occupancy)

	// Move 1
	testBus.Move()
	assert.Equal(t, path[0], testBus.Bus.Location)
	assert.Equal(t, client, testBus.Bus.Clients)

	// Move 2
	testBus.Move()
	assert.Equal(t, path[1], testBus.Bus.Location)
	assert.Equal(t, client, testBus.Bus.Clients)

	// Pickup client
	testBus.Move()
	assert.Equal(t, path[2], testBus.Bus.Location)
	assert.Equal(t, []s2.LatLng{}, testBus.Bus.Clients)
	assert.Equal(t, []s2.LatLng{}, testBus.Bus.Path) // No path left
	time.Sleep(10 * time.Millisecond)
	assert.Equal(t, append([]s2.LatLng{path[2]}, dest...), testBus.Bus.Path) // After rerouting

	// 1 Move wasted on pickup
	testBus.Move()
	assert.Equal(t, path[2], testBus.Bus.Location)
	assert.Equal(t, 1, testBus.Bus.Occupancy)

	// Reach destination
	testBus.Move()
	assert.Equal(t, dest[0], testBus.Bus.Location)
	assert.Equal(t, []s2.LatLng{}, testBus.Bus.Destinations)
	assert.Equal(t, 0, testBus.Bus.Occupancy)
	assert.Equal(t, []s2.LatLng{}, testBus.Bus.Path) // No path left

}

type mockPathFinder struct {
}

func (pf mockPathFinder) GetPath(points []s2.LatLng) ([]s2.LatLng, error) {
	return points, nil
}
