package path_test

import (
	"testing"

	"github.com/eisandbar/BusPool/lion/bus"
	"github.com/eisandbar/BusPool/lion/path"
	"github.com/eisandbar/BusPool/lion/types"
	"github.com/golang/geo/s2"
	"github.com/stretchr/testify/assert"
)

// This is an integration test with GraphHopper
func TestGetPath(t *testing.T) {
	pf := path.DumbPathFinder{}
	testBus := bus.Bus{
		Points: []s2.LatLng{s2.LatLngFromDegrees(52.52, 13.37)},
	}
	point := types.GeoPoint{LatLng: s2.LatLngFromDegrees(52.5, 13.4)}

	// The \\ in the expected string is actually just a single \, but it had to be escaped
	expectedString := "gvp_I{nrpAnI@TFBsd@vCgAf@CTNTh@X^JFVBTCXMLUX_ATc@XKdFXvGv@xCFpLzAX@lBG`AK\\K^YbJsLhPgU`UyZzFiGT_@BY?O~@?jDKdB]X_DTsFAyEKqDa@wJGg@m@gD}AUaAHw@TQsAEIQGqB}A_J_EFw@Ci@aCcNLKNApA?FBBFZtBdBBPDnA~Az@@DGBKDw@BG"
	res, err := pf.GetPath(testBus, point)
	assert.NoError(t, err)
	assert.Equal(t, expectedString, res)
}
