package client

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
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

	err = channel.ExchangeDeclare(
		"logs",   // Exchange name
		"direct", // Exchange type
		true,     // Durable
		false,    // Auto-deleted
		false,    // Internal
		false,    // No-wait
		nil,      // Arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare an exchange: %v", err)
	}

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

func TestNewLoggingClientWithBadURL(t *testing.T) {
	conn, err := amqp.Dial("amqp://client-user:client-pass@localhost:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	message := "Hello, RabbitMQ!"
	err = ch.Publish(
		"my-exchange", // Exchange name
		"",            // Routing key
		false,         // Mandatory
		false,         // Immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)
	if err != nil {
		log.Fatalf("Failed to publish a message: %v", err)
	}
}
