package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"sync"
	"time"
)

const (
	buffMax        = 50
	buffMin        = 1
	portsMax       = 4
	portNumberMax  = 65535
	portNumberMin  = 1
	argsErrMsg     = "[ERROR] Invalid arguments"
	buffSizeErrMsg = "[ERROR] Invalid buffer size, you must pass an integer value lower or equal than 50 and greater than 0"
	portSizeErrMsg = "[ERROR] Invalid number of ports, you must pass until 4 ports at max"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println(argsErrMsg)
		printHelp()
		os.Exit(1)
	}

	ports := os.Args[1 : len(os.Args)-1]
	if len(ports) > portsMax {
		fmt.Println(portSizeErrMsg)
		os.Exit(1)
	}
	if !validatePorts(ports) {
		os.Exit(1)
	}

	bufferSize, err := strconv.Atoi(os.Args[len(os.Args)-1])
	if err != nil {
		fmt.Println("Invalid buffer size:", err)
		os.Exit(1)
	}

	if bufferSize > buffMax || bufferSize < buffMin {
		fmt.Println(buffSizeErrMsg)
		os.Exit(1)
	}

	networkTraffic := make(map[string][]string)
	var mutex sync.Mutex

	for _, port := range ports {
		go captureTraffic(port, bufferSize, &networkTraffic, &mutex)
	}

	select {}
}

func printHelp() {
	fmt.Println("To execute this program you must do:\n\ngo run main.go <port1> <port2> ... <portN> <bufferSize>\n\nYou can pass multiple ports and the buffer size to limit the information stored in memory")
}

func isValidPort(port string) bool {
	portNumber, err := strconv.Atoi(port)
	if err != nil {
		return false
	}

	return portNumber >= portNumberMin && portNumber <= portNumberMax
}

func validatePorts(ports []string) bool {
	seen := make(map[string]bool) // to validate repeated ports

	for _, port := range ports {
		if !isValidPort(port) {
			fmt.Printf("Invalid port: %s\n", port)
			return false
		}
		if seen[port] {
			fmt.Printf("Repeated port: %s\n", port)
			return false
		}
		seen[port] = true
	}

	return true
}

func captureTraffic(port string, bufferSize int, networkTraffic *map[string][]string, mutex *sync.Mutex) {
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
	formattedDateTime := currentTime.Format("2006-01-02 15:04:05.000")
	return formattedDateTime
}
