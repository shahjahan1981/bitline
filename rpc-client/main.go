package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/rpc"
	"regexp"
)

// Registration request
type RegistrationRequest struct {
	FullName        string `json:"full_name"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

// Registration response
type RegistrationResponse struct {
	Message string `json:"message"`
}

// Validate input fields before sending request
func validateRegistration(req RegistrationRequest) error {
	if req.FullName == "" {
		return fmt.Errorf("full name cannot be empty")
	}
	if req.Email == "" || !isValidEmail(req.Email) {
		return fmt.Errorf("invalid email format")
	}
	if req.Password == "" || req.ConfirmPassword == "" {
		return fmt.Errorf("password fields cannot be empty")
	}
	if req.Password != req.ConfirmPassword {
		return fmt.Errorf("password and confirm password do not match")
	}
	return nil
}

// Email validation function
func isValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

func main() {
	client, err := rpc.Dial("tcp", "localhost:5000") // Connect to User Service
	if err != nil {
		log.Fatal("Error connecting to User Service:", err)
	}
	defer client.Close()

	// Step 1: Register User
	regReq := RegistrationRequest{
		FullName:        "Shahjahan Aslam",
		Email:           "shahjahanaslam12345@gmail.com",
		Password:        "12345678",
		ConfirmPassword: "12345678",
	}

	// Validate registration request
	if err := validateRegistration(regReq); err != nil {
		log.Fatal("Validation Error:", err)
	}

	var regRes RegistrationResponse
	err = client.Call("UserService.RegisterUser", regReq, &regRes)
	if err != nil {
		log.Fatal("Error calling RegisterUser:", err)
	}

	fmt.Println("Registration Response:", regRes.Message)

	// Output JSON format response
	jsonOutput, _ := json.MarshalIndent(regRes, "", "  ")
	fmt.Println("JSON Output:", string(jsonOutput))
}
