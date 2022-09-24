package main

import (
	"encoding/json"
	"io"
	"log"
	"os"
)

func loadData(file string) [][]float64 {
	dataFile, err := os.Open(file)
	if err != nil {
		log.Fatalf("Failed to open data file, %s", err)
	}
	defer dataFile.Close()

	byteData, err := io.ReadAll(dataFile)
	if err != nil {
		log.Fatalf("Failed to read data file, %s", err)
	}

	var geoData geoJson
	err = json.Unmarshal(byteData, &geoData)
	if err != nil {
		log.Fatalf("Failed to unmarshal geojson, %s", err)
	}

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
