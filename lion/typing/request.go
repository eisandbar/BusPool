package typing

import "github.com/golang/geo/s2"

type Request struct {
	Client []float64 `json:"client"` // [lat, long]
	Dest   []float64 `json:"dest"`   // [lat, long]
}

type Instruction struct {
	Client s2.LatLng
	Dest   s2.LatLng
}
