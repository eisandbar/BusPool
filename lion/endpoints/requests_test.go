package endpoints_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/eisandbar/BusPool/lion/endpoints"
	"github.com/eisandbar/BusPool/lion/typing"
	. "github.com/eisandbar/BusPool/lion/typing"
	"github.com/golang/geo/s2"
	"github.com/stretchr/testify/assert"
)

func TestRequestPost(t *testing.T) {
	id := 5
	point := typing.Request{
		Client: []float64{15, 15},
		Dest:   []float64{14, 16},
	}
	pub := mockPublisher{}

	body, err := json.Marshal(point)
	assert.NoError(t, err)

	request, _ := http.NewRequest(http.MethodPost, "/requests", bytes.NewBuffer(body))
	response := httptest.NewRecorder()

	rs := endpoints.RequestServer{
		BusStore: mockBusStore(Bus{Id: id}),
		Pub:      &pub,
	}
	rs.RequestPost(response, request)
	assert.Equal(t, http.StatusOK, response.Code)

	assert.Equal(t, 1, len(pub.calls))
	assert.Equal(t, 1, len(pub.paths))

	assert.Equal(t, id, pub.calls[0])
	assert.Equal(t, endpoints.ReqToInstruction(point), pub.paths[0])
}

func TestRequestPostFail(t *testing.T) {

	rs := endpoints.RequestServer{
		BusStore: mockBusStore(Bus{Id: 0}),
		Pub:      &mockPublisher{},
	}
	t.Run("Bad request", func(t *testing.T) {
		body, err := json.Marshal(struct{ BadField int }{2})
		assert.NoError(t, err)
		request, _ := http.NewRequest(http.MethodPost, "/requests", bytes.NewBuffer(body))
		response := httptest.NewRecorder()
		rs.RequestPost(response, request)
		assert.Equal(t, http.StatusBadRequest, response.Code)
	})
	t.Run("Missing fields in request", func(t *testing.T) {
		body, err := json.Marshal(Request{Client: []float64{12, 12}})
		assert.NoError(t, err)
		request, _ := http.NewRequest(http.MethodPost, "/requests", bytes.NewBuffer(body))
		response := httptest.NewRecorder()
		rs.RequestPost(response, request)
		assert.Equal(t, http.StatusBadRequest, response.Code)
	})
}

type mockBusStore Bus

func (bs mockBusStore) FindBus(point s2.LatLng) (Bus, error) {
	if point.Lat == -1 {
		return Bus{}, errors.New("No buses available")
	}
	return Bus(bs), nil
}

func (bs mockBusStore) Store(bus Bus) {

}

type mockPublisher struct {
	calls []int
	paths []Instruction
}

func (pub *mockPublisher) Publish(bus Bus, inst Instruction) {
	pub.calls = append(pub.calls, bus.Id)
	pub.paths = append(pub.paths, inst)
}
