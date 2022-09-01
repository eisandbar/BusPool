package path

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/eisandbar/BusPool/lion/bus"
	"github.com/eisandbar/BusPool/lion/types"
)

var GraphHopper = "http://localhost:8989"

type PathFinder interface {
	GetPath(bus.Bus, types.GeoPoint) (string, error)
}

// This path finder will add a point to the end of points needed to visit
type DumbPathFinder struct {
}

// Returns a polyline encoded array of coordinates
func (pf DumbPathFinder) GetPath(bus bus.Bus, point types.GeoPoint) (string, error) {
	req := generateRequest(bus, point)
	resp, err := http.Post(fmt.Sprintf("%s/route", GraphHopper), "application/json", req)
	if err != nil {
		return "", errors.New("Failed to get directions from GraphHopper")
	}
	var res response
	body, err := io.ReadAll(resp.Body)

	fmt.Printf("%+v\n %s", resp, body)
	if err != nil {
		return "", errors.New("Failed to read response body")
	}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return "", errors.New("Failed to unmarshal json")
	}
	if len(res.Paths) < 1 {
		return "", errors.New("Failed to find a path")
	}
	return res.Paths[0].Points, nil
}

func generateRequest(bus bus.Bus, point types.GeoPoint) io.Reader {
	points := make([][]float64, 0, len(bus.Points)+1)
	for _, p := range bus.Points {
		points = append(points, []float64{p.Lng.Degrees(), p.Lat.Degrees()})
	}
	points = append(points, []float64{point.Lng.Degrees(), point.Lat.Degrees()})
	req := request{
		Points:       points,
		Instructions: false,
		Optimize:     false,
	}
	body, _ := json.Marshal(req)
	return bytes.NewBuffer(body)
}

type request struct {
	Points       [][]float64 `json:"points"`
	Instructions bool        `json:"instructions"`
	Optimize     bool        `json:"optimize"`
}

type response struct {
	Paths []struct {
		Points string
	}
}
