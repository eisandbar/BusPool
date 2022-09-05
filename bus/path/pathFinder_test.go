package path_test

import (
	"testing"

	"github.com/eisandbar/BusPool/bus/path"
	"github.com/golang/geo/s2"
	"github.com/stretchr/testify/assert"
)

// This is an integration test with GraphHopper
func TestGetPath(t *testing.T) {
	points := []s2.LatLng{
		s2.LatLngFromDegrees(52.52, 13.37),
		s2.LatLngFromDegrees(52.519, 13.37),
	}

	want := []s2.LatLng{
		s2.LatLngFromDegrees(52.51957, 13.370862),
		s2.LatLngFromDegrees(52.518999, 13.370862),
	}
	pf := path.DumbPathFinder{}
	res, err := pf.GetPath(points)
	assert.NoError(t, err)
	assert.Equal(t, want, res)
}
