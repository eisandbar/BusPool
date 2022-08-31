package publisher

import "github.com/eisandbar/BusPool/lion/types"

type Publisher interface {
	Publish(point types.GeoPoint, id int)
}
