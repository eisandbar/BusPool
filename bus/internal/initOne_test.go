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
		CheckCalls(t, &reporter, want)
		testChan <- time.Now()
		want = append(want, "Test")
	}
	close(testChan)
}

func CheckCalls(t testing.TB, r *mockReportable, want interface{}) {
	t.Helper()
	r.Lock()
	defer r.Unlock()
	assert.Equal(t, want, r.calls)
}

type mockReportable struct {
	s     string
	calls []string
	sync.Mutex
}

func (r *mockReportable) Id() string {
	return r.s
}

func (r *mockReportable) Report(mqtt.Client) {
	r.Lock()
	defer r.Unlock()
	r.calls = append(r.calls, r.s)
}
