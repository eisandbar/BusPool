package main

import (
	"context"
	"encoding/json"
	"log"
	"sync"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	. "github.com/eisandbar/BusPool/bus/typing"
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
	// Unmarshal data from mqtt message
	var bus Bus
	err := json.Unmarshal(msg.Payload(), &bus)
	if err != nil {
		log.Fatalf("Failed to unmarshal payload from mqtt topic, %s\n", msg.Payload())
	}

	// Transform data for later consumers
	eBus := elasticBus{
		ID:   bus.Id,
		Time: bus.Time.UnixMilli(),
		Location: geoPoint{
			Lat: bus.Location.Lat.Degrees(),
			Lon: bus.Location.Lng.Degrees(),
		},
	}
	body, err := json.Marshal(eBus)
	if err != nil {
		log.Fatalf("Failed to marshal kafkaBus, %+v\n", eBus)
	}

	// Adding data to kafka topic
	var wg sync.WaitGroup
	wg.Add(1)
	// Send raw data into bus-positions topic
	record := &kgo.Record{Topic: "bus-positions", Value: msg.Payload()}
	r.client.Produce(r.ctx, record, func(_ *kgo.Record, err error) {
		defer wg.Done()
		if err != nil {
			log.Fatalf("record had a produce error: %v\n", err)
		}
	})
	wg.Add(1)
	// Send transformed data into topic for elastic
	record = &kgo.Record{Topic: "bus-positions-elastic", Value: body}
	r.client.Produce(r.ctx, record, func(_ *kgo.Record, err error) {
		defer wg.Done()
		if err != nil {
			log.Fatalf("record had a produce error: %v\n", err)
		}
	})
	wg.Wait()
}

type geoPoint struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

type elasticBus struct {
	ID       int      `json:"id"`
	Time     int64    `json:"time"`
	Location geoPoint `json:"location"`
}
