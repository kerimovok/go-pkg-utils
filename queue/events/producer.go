package events

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/kerimovok/go-pkg-utils/queue"
)

// Event represents a generic event structure for the event queue
type Event struct {
	Service string         `json:"service"`
	Type    string         `json:"type"`
	Payload map[string]any `json:"payload"`
}

// Default queue configuration for events
var defaultQueueConfig = &queue.Config{
	ExchangeName:    "events",
	ExchangeType:    "direct", // Direct exchange for exact routing key match
	QueueName:       "",       // Producers don't create queues, consumers do
	RoutingKey:      "event",
	DLXExchangeName: "events.dlx",
	DLQName:         "",
	DLQRoutingKey:   "event.failed",
}

// Producer wraps the base queue producer for publishing events
type Producer struct {
	producer *queue.Producer
	service  string
}

// ProducerConfig holds configuration for the event producer
type ProducerConfig struct {
	// ServiceName is the name of the service publishing events
	ServiceName string

	// QueueConfig allows overriding the default queue configuration (optional)
	QueueConfig *queue.Config
}

// NewProducer creates a new event producer
func NewProducer(connConfig queue.ConnectionConfig, config ProducerConfig) (*Producer, error) {
	if config.ServiceName == "" {
		return nil, fmt.Errorf("service name is required")
	}

	queueConfig := defaultQueueConfig
	if config.QueueConfig != nil {
		queueConfig = config.QueueConfig
	}

	producer, err := queue.NewProducer(connConfig, queueConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create event producer: %w", err)
	}

	return &Producer{
		producer: producer,
		service:  config.ServiceName,
	}, nil
}

// Publish publishes an event to the queue
func (p *Producer) Publish(ctx context.Context, eventType string, payload map[string]any) error {
	if payload == nil {
		payload = make(map[string]any)
	}

	// Add timestamp if not present
	if _, ok := payload["timestamp"]; !ok {
		payload["timestamp"] = time.Now().UTC().Format(time.RFC3339)
	}

	event := Event{
		Service: p.service,
		Type:    eventType,
		Payload: payload,
	}

	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	return p.producer.Publish(ctx, data, nil)
}

// PublishAsync publishes an event asynchronously (fire and forget)
func (p *Producer) PublishAsync(eventType string, payload map[string]any) {
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		_ = p.Publish(ctx, eventType, payload)
	}()
}

// IsConnected returns true if the producer has a valid connection
func (p *Producer) IsConnected() bool {
	return p.producer.IsConnected()
}

// Close closes the producer connection
func (p *Producer) Close() error {
	return p.producer.Close()
}

// ServiceName returns the service name configured for this producer
func (p *Producer) ServiceName() string {
	return p.service
}
