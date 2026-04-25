package sns

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	snstypes "github.com/aws/aws-sdk-go-v2/service/sns/types"

	"github.com/blazeisclone/vehicle-dms-inventory/events"
)

// Publisher publishes domain events to an AWS SNS topic.
type Publisher struct {
	client   *sns.Client
	topicARN string
}

// New constructs a Publisher from an aws.Config and a topic ARN.
func New(awsCfg aws.Config, topicARN string) *Publisher {
	return &Publisher{
		client:   sns.NewFromConfig(awsCfg),
		topicARN: topicARN,
	}
}

// Publish JSON-encodes the DomainEvent and publishes it to the configured SNS topic.
// The event_type message attribute enables future SNS subscription filter policies.
func (p *Publisher) Publish(ctx context.Context, event events.DomainEvent) error {
	body, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("sns publisher: marshal event: %w", err)
	}

	_, err = p.client.Publish(ctx, &sns.PublishInput{
		TopicArn: aws.String(p.topicARN),
		Message:  aws.String(string(body)),
		MessageAttributes: map[string]snstypes.MessageAttributeValue{
			"event_type": {
				DataType:    aws.String("String"),
				StringValue: aws.String(string(event.Type)),
			},
		},
	})
	if err != nil {
		return fmt.Errorf("sns publisher: publish: %w", err)
	}

	return nil
}
