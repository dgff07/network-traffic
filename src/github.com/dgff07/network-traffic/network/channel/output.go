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
	Write(ns NetworkRecorder)
}

type consoleOutput struct{}

type memoryOutput struct{}

type kafkaOutput struct{}

type webSocketOutput struct{}

type OutputFactory struct {
	instances map[int]OutputWriter
	mutex     sync.Mutex
}

func NewOutputFactory() *OutputFactory {
	return &OutputFactory{
		instances: make(map[int]OutputWriter),
	}
}

func (o *consoleOutput) Write(ns NetworkRecorder) {

	go func(ns NetworkRecorder) {
		for {
			netInfo := <-ns.GetChannel()
			fmt.Println(netInfo)
		}
	}(ns)
}

func (o *memoryOutput) Write(ns NetworkRecorder) {
	go func(ns NetworkRecorder) {
		for {
			netInfo := <-ns.GetChannel()
			ns.SaveNetworkTraffic(&netInfo)
		}
	}(ns)
}

func (o *webSocketOutput) Write(ns NetworkRecorder) {
	fmt.Println("Not implemented yet")
	os.Exit(1)
}

func (o *kafkaOutput) Write(ns NetworkRecorder) {
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
		output = &consoleOutput{}
		break
	case Memory:
		output = &memoryOutput{}
		break
	case WebSocket:
		output = &webSocketOutput{}
		break
	case Kafka:
		output = &kafkaOutput{}
		break
	default:
		return nil, fmt.Errorf("There is no output kind with id '%d'", id)
	}

	of.instances[id] = output
	return output, nil
}
