package awscloud

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
)

// Config holds resolved AWS configuration for this service.
type Config struct {
	Region      string
	EndpointURL string // empty means use real AWS
	TopicARN    string
	QueueURL    string
}

// LoadFromEnv reads AWS configuration from environment variables.
// Required: AWS_REGION, SNS_TOPIC_ARN, SQS_QUEUE_URL.
// Optional: AWS_ENDPOINT_URL (LocalStack), AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY.
func LoadFromEnv() (*Config, error) {
	region := os.Getenv("AWS_REGION")
	if region == "" {
		return nil, fmt.Errorf("awscloud: AWS_REGION is required")
	}
	topicARN := os.Getenv("SNS_TOPIC_ARN")
	if topicARN == "" {
		return nil, fmt.Errorf("awscloud: SNS_TOPIC_ARN is required")
	}
	queueURL := os.Getenv("SQS_QUEUE_URL")
	if queueURL == "" {
		return nil, fmt.Errorf("awscloud: SQS_QUEUE_URL is required")
	}
	return &Config{
		Region:      region,
		EndpointURL: os.Getenv("AWS_ENDPOINT_URL"),
		TopicARN:    topicARN,
		QueueURL:    queueURL,
	}, nil
}

// NewAWSConfig builds an aws.Config ready for SNS and SQS clients.
// When EndpointURL is set, all service calls are routed to that endpoint (LocalStack).
func NewAWSConfig(ctx context.Context, cfg *Config) (aws.Config, error) {
	opts := []func(*config.LoadOptions) error{
		config.WithRegion(cfg.Region),
	}

	keyID := os.Getenv("AWS_ACCESS_KEY_ID")
	secret := os.Getenv("AWS_SECRET_ACCESS_KEY")
	if keyID != "" && secret != "" {
		opts = append(opts, config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(keyID, secret, ""),
		))
	}

	awsCfg, err := config.LoadDefaultConfig(ctx, opts...)
	if err != nil {
		return aws.Config{}, fmt.Errorf("awscloud: load config: %w", err)
	}

	if cfg.EndpointURL != "" {
		awsCfg.BaseEndpoint = aws.String(cfg.EndpointURL)
	}

	return awsCfg, nil
}
