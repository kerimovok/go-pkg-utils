package queue

import (
	"fmt"
	"log"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// ConnectionConfig holds RabbitMQ connection details
type ConnectionConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	VHost    string
}

// MessageHandler is a function that processes a message
// It should return an error if processing failed and retry is needed
type MessageHandler func(msg amqp.Delivery) error

// Consumer is a RabbitMQ consumer with automatic reconnection
type Consumer struct {
	conn      *amqp.Connection
	channel   *amqp.Channel
	mu        sync.RWMutex
	config    *Config
	connConfig ConnectionConfig
	retryConfig RetryConfig
	handler   MessageHandler
	consuming bool
	stopChan  chan struct{}
	stopOnce  sync.Once
}

// NewConsumer creates a new RabbitMQ consumer with automatic reconnection
func NewConsumer(connConfig ConnectionConfig, queueConfig *Config, retryConfig RetryConfig, handler MessageHandler) (*Consumer, error) {
	url := fmt.Sprintf("amqp://%s:%s@%s:%s/%s",
		connConfig.Username,
		connConfig.Password,
		connConfig.Host,
		connConfig.Port,
		connConfig.VHost,
	)

	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %v", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open channel: %v", err)
	}

	// Setup all queues and exchanges
	if err := queueConfig.SetupAllQueues(ch); err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to setup queues: %v", err)
	}

	consumer := &Consumer{
		conn:       conn,
		channel:    ch,
		config:     queueConfig,
		connConfig: connConfig,
		retryConfig: retryConfig,
		handler:    handler,
		consuming:  false,
		stopChan:   make(chan struct{}),
	}

	consumer.setupConnectionRecovery()

	return consumer, nil
}

// StartConsuming starts consuming messages from the queue
func (c *Consumer) StartConsuming() error {
	c.mu.Lock()
	if c.consuming {
		c.mu.Unlock()
		return nil // Already consuming
	}
	c.consuming = true
	c.mu.Unlock()

	go c.consumeLoop()
	return nil
}

// consumeLoop handles the actual message consumption loop
func (c *Consumer) consumeLoop() {
	for {
		select {
		case <-c.stopChan:
			log.Println("Stopping message consumption...")
			return
		default:
		}

		c.mu.RLock()
		if c.conn == nil || c.conn.IsClosed() || c.channel == nil || c.channel.IsClosed() {
			c.mu.RUnlock()
			log.Println("RabbitMQ connection is not available, waiting...")
			time.Sleep(5 * time.Second)
			continue
		}
		channel := c.channel
		c.mu.RUnlock()

		// Set QoS
		err := channel.Qos(1, 0, false)
		if err != nil {
			log.Printf("Failed to set QoS: %v, retrying...", err)
			time.Sleep(5 * time.Second)
			continue
		}

		msgs, err := channel.Consume(
			c.config.QueueName,
			"",
			false, // auto-ack
			false, // exclusive
			false, // no-local
			false, // no-wait
			nil,
		)
		if err != nil {
			log.Printf("Failed to register a consumer: %v, retrying...", err)
			time.Sleep(5 * time.Second)
			continue
		}

		log.Printf("Starting to consume messages from queue: %s", c.config.QueueName)

		for {
			select {
			case <-c.stopChan:
				log.Println("Stopping message consumption...")
				return
			case msg, ok := <-msgs:
				if !ok {
					log.Println("Message channel closed, will retry consumption...")
					time.Sleep(2 * time.Second)
					break // Break inner loop to retry
				}
				go c.processMessage(msg)
			}
		}
	}
}

// processMessage processes a single message with retry logic
func (c *Consumer) processMessage(msg amqp.Delivery) {
	retryCount := GetRetryCount(msg)

	if retryCount >= c.retryConfig.MaxRetries {
		log.Printf("Max retries exceeded for message, sending to DLQ")
		if err := msg.Reject(false); err != nil {
			log.Printf("Failed to reject message after max retries: %v", err)
		}
		return
	}

	// Process the message using the handler
	err := c.handler(msg)
	if err != nil {
		log.Printf("Failed to process message (attempt %d/%d): %v", retryCount+1, c.retryConfig.MaxRetries, err)

		// Prepare retry headers
		newHeaders := amqp.Table{}
		if msg.Headers == nil {
			msg.Headers = amqp.Table{}
		}
		newHeaders["x-retry-count"] = retryCount + 1
		newHeaders["x-last-error"] = err.Error()
		newHeaders["x-last-retry"] = time.Now().Unix()

		// Reject message
		if err := msg.Reject(false); err != nil {
			log.Printf("Failed to reject message for retry: %v", err)
		}

		// Schedule retry
		c.mu.RLock()
		channel := c.channel
		c.mu.RUnlock()
		delay := CalculateRetryDelay(retryCount, c.retryConfig)
		ScheduleRetry(channel, c.config, msg.Body, newHeaders, delay)
		return
	}

	// Success - acknowledge message
	if err := msg.Ack(false); err != nil {
		log.Printf("Failed to acknowledge message: %v", err)
	}
}

// IsConnected returns true if the consumer has a valid connection
func (c *Consumer) IsConnected() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.conn != nil && !c.conn.IsClosed() && c.channel != nil && !c.channel.IsClosed()
}

// Close closes the consumer and its connections
func (c *Consumer) Close() error {
	c.mu.Lock()
	c.consuming = false
	c.mu.Unlock()

	c.stopOnce.Do(func() {
		close(c.stopChan)
	})

	time.Sleep(100 * time.Millisecond)

	c.mu.Lock()
	defer c.mu.Unlock()

	if c.channel != nil {
		if err := c.channel.Close(); err != nil {
			return err
		}
	}
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// setupConnectionRecovery sets up automatic reconnection
func (c *Consumer) setupConnectionRecovery() {
	go func() {
		for err := range c.conn.NotifyClose(make(chan *amqp.Error)) {
			if err != nil {
				log.Printf("RabbitMQ connection lost: %v, attempting to reconnect...", err)
				c.reconnect()
			}
		}
	}()

	go func() {
		for err := range c.channel.NotifyClose(make(chan *amqp.Error)) {
			if err != nil {
				log.Printf("RabbitMQ channel lost: %v, attempting to reconnect...", err)
				c.reconnect()
			}
		}
	}()
}

// reconnect attempts to reconnect to RabbitMQ
func (c *Consumer) reconnect() {
	for {
		log.Println("Attempting to reconnect to RabbitMQ...")

		c.mu.Lock()
		if c.channel != nil {
			c.channel.Close()
		}
		if c.conn != nil {
			c.conn.Close()
		}
		c.mu.Unlock()

		time.Sleep(5 * time.Second)

		url := fmt.Sprintf("amqp://%s:%s@%s:%s/%s",
			c.connConfig.Username,
			c.connConfig.Password,
			c.connConfig.Host,
			c.connConfig.Port,
			c.connConfig.VHost,
		)

		conn, err := amqp.Dial(url)
		if err != nil {
			log.Printf("Failed to reconnect: %v, retrying in 5 seconds...", err)
			continue
		}

		ch, err := conn.Channel()
		if err != nil {
			log.Printf("Failed to create channel: %v, retrying in 5 seconds...", err)
			conn.Close()
			continue
		}

		if err := c.config.SetupAllQueues(ch); err != nil {
			log.Printf("Failed to setup queues: %v, retrying in 5 seconds...", err)
			ch.Close()
			conn.Close()
			continue
		}

		c.mu.Lock()
		c.conn = conn
		c.channel = ch
		wasConsuming := c.consuming
		c.mu.Unlock()

		log.Println("Successfully reconnected to RabbitMQ")

		if wasConsuming {
			log.Println("Restarting message consumption after reconnection...")
			go func() {
				if err := c.StartConsuming(); err != nil {
					log.Printf("Failed to restart consumption after reconnection: %v", err)
				}
			}()
		}

		break
	}
}
