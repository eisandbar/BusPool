package main

import (
	"testing"
	"time"

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
		Bus{
			Location:     start,
			Path:         path,
			Clients:      client,
			Destinations: dest,
		},
		mockPathFinder{},
	}
	// Start
	assert.Equal(t, start, testBus.Location)
	assert.Equal(t, client, testBus.Clients)

	// Move 1
	testBus.Move()
	assert.Equal(t, path[0], testBus.Location)
	assert.Equal(t, client, testBus.Clients)

	// Move 2
	testBus.Move()
	assert.Equal(t, path[1], testBus.Location)
	assert.Equal(t, client, testBus.Clients)

	// Pickup client
	testBus.Move()
	assert.Equal(t, path[2], testBus.Location)
	assert.Equal(t, []s2.LatLng{}, testBus.Clients)
	assert.Equal(t, []s2.LatLng{}, testBus.Path) // No path left
	time.Sleep(10 * time.Millisecond)
	assert.Equal(t, append([]s2.LatLng{path[2]}, dest...), testBus.Path) // After rerouting

	// 1 Move wasted on pickup
	testBus.Move()
	assert.Equal(t, path[2], testBus.Location)

	// Reach destination
	testBus.Move()
	assert.Equal(t, dest[0], testBus.Location)
	assert.Equal(t, []s2.LatLng{}, testBus.Destinations)
	assert.Equal(t, []s2.LatLng{}, testBus.Path) // No path left

}

type mockPathFinder struct {
}

func (pf mockPathFinder) GetPath(points []s2.LatLng) ([]s2.LatLng, error) {
	return points, nil
}
