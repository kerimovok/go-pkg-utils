package queue

import (
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// RetryConfig holds retry configuration
type RetryConfig struct {
	MaxRetries     int
	RetryDelayBase int
	MaxRetryDelay  int
}

// GetRetryCount extracts the retry count from message headers
func GetRetryCount(msg amqp.Delivery) int {
	if msg.Headers != nil {
		if retryCount, exists := msg.Headers["x-retry-count"]; exists {
			if count, ok := retryCount.(int32); ok {
				return int(count)
			}
		}
	}
	return 0
}

// CalculateRetryDelay calculates delay with exponential backoff
func CalculateRetryDelay(retryCount int, config RetryConfig) time.Duration {
	// Exponential backoff: baseDelay * 2^retryCount
	delay := time.Duration(config.RetryDelayBase) * time.Duration(1<<retryCount) * time.Second

	// Cap at max delay
	if delay > time.Duration(config.MaxRetryDelay)*time.Second {
		delay = time.Duration(config.MaxRetryDelay) * time.Second
	}

	return delay
}

// ScheduleRetry schedules a message for retry with delay
func ScheduleRetry(
	channel *amqp.Channel,
	config *Config,
	body []byte,
	headers amqp.Table,
	delay time.Duration,
) {
	// In a production system, you might want to use a proper delay queue
	// For now, we'll use a simple goroutine with sleep
	go func() {
		time.Sleep(delay)

		// Publish back to the main queue with updated headers
		err := channel.Publish(
			config.ExchangeName, // exchange
			config.RoutingKey,   // routing key
			false,                // mandatory
			false,                // immediate
			amqp.Publishing{
				ContentType:  "application/json",
				Body:         body,
				Headers:      headers,
				DeliveryMode: amqp.Persistent,
			},
		)

		if err != nil {
			log.Printf("Failed to schedule retry: %v", err)
		} else {
			log.Printf("Scheduled retry with delay %v", delay)
		}
	}()
}
