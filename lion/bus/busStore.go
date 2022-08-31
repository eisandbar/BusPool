package bus

import (
	"github.com/eisandbar/BusPool/lion/types"
)

type BusStore interface {
	FindBus(point types.GeoPoint) int
}
