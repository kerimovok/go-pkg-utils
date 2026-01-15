package queue

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Producer is a RabbitMQ producer with automatic reconnection
type Producer struct {
	conn       *amqp.Connection
	channel    *amqp.Channel
	mu         sync.RWMutex
	config     *Config
	connConfig ConnectionConfig
}

// NewProducer creates a new RabbitMQ producer with automatic reconnection
func NewProducer(connConfig ConnectionConfig, queueConfig *Config) (*Producer, error) {
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

	producer := &Producer{
		conn:       conn,
		channel:    ch,
		config:     queueConfig,
		connConfig: connConfig,
	}

	producer.setupConnectionRecovery()

	return producer, nil
}

// Publish publishes a message to the queue
func (p *Producer) Publish(ctx context.Context, body []byte, headers amqp.Table) error {
	// Check connection health before publishing
	p.mu.RLock()
	if p.conn == nil || p.conn.IsClosed() || p.channel == nil || p.channel.IsClosed() {
		p.mu.RUnlock()
		return fmt.Errorf("RabbitMQ connection is not available")
	}
	channel := p.channel
	p.mu.RUnlock()

	// Use context with timeout if not provided
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
	}

	err := channel.PublishWithContext(ctx,
		p.config.ExchangeName, // exchange
		p.config.RoutingKey,   // routing key
		false,                 // mandatory
		false,                 // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			Headers:      headers,
			DeliveryMode: amqp.Persistent, // Make messages persistent
		})

	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	return nil
}

// IsConnected returns true if the producer has a valid connection
func (p *Producer) IsConnected() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.conn != nil && !p.conn.IsClosed() && p.channel != nil && !p.channel.IsClosed()
}

// Close closes the producer and its connections
func (p *Producer) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.channel != nil {
		if err := p.channel.Close(); err != nil {
			return err
		}
	}
	if p.conn != nil {
		return p.conn.Close()
	}
	return nil
}

// setupConnectionRecovery sets up automatic reconnection
func (p *Producer) setupConnectionRecovery() {
	go func() {
		for err := range p.conn.NotifyClose(make(chan *amqp.Error)) {
			if err != nil {
				log.Printf("RabbitMQ connection lost: %v, attempting to reconnect...", err)
				p.reconnect()
			}
		}
	}()

	go func() {
		for err := range p.channel.NotifyClose(make(chan *amqp.Error)) {
			if err != nil {
				log.Printf("RabbitMQ channel lost: %v, attempting to reconnect...", err)
				p.reconnect()
			}
		}
	}()
}

// reconnect attempts to reconnect to RabbitMQ
func (p *Producer) reconnect() {
	for {
		log.Println("Attempting to reconnect to RabbitMQ...")

		p.mu.Lock()
		if p.channel != nil {
			p.channel.Close()
		}
		if p.conn != nil {
			p.conn.Close()
		}
		p.mu.Unlock()

		time.Sleep(5 * time.Second)

		url := fmt.Sprintf("amqp://%s:%s@%s:%s/%s",
			p.connConfig.Username,
			p.connConfig.Password,
			p.connConfig.Host,
			p.connConfig.Port,
			p.connConfig.VHost,
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

		if err := p.config.SetupAllQueues(ch); err != nil {
			log.Printf("Failed to setup queues: %v, retrying in 5 seconds...", err)
			ch.Close()
			conn.Close()
			continue
		}

		p.mu.Lock()
		p.conn = conn
		p.channel = ch
		p.mu.Unlock()

		log.Println("Successfully reconnected to RabbitMQ")
		break
	}
}
