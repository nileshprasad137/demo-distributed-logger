package main

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

func main() {
	// Connect to Redis
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // Replace with your Redis server address
	})

	// Create a new context
	ctx := context.Background()

	// Define the log stream name
	logStream := "log_stream"

	// Start the log subscriber in the background
	go logSubscriber(client, ctx, logStream)

	// Simulate sending log messages from multiple services
	for i := 1; i <= 10; i++ {
		serviceID := fmt.Sprintf("service_%d", i)
		logMessage := fmt.Sprintf("Log message from %s: Log entry %d", serviceID, i)

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

func logSubscriber(client *redis.Client, ctx context.Context, logStream string) {
	// Create the consumer group and consumer name
	consumerGroup := "log_group"
	consumerName := "consumer"

	// Create the consumer group (if not exists)
	// This command creates a new consumer group named log_group for the log_stream stream.
	// The special ID "$" represents the last ID in the stream, so any new messages added to the stream after the creation of the group will be available for consumption by this group.
	client.XGroupCreateMkStream(ctx, logStream, consumerGroup, "$")

	// Poll for new log messages from the Redis Stream
	for {
		// Read log messages from the Redis Stream
		// This command reads up to 10 new log messages from the log_stream for the log_group consumer group.
		// The ">" symbol represents the last received message ID, so it will fetch only new messages that haven't been processed by any consumer in the group.
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
			// For each batch of log messages obtained from the stream, we iterate over the messages and process them.
			// Once a message has been successfully processed, we acknowledge it using the XAck command to mark it as processed by the consumer.
			// This prevents the same message from being processed again:
			for _, message := range logMsg.Messages {
				serviceID, _ := message.Values["serviceID"].(string)
				logText, _ := message.Values["message"].(string)
				fmt.Printf("[%s] %s\n", serviceID, logText)

				// TODO: You can add your custom processing logic here,
				// e.g., store log messages in a database or display them on a dashboard.

				// Acknowledge the log message to mark it as processed
				// The XAck command is called with the message's ID to acknowledge it.
				client.XAck(ctx, logStream, consumerGroup, message.ID)
			}
		}
	}
}
