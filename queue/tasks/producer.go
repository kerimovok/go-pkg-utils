package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/kerimovok/go-pkg-utils/queue"
)

// Task represents a generic task structure for the tasks queue
type Task struct {
	Service string         `json:"service"`
	Type    string         `json:"type"`
	Payload map[string]any `json:"payload"`
}

// Default queue configuration for tasks
var defaultQueueConfig = &queue.Config{
	ExchangeName:    "tasks",
	ExchangeType:    "topic", // Topic exchange for wildcard routing
	QueueName:       "",      // Producers don't create queues, consumers do
	RoutingKey:      "",      // Dynamic routing key set per-message
	DLXExchangeName: "tasks.dlx",
	DLQName:         "",
	DLQRoutingKey:   "",
}

// Producer wraps the base queue producer for publishing tasks
type Producer struct {
	producer *queue.Producer
	service  string
}

// ProducerConfig holds configuration for the task producer
type ProducerConfig struct {
	// ServiceName is the name of the service publishing tasks
	ServiceName string

	// QueueConfig allows overriding the default queue configuration (optional)
	QueueConfig *queue.Config
}

// NewProducer creates a new task producer
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
		return nil, fmt.Errorf("failed to create task producer: %w", err)
	}

	return &Producer{
		producer: producer,
		service:  config.ServiceName,
	}, nil
}

// Publish publishes a task to the queue with a dynamic routing key
// The routing key is constructed as "tasks.<taskType>"
// Example: taskType "email.verify" becomes routing key "tasks.email.verify"
func (p *Producer) Publish(ctx context.Context, taskType string, payload map[string]any) error {
	if payload == nil {
		payload = make(map[string]any)
	}

	// Add timestamp if not present
	if _, ok := payload["timestamp"]; !ok {
		payload["timestamp"] = time.Now().UTC().Format(time.RFC3339)
	}

	task := Task{
		Service: p.service,
		Type:    taskType,
		Payload: payload,
	}

	data, err := json.Marshal(task)
	if err != nil {
		return fmt.Errorf("failed to marshal task: %w", err)
	}

	// Build routing key: tasks.<taskType>
	// Example: "email.verify" -> "tasks.email.verify"
	routingKey := "tasks." + taskType

	return p.producer.PublishWithRoutingKey(ctx, data, nil, routingKey)
}

// PublishAsync publishes a task asynchronously (fire and forget)
func (p *Producer) PublishAsync(taskType string, payload map[string]any) {
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		_ = p.Publish(ctx, taskType, payload)
	}()
}

// PublishWithCustomRoutingKey publishes a task with a custom routing key
// Use this if you need to override the default "tasks.<taskType>" pattern
func (p *Producer) PublishWithCustomRoutingKey(ctx context.Context, taskType string, payload map[string]any, routingKey string) error {
	if payload == nil {
		payload = make(map[string]any)
	}

	// Add timestamp if not present
	if _, ok := payload["timestamp"]; !ok {
		payload["timestamp"] = time.Now().UTC().Format(time.RFC3339)
	}

	task := Task{
		Service: p.service,
		Type:    taskType,
		Payload: payload,
	}

	data, err := json.Marshal(task)
	if err != nil {
		return fmt.Errorf("failed to marshal task: %w", err)
	}

	return p.producer.PublishWithRoutingKey(ctx, data, nil, routingKey)
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
