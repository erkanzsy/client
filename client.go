package client

import (
	"fmt"
	"github.com/erkanzsy/client/utils"
	"github.com/streadway/amqp"
	"net/http"
)

type LoggingClient struct {
	client *http.Client
}

func NewLoggingClient(channel *amqp.Channel) *LoggingClient {
	return &LoggingClient{
		client: &http.Client{
			Transport: utils.NewClientLogger(channel),
		},
	}
}

func (lc *LoggingClient) Do(req *http.Request) (*http.Response, error) {
	fmt.Printf("Sending request: %s %s\n", req.Method, req.URL.String())

	resp, err := lc.client.Do(req)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return nil, err
	}

	fmt.Printf("Received response: %s\n", resp.Status)
	return resp, nil
}
