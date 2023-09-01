package channel

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
	OutputWriter
}

func BuildNetworkService(bs int, outputId int) (NetworkRecorder, error) {
	output, err := GetOutput(outputId)

	if err != nil {
		return nil, err
	}

	return &NetworkService{
		NetworkTrafficChan: make(chan NetworkInfo, bs),
		BufferSize:         bs,
		NetworkTraffic:     make(NetworkTraffic),
		OutputWriter:       output,
	}, nil
}

func (ns *NetworkService) InitChannelReader() {
	ns.OutputWriter.write(ns)
}

func (ns *NetworkService) CaptureTraffic(port string, bufferSize int) {
	go captureTrafficRoutine(port, ns.NetworkTrafficChan)
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
		conn.Close()

		fromIP, _, _ := net.SplitHostPort(remoteAddr)

		currentDate := getCurrentDateTimeFormatted()

		networkInfo := fmt.Sprintf("From %s to %s at %s", fromIP, port, currentDate)

		// Send the network traffic data to the channel
		chanVal := NetworkInfo{
			port: port,
			info: networkInfo,
		}
		trafficChan <- chanVal
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
