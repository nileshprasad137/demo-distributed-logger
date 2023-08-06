// main.go

package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/redis/go-redis/v9"
)

func main() {
	// Initialize log producer
	go produceLogs()

	// Start log consumer
	go consumeLogs()

	// Wait for termination signal to gracefully stop consumer
	waitForTermination()
}

func produceLogs() {
	// Connect to Redis
	client := redis.NewClient(&redis.Options{
		Addr: "redis:6379", // Use "redis" as the Redis server address (Docker service name)
	})

	// Create a new context
	ctx := context.Background()

	// Define the log stream name
	logStream := "log_stream"

	// Simulate sending log messages from multiple services
	for i := 1; i <= 100; i++ {
		serviceID := fmt.Sprintf("service_%d", i)
		logMessage := fmt.Sprintf("Log message from %s: Log entry %d", serviceID, rand.Intn(1000))

		// Publish the log message to the Redis Stream
		client.XAdd(ctx, &redis.XAddArgs{
			Stream: logStream,
			Values: map[string]interface{}{
				"serviceID": serviceID,
				"message":   logMessage,
			},
		})

		fmt.Printf("Log message sent: %s\n", logMessage)
		time.Sleep(time.Second)
	}

	// Close the Redis client connection
	client.Close()
}

func consumeLogs() {
	// Connect to Redis
	client := redis.NewClient(&redis.Options{
		Addr: "redis:6379", // Use "redis" as the Redis server address (Docker service name)
	})

	// Create a new context
	ctx := context.Background()

	// Define the log stream name
	logStream := "log_stream"

	// Define the consumer group
	consumerGroup := "log_group"

	// Create the consumer group (if not exists)
	client.XGroupCreateMkStream(ctx, logStream, consumerGroup, "$")

	// Get the hostname of the container to generate a unique consumer name
	hostname, err := os.Hostname()
	if err != nil {
		fmt.Println("Error getting hostname:", err)
		return
	}
	consumerName := fmt.Sprintf("consumer_%s", strings.TrimSpace(hostname))

	// Poll for new log messages from the Redis Stream
	for {
		// Read log messages from the Redis Stream
		logMessages, err := client.XReadGroup(ctx, &redis.XReadGroupArgs{
			Group:    consumerGroup,
			Consumer: consumerName,
			Streams:  []string{logStream, ">"},
			Block:    0, // Block and wait for new log messages indefinitely
			Count:    10,
			NoAck:    false,
		}).Result()

		if err != nil {
			fmt.Println("Error reading log messages:", err)
			break
		}

		// Process log messages
		for _, logMsg := range logMessages {
			for _, message := range logMsg.Messages {
				serviceID, _ := message.Values["serviceID"].(string)
				logText, _ := message.Values["message"].(string)
				fmt.Printf("[%s] %s\n", serviceID, logText)

				// TODO: You can add your custom processing logic here,
				// e.g., store log messages in a database or display them on a dashboard.

				// Acknowledge the log message to mark it as processed
				client.XAck(ctx, logStream, consumerGroup, message.ID)
			}
		}
	}
}

func waitForTermination() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
	fmt.Println("Termination signal received. Stopping consumer gracefully.")
	os.Exit(0)
}
