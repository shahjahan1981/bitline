package main

import (
	"fmt"
	"net"
	"net/rpc"
	"time"
)

// MicroserviceB provides RPC methods
type MicroserviceB struct{}

// Greet method responds to RPC calls with a message
func (s *MicroserviceB) Greet(req string, res *string) error {
	*res = "Hello from Microservice B!"
	return nil
}

// startMicroserviceB initializes the service and listens on port 5678
func startMicroserviceB() {
	service := new(MicroserviceB)
	rpc.Register(service)

	listener, err := net.Listen("tcp", ":5678")
	if err != nil {
		fmt.Println("Microservice B failed to start:", err)
		return
	}
	fmt.Println("Microservice B is running on port 5678...")

	// Start trying to connect to Microservice A
	go connectToMicroserviceA()

	rpc.Accept(listener)
}

// connectToMicroserviceA keeps trying to connect to Microservice A
func connectToMicroserviceA() {
	for {
		client, err := rpc.Dial("tcp", "localhost:1234")
		if err == nil {
			var response string
			err = client.Call("MicroserviceA.Greet", "", &response)
			if err == nil {
				fmt.Println("Message from Microservice A:", response)
				break
			}
		}
		fmt.Println("Retrying connection to Microservice A...")
		time.Sleep(2 * time.Second)
	}
}

func main() {
	startMicroserviceB()
}
