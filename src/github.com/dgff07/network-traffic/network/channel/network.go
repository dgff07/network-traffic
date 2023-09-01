package channel

import (
	"fmt"
	"net"
	"sync"
	"time"
)

const (
	dateFmt = "2006-01-02 15:04:05.000"
)

type NetworkTraffic map[string][]string

type NetworkRecorder interface {
	CaptureTraffic(string, int)
	GetChannel() chan NetworkInfo
	GetBufferSize() int
	SaveNetworkTraffic(*NetworkInfo)
}

// Used to pass info through the channel
type NetworkInfo struct {
	port string
	info string
}

type NetworkService struct {
	networkTrafficChan chan NetworkInfo
	bufferSize         int
	netTraffic         NetworkTraffic
	mutex              sync.RWMutex
}

func BuildNetworkService(bs int) NetworkRecorder {

	return &NetworkService{
		networkTrafficChan: make(chan NetworkInfo),
		bufferSize:         bs,
		netTraffic:         make(NetworkTraffic),
	}
}

func (ns *NetworkService) SaveNetworkTraffic(netInfo *NetworkInfo) {
	ns.mutex.Lock() //Lock write
	defer ns.mutex.Unlock()
	ns.netTraffic[netInfo.port] = append(ns.netTraffic[netInfo.port], netInfo.info)

	// Check if the buffer size is exceeded and wrap around if needed
	if len(ns.netTraffic[netInfo.port]) > ns.GetBufferSize() {
		ns.netTraffic[netInfo.port] = ns.netTraffic[netInfo.port][1:]
	}
	fmt.Println(ns.netTraffic)
}

func (ns *NetworkService) GetChannel() chan NetworkInfo {
	return ns.networkTrafficChan
}

func (ns *NetworkService) GetBufferSize() int {
	return ns.bufferSize
}

func (ns *NetworkService) CaptureTraffic(port string, bufferSize int) {
	go captureTrafficRoutine(port, ns.networkTrafficChan)
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
	close(ns.networkTrafficChan)
}
