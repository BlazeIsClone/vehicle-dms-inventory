package aws

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	sqstypes "github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

const (
	vehicleTopicName = "vehicle-events"
	vehicleQueueName = "vehicle-events-queue"
)

// EnsureResources idempotently creates the SNS topic, SQS queue, and SNS→SQS
// subscription. Safe to call on every startup — CreateTopic and CreateQueue are
// idempotent per the AWS/LocalStack spec.
func EnsureResources(ctx context.Context, awsCfg aws.Config) (topicARN, queueURL string, err error) {
	snsClient := sns.NewFromConfig(awsCfg)
	sqsClient := sqs.NewFromConfig(awsCfg)

	topicOut, err := snsClient.CreateTopic(ctx, &sns.CreateTopicInput{
		Name: aws.String(vehicleTopicName),
	})
	if err != nil {
		return "", "", fmt.Errorf("create sns topic: %w", err)
	}
	topicARN = aws.ToString(topicOut.TopicArn)
	log.Printf("awscloud: SNS topic ready: %s", topicARN)

	queueOut, err := sqsClient.CreateQueue(ctx, &sqs.CreateQueueInput{
		QueueName: aws.String(vehicleQueueName),
	})
	if err != nil {
		return "", "", fmt.Errorf("create sqs queue: %w", err)
	}
	queueURL = aws.ToString(queueOut.QueueUrl)
	log.Printf("awscloud: SQS queue ready: %s", queueURL)

	attrOut, err := sqsClient.GetQueueAttributes(ctx, &sqs.GetQueueAttributesInput{
		QueueUrl:       aws.String(queueURL),
		AttributeNames: []sqstypes.QueueAttributeName{sqstypes.QueueAttributeNameQueueArn},
	})
	if err != nil {
		return "", "", fmt.Errorf("get queue arn: %w", err)
	}
	queueARN := attrOut.Attributes[string(sqstypes.QueueAttributeNameQueueArn)]

	_, err = snsClient.Subscribe(ctx, &sns.SubscribeInput{
		TopicArn: aws.String(topicARN),
		Protocol: aws.String("sqs"),
		Endpoint: aws.String(queueARN),
		Attributes: map[string]string{
			"RawMessageDelivery": "true",
		},
	})
	if err != nil {
		return "", "", fmt.Errorf("subscribe sqs to sns: %w", err)
	}
	log.Printf("awscloud: SQS queue subscribed to SNS topic")

	return topicARN, queueURL, nil
}
