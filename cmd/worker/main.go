package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"

	"github.com/blazeisclone/vehicle-dms-inventory/internal/awscloud"
	"github.com/blazeisclone/vehicle-dms-inventory/worker"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("worker: no .env file, using environment variables")
	}

	awsCfg, err := awscloud.LoadFromEnv()
	if err != nil {
		log.Fatalf("worker: aws config: %v", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	sdkCfg, err := awscloud.NewAWSConfig(ctx, awsCfg)
	if err != nil {
		log.Fatalf("worker: build aws sdk config: %v", err)
	}

	topicARN, queueURL, err := awscloud.EnsureResources(ctx, sdkCfg)
	if err != nil {
		log.Fatalf("worker: init aws resources: %v", err)
	}
	log.Printf("worker: topic=%s queue=%s", topicARN, queueURL)

	c := worker.NewConsumer(sdkCfg, queueURL)
	for eventType, fn := range worker.VehicleEventHandlers() {
		c.Register(eventType, fn)
	}

	if err := c.Run(ctx); err != nil {
		log.Fatalf("worker: run: %v", err)
	}
	log.Println("worker: stopped cleanly")
}
