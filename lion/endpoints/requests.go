package endpoints

import (
	"encoding/json"
	"net/http"

	"github.com/eisandbar/BusPool/lion/bus"
	"github.com/eisandbar/BusPool/lion/publisher"
	. "github.com/eisandbar/BusPool/lion/typing"
	"github.com/golang/geo/s2"
)

type RequestServer struct {
	BusStore bus.BusStore
	Pub      publisher.Publisher
}

// Handler for client requests
func (rs RequestServer) RequestPost(w http.ResponseWriter, r *http.Request) {
	var req Request

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if !nonEmptyRequest(req) {
		http.Error(w, "Missing fields", http.StatusBadRequest)
		return
	}

	instruction := ReqToInstruction(req)

	bus, err := rs.BusStore.FindBus(instruction.Client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	rs.Pub.Publish(bus, instruction)
}

func nonEmptyRequest(req Request) bool {
	return req.Client != nil && len(req.Client) == 2 && req.Dest != nil && len(req.Dest) == 2
}

func ReqToInstruction(req Request) Instruction {
	return Instruction{
		Client: s2.LatLngFromDegrees(req.Client[0], req.Client[1]),
		Dest:   s2.LatLngFromDegrees(req.Dest[0], req.Dest[1]),
	}
}
