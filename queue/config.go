package queue

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

// Config holds the configuration for RabbitMQ queues and exchanges
type Config struct {
	ExchangeName    string
	ExchangeType    string // "direct", "topic", "fanout", "headers" - defaults to "direct"
	QueueName       string
	RoutingKey      string
	DLXExchangeName string
	DLQName         string
	DLQRoutingKey   string
}

// getExchangeType returns the exchange type, defaulting to "direct" if not set
func (qc *Config) getExchangeType() string {
	if qc.ExchangeType == "" {
		return "direct"
	}
	return qc.ExchangeType
}

// GetQueueArguments returns the queue arguments for both main queue and DLQ
func (qc *Config) GetQueueArguments() amqp.Table {
	return amqp.Table{
		"x-message-ttl":             int32(24 * 60 * 60 * 1000), // 24 hours TTL
		"x-max-priority":            int32(10),                  // Priority support
		"x-overflow":                "drop-head",                // Drop oldest when full
		"x-dead-letter-exchange":    qc.DLXExchangeName,         // Dead letter exchange
		"x-dead-letter-routing-key": qc.DLQRoutingKey,           // Routing key for DLQ
	}
}

// SetupExchange declares the main exchange
func (qc *Config) SetupExchange(ch *amqp.Channel) error {
	return ch.ExchangeDeclare(
		qc.ExchangeName,      // name
		qc.getExchangeType(), // type (defaults to "direct")
		true,                 // durable
		false,                // auto-deleted
		false,                // internal
		false,                // no-wait
		nil,                  // arguments
	)
}

// SetupDeadLetterExchange declares the dead letter exchange
func (qc *Config) SetupDeadLetterExchange(ch *amqp.Channel) error {
	return ch.ExchangeDeclare(
		qc.DLXExchangeName, // dead letter exchange name
		"direct",           // exchange type
		true,               // durable
		false,              // auto-deleted
		false,              // internal
		false,              // no-wait
		nil,                // arguments
	)
}

// SetupDeadLetterQueue declares the dead letter queue
func (qc *Config) SetupDeadLetterQueue(ch *amqp.Channel) error {
	dlq, err := ch.QueueDeclare(
		qc.DLQName, // dead letter queue name
		true,       // durable
		false,      // delete when unused
		false,      // exclusive
		false,      // no-wait
		nil,        // arguments
	)
	if err != nil {
		return err
	}

	// Bind dead letter queue to exchange
	return ch.QueueBind(
		dlq.Name,           // queue name
		qc.DLQRoutingKey,   // routing key
		qc.DLXExchangeName, // exchange
		false,              // no-wait
		nil,                // arguments
	)
}

// SetupMainQueue declares the main queue with DLQ configuration
func (qc *Config) SetupMainQueue(ch *amqp.Channel) error {
	_, err := ch.QueueDeclare(
		qc.QueueName, // name
		true,         // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		qc.GetQueueArguments(),
	)
	if err != nil {
		return err
	}

	// Bind queue to exchange
	return ch.QueueBind(
		qc.QueueName,    // queue name
		qc.RoutingKey,   // routing key
		qc.ExchangeName, // exchange
		false,
		nil,
	)
}

// SetupAllQueues sets up all exchanges and queues
func (qc *Config) SetupAllQueues(ch *amqp.Channel) error {
	// Setup dead letter exchange and queue first (if configured)
	if qc.DLXExchangeName != "" {
		if err := qc.SetupDeadLetterExchange(ch); err != nil {
			return err
		}
	}

	if qc.DLQName != "" {
		if err := qc.SetupDeadLetterQueue(ch); err != nil {
			return err
		}
	}

	// Setup main exchange
	if err := qc.SetupExchange(ch); err != nil {
		return err
	}

	// Setup main queue (skip if not configured - producer-only mode)
	if qc.QueueName != "" {
		return qc.SetupMainQueue(ch)
	}

	return nil
}
