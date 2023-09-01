package channel

import (
	"fmt"
	"os"
	"sync"
)

const (
	Console   = 0
	Memory    = 1
	WebSocket = 2
	Kafka     = 3
)

type OutputWriter interface {
	Write(ns *NetworkService)
}

type ConsoleOutput struct{}

type MemoryOutput struct{}

type KafkaOutput struct{}

type WebSocketOutput struct{}

type OutputFactory struct {
	instances map[int]OutputWriter
	mutex     sync.Mutex
}

func NewOutputFactory() *OutputFactory {
	return &OutputFactory{
		instances: make(map[int]OutputWriter),
	}
}

func (o *ConsoleOutput) Write(ns *NetworkService) {

	go func(ns *NetworkService) {
		for {
			netInfo := <-ns.NetworkTrafficChan
			fmt.Println(netInfo)
		}
	}(ns)
}

func (o *MemoryOutput) Write(ns *NetworkService) {
	go func(ns *NetworkService) {
		for {
			netInfo := <-ns.NetworkTrafficChan

			ns.NetworkTraffic[netInfo.port] = append(ns.NetworkTraffic[netInfo.port], netInfo.info)

			// Check if the buffer size is exceeded and wrap around if needed
			if len(ns.NetworkTraffic[netInfo.port]) > ns.BufferSize {
				ns.NetworkTraffic[netInfo.port] = ns.NetworkTraffic[netInfo.port][1:]
			}

			fmt.Println(ns.NetworkTraffic)
		}
	}(ns)

}

func (o *WebSocketOutput) Write(ns *NetworkService) {
	fmt.Println("Not implemented yet")
	os.Exit(1)
}

func (o *KafkaOutput) Write(ns *NetworkService) {
	fmt.Println("Not implemented yet")
	os.Exit(1)
}

func (of *OutputFactory) GetOutput(id int) (OutputWriter, error) {
	of.mutex.Lock()
	defer of.mutex.Unlock()

	if _, ok := of.instances[id]; ok {
		return nil, fmt.Errorf("There is already a output writer with the id '%d'", id)
	}

	var output OutputWriter
	switch id {
	case Console:
		output = &ConsoleOutput{}
	case Memory:
		output = &MemoryOutput{}
	case WebSocket:
		output = &WebSocketOutput{}
	case Kafka:
		output = &KafkaOutput{}
	default:
		return nil, fmt.Errorf("There is no output kind with id '%d'", id)
	}

	of.instances[id] = output
	return output, nil
}
