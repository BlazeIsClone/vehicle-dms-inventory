package sns

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ThreeDotsLabs/watermill"
	sns "github.com/ThreeDotsLabs/watermill-aws/sns"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/aws/aws-sdk-go-v2/aws"

	"github.com/blazeisclone/vehicle-dms-inventory/internal/events"
)

type Publisher struct {
	pub      *sns.Publisher
	topicARN string
}

func New(conf aws.Config, topicARN string) (*Publisher, error) {
	pub, err := sns.NewPublisher(sns.PublisherConfig{
		AWSConfig:                   conf,
		TopicResolver:               sns.TransparentTopicResolver{},
		DoNotCreateTopicIfNotExists: true,
	}, watermill.NewStdLogger(false, false))
	if err != nil {
		return nil, fmt.Errorf("sns publisher: %w", err)
	}

	return &Publisher{pub: pub, topicARN: topicARN}, nil
}

func (p *Publisher) Publish(_ context.Context, event events.DomainEvent) error {
	body, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("sns publisher: marshal: %w", err)
	}

	msg := message.NewMessage(event.ID, body)

	msg.Metadata.Set("event_type", string(event.Type))
	if err := p.pub.Publish(p.topicARN, msg); err != nil {
		return fmt.Errorf("sns publisher: publish: %w", err)
	}

	return nil
}
