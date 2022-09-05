package internal_test

import (
	"sync"
	"testing"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/eisandbar/BusPool/bus/internal"
	"github.com/stretchr/testify/assert"
)

// Tests that InitOne will call r.Report() every time something is sent to channel tick
func TestInitOne(t *testing.T) {
	client := mqtt.NewClient(mqtt.NewClientOptions())
	reporter := mockReportable{s: "Test"}
	testChan := make(chan time.Time)

	go internal.InitOne(client, &reporter, testChan)

	var want []string
	for i := 0; i < 5; i++ { // 5 is a random number
		time.Sleep(time.Millisecond) // Wait for InitOne
		reporter.Lock()
		assert.Equal(t, []string{"Test"}, reporter.subscribes)
		assert.Equal(t, want, reporter.reports)
		assert.Equal(t, want, reporter.moves)
		reporter.Unlock()
		testChan <- time.Now()
		want = append(want, "Test")
	}
	close(testChan)
}

type mockReportable struct {
	s          string
	reports    []string
	subscribes []string
	moves      []string
	sync.Mutex
}

func (r *mockReportable) Id() string {
	return r.s
}

func (r *mockReportable) Report(mqtt.Client) {
	r.Lock()
	defer r.Unlock()
	r.reports = append(r.reports, r.s)
}

func (r *mockReportable) Subscribe(mqtt.Client) {
	r.Lock()
	defer r.Unlock()
	r.subscribes = append(r.subscribes, r.s)
}

func (r *mockReportable) Move() {
	r.Lock()
	defer r.Unlock()
	r.moves = append(r.moves, r.s)
}
