package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	sqstypes "github.com/aws/aws-sdk-go-v2/service/sqs/types"

	"github.com/blazeisclone/vehicle-dms-inventory/events"
)

const (
	maxMessages       = 10
	waitTimeSeconds   = 20
	visibilityTimeout = 30
)

// HandlerFunc processes a single domain event. Return a non-nil error to leave
// the message in the queue; it becomes visible again after visibilityTimeout.
type HandlerFunc func(ctx context.Context, event events.DomainEvent) error

// Consumer is a long-polling SQS consumer that dispatches each message to a
// registered HandlerFunc based on event type.
type Consumer struct {
	client   *sqs.Client
	queueURL string
	handlers map[events.EventType]HandlerFunc
}

// NewConsumer constructs a Consumer from an aws.Config and a queue URL.
func NewConsumer(awsCfg aws.Config, queueURL string) *Consumer {
	return &Consumer{
		client:   sqs.NewFromConfig(awsCfg),
		queueURL: queueURL,
		handlers: make(map[events.EventType]HandlerFunc),
	}
}

// Register associates an EventType with a handler. Must be called before Run.
func (c *Consumer) Register(eventType events.EventType, fn HandlerFunc) {
	c.handlers[eventType] = fn
}

// Run starts the long-polling loop and blocks until ctx is cancelled.
func (c *Consumer) Run(ctx context.Context) error {
	log.Println("worker: consumer started, polling for events...")
	for {
		select {
		case <-ctx.Done():
			log.Println("worker: consumer shutting down")
			return nil
		default:
		}

		messages, err := c.receive(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return nil
			}
			log.Printf("worker: receive error: %v", err)
			continue
		}

		for _, msg := range messages {
			if err := c.process(ctx, msg); err != nil {
				log.Printf("worker: process error (message stays in queue): %v", err)
				continue
			}
			c.delete(ctx, msg)
		}
	}
}

func (c *Consumer) receive(ctx context.Context) ([]sqstypes.Message, error) {
	out, err := c.client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
		QueueUrl:              aws.String(c.queueURL),
		MaxNumberOfMessages:   maxMessages,
		WaitTimeSeconds:       waitTimeSeconds,
		VisibilityTimeout:     visibilityTimeout,
		MessageAttributeNames: []string{"All"},
	})
	if err != nil {
		return nil, fmt.Errorf("sqs receive: %w", err)
	}
	return out.Messages, nil
}

func (c *Consumer) process(ctx context.Context, msg sqstypes.Message) error {
	var event events.DomainEvent
	if err := json.Unmarshal([]byte(aws.ToString(msg.Body)), &event); err != nil {
		// Malformed message — delete to prevent poison-pill redelivery loop.
		log.Printf("worker: malformed message %s, deleting: %v", aws.ToString(msg.MessageId), err)
		c.delete(ctx, msg)
		return nil
	}

	handler, ok := c.handlers[event.Type]
	if !ok {
		log.Printf("worker: no handler for event type %q, discarding", event.Type)
		return nil
	}

	return handler(ctx, event)
}

func (c *Consumer) delete(ctx context.Context, msg sqstypes.Message) {
	_, err := c.client.DeleteMessage(ctx, &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(c.queueURL),
		ReceiptHandle: msg.ReceiptHandle,
	})
	if err != nil {
		log.Printf("worker: delete message error: %v", err)
	}
}
