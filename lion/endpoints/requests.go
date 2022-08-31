package endpoints

import (
	"encoding/json"
	"net/http"

	"github.com/eisandbar/BusPool/lion/bus"
	"github.com/eisandbar/BusPool/lion/publisher"
	"github.com/eisandbar/BusPool/lion/types"
)

type RequestServer struct {
	BusStore bus.BusStore
	Pub      publisher.Publisher
}

// Handler for client requests
func (rs RequestServer) RequestPost(w http.ResponseWriter, r *http.Request) {
	var point types.GeoPoint

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&point)
	if err != nil {

	}
	id := rs.BusStore.FindBus(point)
	rs.Pub.Publish(point, id)
}
