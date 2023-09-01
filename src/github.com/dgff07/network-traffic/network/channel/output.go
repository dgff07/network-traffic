package channel

import "fmt"

const (
	Console   = 0
	Memory    = 1
	Websocket = 2
	Kafka     = 3
)

type OutputKind struct {
	id     int
	name   string
	writer func(ns *NetworkService)
}

type OutputWriter interface {
	write(ns *NetworkService)
}

type Output struct {
	kind *OutputKind
}

func (o *Output) write(ns *NetworkService) {
	go o.kind.writer(ns)
}

func consoleWriter(ns *NetworkService) {
	for {
		netInfo := <-ns.NetworkTrafficChan
		fmt.Println(netInfo)
	}
}

func memoryWriter(ns *NetworkService) {
	for {
		netInfo := <-ns.NetworkTrafficChan

		ns.NetworkTraffic[netInfo.port] = append(ns.NetworkTraffic[netInfo.port], netInfo.info)

		// Check if the buffer size is exceeded and wrap around if needed
		if len(ns.NetworkTraffic[netInfo.port]) > ns.BufferSize {
			ns.NetworkTraffic[netInfo.port] = ns.NetworkTraffic[netInfo.port][1:]
		}

		fmt.Println(ns.NetworkTraffic)
	}

}

func websocketWriter(ns *NetworkService) {

}

func kafkaWriter(ns *NetworkService) {

}

func GetOutput(id int) (*Output, error) {
	switch id {
	case 0:
		kind := &OutputKind{id: 0, name: "console", writer: consoleWriter}
		return &Output{kind: kind}, nil
	case 1:
		kind := &OutputKind{id: 1, name: "memory", writer: memoryWriter}
		return &Output{kind: kind}, nil
	case 2:
		kind := &OutputKind{id: 2, name: "websocket", writer: websocketWriter}
		return &Output{kind: kind}, nil
	case 3:
		kind := &OutputKind{id: 3, name: "kafka", writer: kafkaWriter}
		return &Output{kind: kind}, nil
	}

	return nil, fmt.Errorf("There is no output kind with id '%d'", id)
}
