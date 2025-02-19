package main

import (
	"fmt"
	"net"
	"net/rpc"
)

// Jumio Service struct
type JumioService struct{}

// VerifyIdentity method
func (j *JumioService) VerifyIdentity(request string, response *string) error {
	*response = "Identity verification complete: " + request
	return nil
}

func startJumioService() {
	service := new(JumioService)
	rpc.Register(service)

	listener, err := net.Listen("tcp", ":5006")
	if err != nil {
		fmt.Println("Error starting Jumio Service:", err)
		return
	}
	fmt.Println("Jumio Service running on port 5006")

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
	startJumioService()
}
