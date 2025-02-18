package main

import (
	"fmt"
	"net"
	"net/rpc"
	"time"
)

// MicroserviceA provides RPC methods
type MicroserviceA struct{}

// Greet method responds to RPC calls with a message
func (s *MicroserviceA) Greet(req string, res *string) error {
	*res = "Hello from Microservice A!"
	return nil
}

// startMicroserviceA initializes the service and listens on port 1234
func startMicroserviceA() {
	service := new(MicroserviceA)
	rpc.Register(service)

	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		fmt.Println("Microservice A failed to start:", err)
		return
	}
	fmt.Println("Microservice A is running on port 1234...")

	// Start trying to connect to Microservice B
	go connectToMicroserviceB()

	rpc.Accept(listener)
}

// connectToMicroserviceB keeps trying to connect to Microservice B
func connectToMicroserviceB() {
	for {
		client, err := rpc.Dial("tcp", "localhost:5678")
		if err == nil {
			var response string
			err = client.Call("MicroserviceB.Greet", "", &response)
			if err == nil {
				fmt.Println(" Message from Microservice B:", response)
				break
			}
		}
		fmt.Println("Retrying connection to Microservice B...")
		time.Sleep(2 * time.Second)
	}
}

func main() {
	startMicroserviceA()
}
