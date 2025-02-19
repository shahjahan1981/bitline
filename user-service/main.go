package main
import (
	"fmt"
	"net"
	"net/rpc"
)

// UserService defines methods for RPC communication
type UserService struct{}

// GetUser returns user details
func (u *UserService) GetUser(request string, response *string) error {
	*response = "User details for: " + request
	return nil
}

// startServer initializes and starts the User Service
func startServer() {
	service := new(UserService)
	rpc.Register(service)

	listener, err := net.Listen("tcp", ":5000")
	if err != nil {
		fmt.Println("Error starting User Service:", err)
		return
	}
	fmt.Println("User Service is running on port 5000...")

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
	startServer()
}
