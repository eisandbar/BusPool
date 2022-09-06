package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/eisandbar/BusPool/lion/bus"
	"github.com/eisandbar/BusPool/lion/endpoints"
	"github.com/eisandbar/BusPool/lion/publisher"
	"github.com/eisandbar/BusPool/lion/subscriber"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

const (
	KAFKA_TOPIC = "bus-positions"
	KAFKA_GROUP = "lion"
	CONN_PORT   = "3333"
)

func main() {
	// Creating new bus store
	bs := bus.NewMemoryBusStore()

	// Creating new subscriber
	sub, err := subscriber.NewSubscriber(bs, KAFKA_TOPIC, KAFKA_GROUP)
	if err != nil {
		panic(err)
	}
	// Start listening on topic in new goroutine
	go sub.Subscribe(KAFKA_TOPIC)

	// Create Request Server
	rs := endpoints.RequestServer{
		BusStore: bs,
		Pub:      publisher.NewMQTTPublisher(),
	}

	// Start http server
	router := mux.NewRouter()
	router.HandleFunc("/requests", rs.RequestPost).Methods("POST")

	handler := cors.Default().Handler(router)

	fmt.Println("Listening on port:", CONN_PORT)
	log.Fatal(http.ListenAndServe(":"+CONN_PORT, handler))
}
