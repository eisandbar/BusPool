package main

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/golang/geo/s2"
	"github.com/twmb/franz-go/pkg/kadm"
	"github.com/twmb/franz-go/pkg/kgo"
)

var seeds = []string{"localhost:9092"}

type rhino struct {
	client *kgo.Client
	ctx    context.Context
}

func (r *rhino) Init() {
	cl, err := kgo.NewClient(
		kgo.SeedBrokers(seeds...),
	)
	if err != nil {
		panic(err)
	}

	r.client = cl
	r.ctx = context.Background()
}

func (r *rhino) Close() {
	r.client.Close()
}

func (r rhino) InitTopic(topic string) {
	admClient := kadm.NewClient(r.client)

	// Check that topic doesn't exist
	topicDetails, err := admClient.ListTopics(r.ctx)
	if err != nil {
		panic(err)
	}
	if topicDetails.Has(topic) { // Topic already exists
		return
	}

	// Create topic
	_, err = admClient.CreateTopics(r.ctx, 3, 1, nil, topic)
	if err != nil {
		panic(err)
	}
}

func (r rhino) positionHandler(client mqtt.Client, msg mqtt.Message) {
	// Unmarshal bus data from mqtt message
	var mBus mqttBus
	err := json.Unmarshal(msg.Payload(), &mBus)
	if err != nil {
		panic("Bad payload in mqtt topic")
	}
	latlng := s2.LatLngFromPoint(mBus.Point)

	// Transform data for later consumers
	kBus := kafkaBus{
		Id: mBus.Id,
		Location: GeoPoint{
			Lat: latlng.Lat.Degrees(), Lon: latlng.Lng.Degrees(), Type: "geo_point",
		},
	}
	body, err := json.Marshal(kBus)

	// Adding data to kafka topic
	var wg sync.WaitGroup
	wg.Add(1)
	record := &kgo.Record{Topic: "bus-positions", Value: body}
	r.client.Produce(r.ctx, record, func(_ *kgo.Record, err error) {
		defer wg.Done()
		if err != nil {
			fmt.Printf("record had a produce error: %v\n", err)
		}

	})
	wg.Wait()
}

type mqttBus struct {
	Id int
	s2.Point
}

type kafkaBus struct {
	Id       int
	Location GeoPoint `json:"location"`
}

type GeoPoint struct {
	Type string  `json:"type"`
	Lat  float64 `json:"lat"`
	Lon  float64 `json:"lon"`
}
