package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"

	"github.com/blazeisclone/vehicle-dms-inventory/infra/sqs"
	"github.com/blazeisclone/vehicle-dms-inventory/internal/aws"
	"github.com/blazeisclone/vehicle-dms-inventory/internal/database"
	"github.com/blazeisclone/vehicle-dms-inventory/internal/outbox"
	"github.com/blazeisclone/vehicle-dms-inventory/inventory/vehicle"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("worker: no .env file, using environment variables")
	}

	awsCfg, err := aws.LoadFromEnv()
	if err != nil {
		log.Fatalf("worker: aws config: %v", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	sdkCfg, err := aws.NewAWSConfig(ctx, awsCfg)
	if err != nil {
		log.Fatalf("worker: build aws sdk config: %v", err)
	}

	topicARN, queueURL, err := aws.EnsureResources(ctx, sdkCfg)
	if err != nil {
		log.Fatalf("worker: init aws resources: %v", err)
	}
	log.Printf("worker: topic=%s queue=%s", topicARN, queueURL)

	db := database.New()
	processed := outbox.NewProcessedStore(db.DB())

	c, err := sqs.NewConsumer(sdkCfg, queueURL, processed)
	if err != nil {
		log.Fatalf("worker: init consumer: %v", err)
	}

	for eventType, fn := range vehicle.EventHandlers() {
		c.Register(eventType, fn)
	}

	if err := c.Run(ctx); err != nil {
		log.Fatalf("worker: run: %v", err)
	}
	log.Println("worker: stopped cleanly")
}
