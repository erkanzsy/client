package client

import (
	"bytes"
	"fmt"
	"github.com/streadway/amqp"
	"net/http"
	"testing"
)

func TestNewLoggingClient(t *testing.T) {
	rabbitMQURL := "amqp://client-user:client-pass@localhost:5672/"
	conn, err := amqp.Dial(rabbitMQURL)
	if err != nil {
		fmt.Printf("Error connecting to RabbitMQ: %v\n", err)
		return
	}
	defer conn.Close()

	channel, err := conn.Channel()
	if err != nil {
		fmt.Printf("Error creating RabbitMQ channel: %v\n", err)
		return
	}

	defer channel.Close()

	client := NewLoggingClient(channel)

	req, err := http.NewRequest("GET", "https://www.example.com", nil)
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error sending request: %v\n", err)
		return
	}
	defer resp.Body.Close()
}

func TestNewLoggingClientPostRequest(t *testing.T) {
	rabbitMQURL := "amqp://client-user:client-pass@localhost:5672/"
	conn, err := amqp.Dial(rabbitMQURL)
	if err != nil {
		fmt.Printf("Error connecting to RabbitMQ: %v\n", err)
		return
	}
	defer conn.Close()

	channel, err := conn.Channel()
	if err != nil {
		fmt.Printf("Error creating RabbitMQ channel: %v\n", err)
		return
	}

	defer channel.Close()

	client := NewLoggingClient(channel)

	url := "https://jsonplaceholder.typicode.com/posts" // Replace with your desired URL

	// Define the request body as a JSON string
	requestBody := []byte(`{
		"title": "foo",
		"body": "bar",
		"userId": 1
	}`)

	// Create a new POST request with the request body
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	// Set the content type header
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error sending request: %v\n", err)
		return
	}
	defer resp.Body.Close()
}
