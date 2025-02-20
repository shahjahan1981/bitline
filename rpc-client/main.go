package main

import (
    "fmt"
    "log"
    "net/rpc"
)

// Registration request
type RegistrationRequest struct {
    Username string
    Password string
}

// Registration response
type RegistrationResponse struct {
    Message string
}

// Login request
type LoginRequest struct {
    Username string
    Password string
}

// Login response
type LoginResponse struct {
    Token string
}

func main() {
    client, err := rpc.Dial("tcp", "localhost:5004") // Connect to Public API Service
    if err != nil {
        log.Fatal("Error connecting to Public API Service:", err)
    }
    defer client.Close()

    // Step 1: Register User
    regReq := RegistrationRequest{
        Username: "shahjahan",
        Password: "12345678",
    }
    var regRes RegistrationResponse

    err = client.Call("PublicAPIService.RegisterUser", regReq, &regRes)
    if err != nil {
        log.Fatal("Error calling RegisterUser:", err)
    }

    fmt.Println("Registration Response:", regRes.Message)

    // Step 2: Login User
    loginReq := LoginRequest{
        Username: "shahjahan",
        Password: "12345678",
    }
    var loginRes LoginResponse

    err = client.Call("PublicAPIService.LoginUser", loginReq, &loginRes)
    if err != nil {
        log.Fatal("Error calling LoginUser:", err)
    }

    fmt.Println("Login Response:", loginRes.Token)
}
