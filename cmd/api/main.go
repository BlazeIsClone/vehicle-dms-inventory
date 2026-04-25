package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"

	snspublisher "github.com/blazeisclone/vehicle-dms-inventory/infra/sns"
	"github.com/blazeisclone/vehicle-dms-inventory/internal/awscloud"
	"github.com/blazeisclone/vehicle-dms-inventory/internal/server"
)

func gracefulShutdown(apiServer *http.Server, done chan bool) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()

	log.Println("shutting down gracefully, press Ctrl+C again to force")
	stop()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := apiServer.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown with error: %v", err)
	}

	log.Println("Server exiting")

	done <- true
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("api: no .env file, using environment variables")
	}

	awsCfg, err := awscloud.LoadFromEnv()
	if err != nil {
		log.Fatalf("api: aws config: %v", err)
	}

	sdkCfg, err := awscloud.NewAWSConfig(context.Background(), awsCfg)
	if err != nil {
		log.Fatalf("api: build aws sdk config: %v", err)
	}

	pub := snspublisher.New(sdkCfg, awsCfg.TopicARN)
	log.Printf("api: SNS publisher ready, topic=%s", awsCfg.TopicARN)

	srv := server.NewServer(pub)

	done := make(chan bool, 1)
	go gracefulShutdown(srv, done)

	err = srv.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		panic(fmt.Sprintf("http server error: %s", err))
	}

	<-done
	log.Println("Graceful shutdown complete.")
}
