package main

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const webPort = "80"

type Config struct {
	Rabbit *amqp.Connection
}

func main() {
	//try to connect to rabbitmq
	rabbitConn, err := connect()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer rabbitConn.Close()
	log.Println("Connected to RabbitMQ")

	// Create an instance of the Config struct.
	app := Config{
		Rabbit: rabbitConn,
	}

	// Print a log message indicating that the broker service is starting on the specified port.
	log.Printf("Starting broker service on port %s", webPort)

	// Create an HTTP server with the specified address and use the app's routes() method as the request handler.
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort), // Bind the server to the specified port.
		Handler: app.routes(),                // Use the routes() method of the Config instance as the request handler.
	}

	// Start the HTTP server and listen for incoming requests.
	err = srv.ListenAndServe()

	// Check for any errors that may occur while starting the server.
	if err != nil {
		log.Panic(err) // If there's an error, log it as a panic.
	}
}

func connect() (*amqp.Connection, error) {
	var counts int64
	var backOff = 1 * time.Second
	var connection *amqp.Connection

	//dont continue unitill rabbit is ready

	for {
		c, err := amqp.Dial("amqp://guest:guest@rabbitmq")
		if err != nil {
			fmt.Println("RabbitMQ not yet ready...")
			counts++
		} else {
			log.Println("Connected to RabbitMQ")
			connection = c
			break
		}
		if counts > 5 {
			fmt.Println(err)
			return nil, err
		}
		backOff = time.Duration(math.Pow(float64(counts), 2)) * time.Second
		log.Println("backing off...")
		time.Sleep(backOff)
		continue
	}
	return connection, nil
}
