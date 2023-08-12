package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type ClientLogger struct {
	Transport http.RoundTripper
	Channel   *amqp.Channel
}

func NewClientLogger(channel *amqp.Channel) *ClientLogger {
	return &ClientLogger{
		Transport: http.DefaultTransport,
		Channel:   channel,
	}
}

type LogData struct {
	RequestURL      string
	RequestMethod   string
	ResponseStatus  string
	ResponseTime    time.Duration
	ResponseHeaders http.Header
	ResponseBody    string
	ResponseError   string
	RequestHeaders  http.Header
	RequestBody     string
}

func (cl *ClientLogger) RoundTrip(req *http.Request) (*http.Response, error) {
	startTime := time.Now()

	resp, err := cl.Transport.RoundTrip(req)

	fmt.Println("Request URL: ", req.URL.String())
	fmt.Println("Request Method: ", req.Method)
	fmt.Println("Response Status: ", resp.Status)
	fmt.Println("Response Headers: ", resp.Header)
	fmt.Println("Response Body: ", resp.Body)
	fmt.Println("Response Error: ", err)
	fmt.Println("Request Headers: ", req.Header)
	fmt.Println("Request Body: ", req.Body)

	endTime := time.Now()
	elapsedTime := endTime.Sub(startTime)

	logData := LogData{
		RequestURL:     req.URL.String(),
		RequestMethod:  req.Method,
		ResponseTime:   elapsedTime,
		RequestHeaders: req.Header,
	}

	if err != nil {
		logData.ResponseError = err.Error()
	} else {
		logData.ResponseStatus = resp.Status
		logData.ResponseHeaders = resp.Header

		respBody, _ := ioutil.ReadAll(resp.Body)
		logData.ResponseBody = string(respBody)
		resp.Body = ioutil.NopCloser(bytes.NewBuffer(respBody))
	}

	cl.sendToRabbitMQ(logData)

	return resp, err
}

func (c *ClientLogger) sendToRabbitMQ(logData LogData) {
	jsonData, err := json.Marshal(logData)
	if err != nil {
		fmt.Printf("Error encoding log data: %v\n", err)
		return
	}

	err = c.Channel.ExchangeDeclare(
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

	err = c.Channel.Publish(
		"logs",
		"",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        jsonData,
		},
	)
	fmt.Printf("pushed data: %v\n", logData)

	if err != nil {
		fmt.Printf("Error sending log data to RabbitMQ: %v\n", err)
	}
}