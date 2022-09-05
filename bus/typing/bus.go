package typing

import (
	"sync"
	"time"

	"github.com/golang/geo/s2"
)

type Bus struct {
	Id           int
	Time         time.Time
	Capacity     int
	Occupancy    int         // Number of passengers including those not yet picked up
	Location     s2.LatLng   // Current location
	Clients      []s2.LatLng // Clients that still need to be picked up
	Destinations []s2.LatLng // Client drop-off locations
	Path         []s2.LatLng
	sync.RWMutex
}
