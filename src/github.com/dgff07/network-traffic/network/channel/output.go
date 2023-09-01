package channel

import (
	"fmt"
	"os"
)

const (
	Console   = 0
	Memory    = 1
	Websocket = 2
	Kafka     = 3
)

type OutputWriter interface {
	write(ns *NetworkService)
}

type ConsoleOutput struct{}

type MemoryOutput struct{}

type KafkaOutput struct{}

type WebSocketOutput struct{}

func (o *ConsoleOutput) write(ns *NetworkService) {
	go func(ns *NetworkService) {
		for {
			netInfo := <-ns.NetworkTrafficChan
			fmt.Println(netInfo)
		}
	}(ns)
}

func (o *MemoryOutput) write(ns *NetworkService) {
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

func (o *WebSocketOutput) write(ns *NetworkService) {
	fmt.Println("Not implemented yet")
	os.Exit(1)
}

func (o *KafkaOutput) write(ns *NetworkService) {
	fmt.Println("Not implemented yet")
	os.Exit(1)
}

func GetOutput(id int) (OutputWriter, error) {
	switch id {
	case 0:
		return &ConsoleOutput{}, nil
	case 1:
		return &MemoryOutput{}, nil
	case 2:
		return &WebSocketOutput{}, nil
	case 3:
		return &KafkaOutput{}, nil
	}

	return nil, fmt.Errorf("There is no output kind with id '%d'", id)
}
