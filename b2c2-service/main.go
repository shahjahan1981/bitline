package main

import (
	"fmt"
	"net"
	"net/rpc"
)

// B2C2 Service struct
type B2C2Service struct{}

// ExecuteTrade method
func (b *B2C2Service) ExecuteTrade(request string, response *string) error {
	*response = "Trade executed: " + request
	return nil
}

func startB2C2Service() {
	service := new(B2C2Service)
	rpc.Register(service)

	listener, err := net.Listen("tcp", ":5007")
	if err != nil {
		fmt.Println("Error starting B2C2 Service:", err)
		return
	}
	fmt.Println("B2C2 Service running on port 5007")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Connection error:", err)
			continue
		}
		go rpc.ServeConn(conn)
	}
}

func main() {
	startB2C2Service()
}
