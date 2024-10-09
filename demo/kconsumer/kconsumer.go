package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"

	"github.com/IBM/sarama"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// Message represents the structure of the JSON message consumed from Kafka.
type Message struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

var dynamoClient *dynamodb.Client
var tableName = "KafkaMessages" // Update with your DynamoDB table name

// initDynamoDB initializes the DynamoDB client.
func initDynamoDB() {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1")) // Set your AWS region
	if err != nil {
		log.Fatalf("Unable to load AWS SDK config: %v", err)
	}
	dynamoClient = dynamodb.NewFromConfig(cfg)
}

// writeToDynamoDB writes a message to DynamoDB.
func writeToDynamoDB(msg Message) error {
	_, err := dynamoClient.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item: map[string]types.AttributeValue{
			"Key":   &types.AttributeValueMemberS{Value: msg.Key},
			"Value": &types.AttributeValueMemberS{Value: msg.Value},
		},
	})
	return err
}

func main() {
	// Initialize DynamoDB client
	initDynamoDB()

	// Kafka consumer configuration
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	// Define Kafka broker and topic
	brokers := []string{"localhost:9092"}  // Update with your Kafka broker(s)
	topic := "my-topic"                    // Update with your Kafka topic

	// Create a new consumer group
	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		log.Fatalf("Error creating Kafka consumer: %v", err)
	}
	defer consumer.Close()

	// Start consuming the specified topic
	partitionConsumer, err := consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatalf("Error creating Kafka partition consumer: %v", err)
	}
	defer partitionConsumer.Close()

	// Handle graceful shutdown
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	// Consume messages from the topic
	log.Println("Consuming messages from Kafka topic and writing to DynamoDB...")
	for {
		select {
		case msg := <-partitionConsumer.Messages():
			var message Message
			err := json.Unmarshal(msg.Value, &message)
			if err != nil {
				log.Printf("Error unmarshaling message: %v", err)
				continue
			}

			// Write the consumed message to DynamoDB
			if err := writeToDynamoDB(message); err != nil {
				log.Printf("Error writing message to DynamoDB: %v", err)
			} else {
				log.Printf("Message written to DynamoDB: Key=%s, Value=%s\n", message.Key, message.Value)
			}

		case err := <-partitionConsumer.Errors():
			log.Printf("Error consuming message: %v", err)

		case <-signals:
			log.Println("Interrupt signal received, shutting down consumer...")
			return
		}
	}
}
