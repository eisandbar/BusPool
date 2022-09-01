package endpoints_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/eisandbar/BusPool/lion/bus"
	"github.com/eisandbar/BusPool/lion/endpoints"
	"github.com/eisandbar/BusPool/lion/types"
	"github.com/golang/geo/s2"
	"github.com/stretchr/testify/assert"
)

func TestRequestPost(t *testing.T) {
	id := 5
	point := types.GeoPoint{LatLng: s2.LatLngFromDegrees(15, 15)}
	pub := mockPublisher{}

	body, err := json.Marshal(point)
	assert.NoError(t, err)

	request, _ := http.NewRequest(http.MethodPost, "/requests", bytes.NewBuffer(body))
	response := httptest.NewRecorder()

	rs := endpoints.RequestServer{
		BusStore:   mockBusStore(bus.Bus{Id: id}),
		PathFinder: mockPathFinder{},
		Pub:        &pub,
	}
	rs.RequestPost(response, request)

	assert.Equal(t, 1, len(pub.calls))
	assert.Equal(t, 1, len(pub.paths))

	assert.Equal(t, id, pub.calls[0])
	assert.Equal(t, strconv.Itoa(id), pub.paths[0])
}

type mockBusStore bus.Bus

func (bs mockBusStore) FindBus(point types.GeoPoint) bus.Bus {
	return bus.Bus(bs)
}

func (bs mockBusStore) Store(bus bus.Bus) {

}

type mockPathFinder struct {
}

func (pf mockPathFinder) GetPath(bus bus.Bus, point types.GeoPoint) (string, error) {
	return strconv.Itoa(bus.Id), nil
}

type mockPublisher struct {
	calls []int
	paths []string
}

func (pub *mockPublisher) Publish(bus bus.Bus, path string) {
	pub.calls = append(pub.calls, bus.Id)
	pub.paths = append(pub.paths, path)
}
