package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/dgff07/network-traffic/src/github.com/dgff07/network-traffic/network"
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

	netService := network.BuildNetworkService()

	for _, port := range ports {
		netService.CaptureTraffic(port, bufferSize)
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
