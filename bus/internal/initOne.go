package internal

import (
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func InitOne(client mqtt.Client, r Reportable, tick <-chan time.Time) {
	r.Subscribe(client)
	for {
		_, more := <-tick
		if more {
			r.Move()
			r.Report(client)
		} else {
			return
		}
	}
}

type Reportable interface {
	Report(mqtt.Client)
	Subscribe(mqtt.Client)
	Move()
}
