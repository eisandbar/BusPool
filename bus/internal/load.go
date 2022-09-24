package internal

import (
	"encoding/json"
	"io"
	"log"
	"os"
)

const (
	file = "data/export.geojson"
)

func LoadData() [][]float64 {
	dataFile, err := os.Open(file)
	defer dataFile.Close()
	if err != nil {
		log.Fatalf("Failed to open data file, %s", err)
	}

	byteData, err := io.ReadAll(dataFile)
	if err != nil {
		log.Fatalf("Failed to read data file, %s", err)
	}

	var geoData geoJson
	json.Unmarshal(byteData, &geoData)

	coordinates := make([][]float64, len(geoData.Features))
	for i, feature := range geoData.Features {
		// The json has it in [lon, lat]
		coordinates[i] = []float64{feature.Geometry.Coordinates[1], feature.Geometry.Coordinates[0]}
	}

	return coordinates
}

type geoJson struct {
	Features []struct {
		Geometry struct {
			Coordinates []float64 `json:"coordinates"`
		} `json:"geometry"`
	} `json:"features"`
}
