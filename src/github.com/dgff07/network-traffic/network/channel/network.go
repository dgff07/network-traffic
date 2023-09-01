package network

import (
	"fmt"
	"net"
	"time"
)

const (
	dateFmt = "2006-01-02 15:04:05.000"
)

type NetworkTraffic map[string][]string

type NetworkRecorder interface {
	CaptureTraffic(string, int)
	InitChannelReader()
}

// Used to pass info through the channel
type NetworkInfo struct {
	port string
	info string
}

type NetworkService struct {
	NetworkTrafficChan chan NetworkInfo
	BufferSize         int
	NetworkTraffic
}

func BuildNetworkService(bs int) NetworkRecorder {
	return &NetworkService{
		NetworkTrafficChan: make(chan NetworkInfo, bs),
		BufferSize:         bs,
		NetworkTraffic:     make(NetworkTraffic),
	}
}

func (ns *NetworkService) InitChannelReader() {
	go appendNetworkData(ns)
}

func (ns *NetworkService) CaptureTraffic(port string, bufferSize int) {
	go captureTrafficRoutine(port, ns.NetworkTrafficChan)
}

func appendNetworkData(ns *NetworkService) {
	for {
		netInfo := <-ns.NetworkTrafficChan

		(ns.NetworkTraffic)[netInfo.port] = append((ns.NetworkTraffic)[netInfo.port], netInfo.info)

		// Check if the buffer size is exceeded and wrap around if needed
		if len((ns.NetworkTraffic)[netInfo.port]) > ns.BufferSize {
			(ns.NetworkTraffic)[netInfo.port] = (ns.NetworkTraffic)[netInfo.port][1:]
		}

		fmt.Println(ns.NetworkTraffic)
	}
}

func captureTrafficRoutine(port string, trafficChan chan NetworkInfo) {
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Println("Error starting listener:", err)
		return
	}
	defer listener.Close()

	fmt.Printf("Listening on port %s\n", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		remoteAddr := conn.RemoteAddr().String()
		fromIP, _, _ := net.SplitHostPort(remoteAddr)

		currentDate := getCurrentDateTimeFormatted()

		networkInfo := fmt.Sprintf("From %s to %s at %s", fromIP, port, currentDate)

		// Send the network traffic data to the channel
		chanVal := NetworkInfo{
			port: port,
			info: networkInfo,
		}
		trafficChan <- chanVal

		conn.Close()
	}
}

func getCurrentDateTimeFormatted() string {
	currentTime := time.Now()
	formattedDateTime := currentTime.Format(dateFmt)
	return formattedDateTime
}

func (ns *NetworkService) Close() {
	close(ns.NetworkTrafficChan)
}
