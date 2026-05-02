package sqs

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-aws/sqs"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/aws/aws-sdk-go-v2/aws"

	"github.com/blazeisclone/vehicle-dms-inventory/internal/events"
	"github.com/blazeisclone/vehicle-dms-inventory/internal/outbox"
)

// Consumer wraps a Watermill SQS subscriber and router, dispatching each
// message to the registered HandlerFunc for its event type.
type Consumer struct {
	router    *message.Router
	sub       *sqs.Subscriber
	queueURL  string
	handlers  map[events.EventType]events.HandlerFunc
	processed *outbox.ProcessedStore
}

// NewConsumer constructs a Consumer from an aws.Config, queue URL, and processed-event store.
func NewConsumer(awsCfg aws.Config, queueURL string, processed *outbox.ProcessedStore) (*Consumer, error) {
	logger := watermill.NewStdLogger(false, false)

	sub, err := sqs.NewSubscriber(sqs.SubscriberConfig{
		AWSConfig:                   awsCfg,
		QueueUrlResolver:            sqs.TransparentUrlResolver{},
		DoNotCreateQueueIfNotExists: true,
	}, logger)
	if err != nil {
		return nil, fmt.Errorf("consumer: init subscriber: %w", err)
	}

	router, err := message.NewRouter(message.RouterConfig{}, logger)
	if err != nil {
		return nil, fmt.Errorf("consumer: init router: %w", err)
	}

	return &Consumer{
		router:    router,
		sub:       sub,
		queueURL:  queueURL,
		handlers:  make(map[events.EventType]events.HandlerFunc),
		processed: processed,
	}, nil
}

// Register associates an EventType with a handler. Must be called before Run.
func (c *Consumer) Register(eventType events.EventType, fn events.HandlerFunc) {
	c.handlers[eventType] = fn
}

// Run wires up the router handler and blocks until ctx is cancelled.
func (c *Consumer) Run(ctx context.Context) error {
	log.Println("worker: consumer started, polling for events...")
	c.router.AddNoPublisherHandler("vehicle-events", c.queueURL, c.sub, c.dispatch)
	return c.router.Run(ctx)
}

func (c *Consumer) dispatch(msg *message.Message) error {
	var event events.DomainEvent
	if err := json.Unmarshal(msg.Payload, &event); err != nil {
		// Malformed message — ack to prevent poison-pill redelivery loop.
		log.Printf("worker: malformed message %s, discarding: %v", msg.UUID, err)
		return nil
	}

	already, err := c.processed.IsProcessed(msg.Context(), event.ID)
	if err != nil {
		return fmt.Errorf("worker: check idempotency for event %s: %w", event.ID, err)
	}
	if already {
		log.Printf("worker: duplicate event %s (%s), skipping", event.ID, event.Type)
		return nil
	}

	handler, ok := c.handlers[event.Type]
	if !ok {
		log.Printf("worker: no handler for event type %q, discarding", event.Type)
		return nil
	}

	if err := handler(msg.Context(), event); err != nil {
		return err
	}

	if err := c.processed.MarkProcessed(msg.Context(), event.ID); err != nil {
		log.Printf("worker: mark processed failed for event %s: %v", event.ID, err)
	}

	return nil
}
