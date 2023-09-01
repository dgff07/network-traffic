package network

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
}

type NetworkService struct {
	*NetworkTraffic
	Mutex *sync.Mutex
}

func BuildNetworkService() NetworkRecorder {
	netTraffic := make(NetworkTraffic)
	return &NetworkService{
		NetworkTraffic: &netTraffic,
		Mutex:          &sync.Mutex{},
	}
}

func (ns *NetworkService) CaptureTraffic(port string, bufferSize int) {
	go captureTrafficRoutine(port, bufferSize, ns.NetworkTraffic, ns.Mutex)
}

func captureTrafficRoutine(port string, bufferSize int, networkTraffic *NetworkTraffic, mutex *sync.Mutex) {
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

		// Lock the mutex before accessing the map
		mutex.Lock()

		// Append the network info to the slice
		(*networkTraffic)[port] = append((*networkTraffic)[port], networkInfo)

		// Check if the buffer size is exceeded and wrap around if needed
		if len((*networkTraffic)[port]) > bufferSize {
			(*networkTraffic)[port] = (*networkTraffic)[port][1:]
		}

		fmt.Println(networkTraffic)

		// Unlock the mutex after updating the map
		mutex.Unlock()

		conn.Close()
	}
}

func getCurrentDateTimeFormatted() string {
	currentTime := time.Now()
	formattedDateTime := currentTime.Format(dateFmt)
	return formattedDateTime
}
