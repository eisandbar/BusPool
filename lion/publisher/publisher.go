package publisher

import (
	"fmt"

	. "github.com/eisandbar/BusPool/lion/typing"
)

type Publisher interface {
	Publish(Bus, Instruction)
}

type EmptyPublisher struct {
}

func (pub EmptyPublisher) Publish(bus Bus, inst Instruction) {
	fmt.Println(bus, inst)
}
