package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"
	// "github.com/eisandbar/BusPool/lion/typing"
)

const (
	file   = "data/export.geojson"       // file with bus stop data
	server = "http://lion:3333/requests" // Request server
	N      = 600                         // How often to generate a request in milliseconds
)

func main() {
	coords := loadData(file)
	// generate request every N milliseconds
	ticker := time.Tick(N * time.Millisecond)
	for {
		_, more := <-ticker
		if more {
			req := generateRequest(coords)
			resp, err := http.Post(server, "application/json", req)
			if err != nil {
				log.Fatalf("Failed to send request, %s", err)
			}
			buf := new(strings.Builder)
			io.Copy(buf, resp.Body)
			fmt.Println(resp.StatusCode, buf.String())
		} else {
			continue
		}
	}
}

type Request struct {
	Client []float64 `json:"client"`
	Dest   []float64 `json:"dest"`
}

func generateRequest(coords [][]float64) io.Reader {
	coordsLength := len(coords)

	// Randomly pick 2 points
	rand.Seed(time.Now().UnixMilli())
	req := Request{
		Client: coords[rand.Intn(coordsLength)],
		Dest:   coords[rand.Intn(coordsLength)],
	}
	log.Println(req.Client)
	body, err := json.Marshal(req)
	if err != nil {
		log.Fatalf("Failed to generate request, %s", err)
	}
	return bytes.NewBuffer(body)
}
