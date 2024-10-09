package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/IBM/sarama"
)

// KafkaConfig holds the Kafka configuration.
type KafkaConfig struct {
	Brokers []string
	Topic   string
}

// Message represents the JSON message to be sent to Kafka.
type Message struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

var kafkaProducer sarama.SyncProducer
var kafkaConfig KafkaConfig

func initKafka() {
	kafkaConfig = KafkaConfig{
		Brokers: []string{"localhost:9092"}, // Update with your Kafka broker(s)
		Topic:   "my-topic",                 // Update with your Kafka topic
	}

	// Kafka producer configuration
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true

	// Initialize the Kafka producer
	var err error
	kafkaProducer, err = sarama.NewSyncProducer(kafkaConfig.Brokers, config)
	if err != nil {
		log.Fatalf("Failed to start Kafka producer: %v", err)
	}
}

func handleKafkaMessage(w http.ResponseWriter, r *http.Request) {
	var msg Message
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&msg); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Prepare the Kafka message
	kafkaMsg := &sarama.ProducerMessage{
		Topic: kafkaConfig.Topic,
		Key:   sarama.StringEncoder(msg.Key),
		Value: sarama.StringEncoder(msg.Value),
		Timestamp: time.Now(),
	}

	// Send the message to Kafka
	partition, offset, err := kafkaProducer.SendMessage(kafkaMsg)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to send message to Kafka: %v", err), http.StatusInternalServerError)
		return
	}

	// Respond with the partition and offset of the sent message
	fmt.Fprintf(w, "Message sent to partition %d, offset %d\n", partition, offset)
}

func main() {
	// Initialize Kafka producer
	initKafka()
	defer kafkaProducer.Close()

	// Define HTTP routes
	http.HandleFunc("/send", handleKafkaMessage)

	// Start the HTTP server
	log.Println("Starting server on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
