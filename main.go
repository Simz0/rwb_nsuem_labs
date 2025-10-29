package main

import (
	"context"
	"fmt"
	"labs/l0/cache"
	"labs/l0/database"
	"labs/l0/database/repository"
	"labs/l0/nats_handler"
	"labs/l0/router"
	"labs/l0/services"
	"os/signal"
	"syscall"

	"log"
	"os"

	"github.com/nats-io/nats.go"
)

func main() {
	database.ConnectDB()
	err := database.InitDatabase()

	ctx := context.Background()

	if err != nil {
		fmt.Print(err)
		log.Fatal("Error with init database")
	} else {
		fmt.Println("Init DB is successfull")
	}

	useRedis := (os.Getenv("USE_REDIS")) == "true"
	redisAddr := os.Getenv("REDIS_ADDR")

	cacheInstanse, err := cache.NewCache(useRedis, redisAddr)
	if err != nil {
		log.Fatal("Failed to create cache: ", err)
	}

	orderRepo := repository.NewGormOrderRepository()
	err = cacheInstanse.Preload(ctx, orderRepo)
	if err != nil {
		log.Printf("Warning: failed to preload orders cache: %v", err)
	}

	orderService := services.NewOrderServiceWithCache(orderRepo, cacheInstanse)
	orderHandler := nats_handler.NewOrderHandler(orderService)

	NATS_URL := os.Getenv("NATS_URL")
	if NATS_URL == "" {
		NATS_URL = nats.DefaultURL
	}

	nc, err := nats.Connect(NATS_URL)
	if err != nil {
		log.Fatal("Failed to connect NATS:", err)
	}

	defer nc.Drain()

	CHANNEL := os.Getenv("NATS_CHANNEL")
	if CHANNEL == "" {
		CHANNEL = "foo"
	}

	sub, err := nc.Subscribe(CHANNEL, orderHandler.HandlerOrderMessage())

	if err != nil {
		log.Fatal("Failed to subscribe:", err)
	}

	defer sub.Unsubscribe()

	nc.Publish(CHANNEL, []byte("Service starting and listening"))

	go func() {
		r := router.SetupRouter()
		port := os.Getenv("PORT")
		if port == "" {
			port = "8080"
		}
		log.Printf("Starting HTTP server on :%s", port)
		if err := r.Run(":" + port); err != nil {
			log.Fatal("HTTP server error: ", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	log.Println("Shutdown signal received, exiting")
}
