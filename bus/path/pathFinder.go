package path

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/golang/geo/s2"
)

var GraphHopper = "http://localhost:8989"

type PathFinder interface {
	GetPath([]s2.LatLng) ([]s2.LatLng, error)
}

// This path finder will add a point to the end of points needed to visit
type DumbPathFinder struct {
}

// Returns a polyline encoded array of coordinates
func (pf DumbPathFinder) GetPath(points []s2.LatLng) ([]s2.LatLng, error) {
	req := generateRequest(points)
	resp, err := http.Post(fmt.Sprintf("%s/route", GraphHopper), "application/json", req)
	if err != nil {
		return nil, errors.New("Failed to get directions from GraphHopper")
	}
	var res response
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, errors.New("Failed to read response body")
	}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return nil, errors.New("Failed to unmarshal json")
	}
	if len(res.Paths) < 1 {
		return nil, errors.New("Failed to find a path")
	}
	return getPoints(res.Paths[0].Points.Coordinates), nil
}

func generateRequest(points []s2.LatLng) io.Reader {
	coords := make([][]float64, len(points))
	for i, p := range points {
		coords[i] = []float64{p.Lng.Degrees(), p.Lat.Degrees()}
	}
	req := request{
		Points:        coords,
		Instructions:  false,
		Optimize:      true,
		PointsEncoded: false,
	}
	body, _ := json.Marshal(req)
	return bytes.NewBuffer(body)
}

func getPoints(coords [][]float64) []s2.LatLng {
	points := make([]s2.LatLng, len(coords))
	for i, coord := range coords {
		fmt.Println(coord)
		points[i] = s2.LatLngFromDegrees(coord[1], coord[0]) // coords come as [lon, lat]
	}
	return points
}

type request struct {
	Points        [][]float64 `json:"points"`
	Instructions  bool        `json:"instructions"`
	Optimize      bool        `json:"optimize"`
	PointsEncoded bool        `json:"points_encoded"`
}

type response struct {
	Paths []struct {
		Points struct {
			Coordinates [][]float64
		}
	}
}
