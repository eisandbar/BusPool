package subscriber

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/eisandbar/BusPool/lion/bus"
	"github.com/twmb/franz-go/pkg/kgo"
)

var seeds = []string{"localhost:9092"}

type Subscriber struct {
	client *kgo.Client
	bs     bus.BusStore
}

// Returns new subscriber
func NewSubscriber(bs bus.BusStore, topic, group string) (Subscriber, error) {
	sub := Subscriber{}

	// Creating new client
	cl, err := kgo.NewClient(
		kgo.SeedBrokers(seeds...),
		kgo.ConsumerGroup(group),
		kgo.ConsumeTopics(topic),
	)

	if err != nil {
		return sub, err
	}

	sub.client = cl
	sub.bs = bs

	return sub, nil
}

// Subscribes to kafka topic and populates bus store
func (sub Subscriber) Subscribe() {
	fmt.Println("Subscribing to topic")
	ctx := context.Background()
	for {
		fetches := sub.client.PollFetches(ctx)
		if errs := fetches.Errors(); len(errs) > 0 {
			panic(fmt.Sprint(errs))
		}

		// Update bus positions in bus store
		fetches.EachPartition(func(p kgo.FetchTopicPartition) {
			p.EachRecord(func(record *kgo.Record) {
				bus := bus.Bus{}
				err := json.Unmarshal(record.Value, &bus)
				if err != nil || bus.Location.Lng == 0 {
					panic(fmt.Sprintf("%s, %+v", err, bus))
				}
				sub.bs.Store(bus)
			})
		})
	}
}
