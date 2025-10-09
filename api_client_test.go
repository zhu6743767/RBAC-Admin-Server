package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const baseURL = "http://127.0.0.1:8080"

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Nickname string `json:"nickname"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func main() {
	fmt.Println("=== RBAC API Test ===")

	// Wait a moment for server to be ready
	time.Sleep(2 * time.Second)

	// Test 1: Health Check
	fmt.Println("\n1. Testing health check...")
	if err := testHealth(); err != nil {
		fmt.Printf("❌ Health check failed: %v\n", err)
	} else {
		fmt.Println("✅ Health check passed!")
	}

	// Test 2: Registration
	fmt.Println("\n2. Testing registration...")
	if err := testRegistration(); err != nil {
		fmt.Printf("❌ Registration failed: %v\n", err)
	} else {
		fmt.Println("✅ Registration successful!")
	}

	// Test 3: Login
	fmt.Println("\n3. Testing login...")
	if err := testLogin(); err != nil {
		fmt.Printf("❌ Login failed: %v\n", err)
	} else {
		fmt.Println("✅ Login successful!")
	}
}

func testHealth() error {
	resp, err := http.Get(baseURL + "/health")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("status code: %d", resp.StatusCode)
	}

	var result Response
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}

	fmt.Printf("   Response: %+v\n", result)
	return nil
}

func testRegistration() error {
	req := RegisterRequest{
		Username: "testuser10",
		Password: "test123456",
		Email:    "test10@example.com",
		Phone:    "13800138010",
		Nickname: "Test User 10",
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		return err
	}

	resp, err := http.Post(baseURL+"/public/register", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var result Response
	if err := json.Unmarshal(body, &result); err != nil {
		return fmt.Errorf("failed to decode response: %v, body: %s", err, string(body))
	}

	fmt.Printf("   Response: %+v\n", result)

	if result.Code != 200 {
		return fmt.Errorf("registration failed with code %d: %s", result.Code, result.Msg)
	}

	return nil
}

func testLogin() error {
	req := LoginRequest{
		Username: "testuser10",
		Password: "test123456",
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		return err
	}

	resp, err := http.Post(baseURL+"/public/login", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var result Response
	if err := json.Unmarshal(body, &result); err != nil {
		return fmt.Errorf("failed to decode response: %v, body: %s", err, string(body))
	}

	fmt.Printf("   Response: %+v\n", result)

	if result.Code != 200 {
		return fmt.Errorf("login failed with code %d: %s", result.Code, result.Msg)
	}

	return nil
}
