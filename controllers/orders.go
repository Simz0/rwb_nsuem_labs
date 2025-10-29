package controllers

import (
	"context"
	"log"
	"net/http"

	"labs/l0/cache"
	"labs/l0/database/repository"
	"labs/l0/services"

	"github.com/gin-gonic/gin"
)

var sharedCache cache.Cache
var sharedService *services.OrderService

func InitControllers(useRedis bool, redisAddr string) error {
	var err error
	sharedCache, err := cache.NewCache(useRedis, redisAddr)
	if err != nil {
		return err
	}

	orderRepo := repository.NewGormOrderRepository()
	sharedService = services.NewOrderServiceWithCache(orderRepo, sharedCache)

	err = sharedCache.Preload(context.Background(), orderRepo)
	if err != nil {
		log.Printf("Warning: failed to preload cache: %v", err)
	}

	return nil
}

func GetOrderByUID(c *gin.Context) {
	ctx := context.Background()

	order_uid := c.Param("uid")

	order, err := sharedService.GetOrder(ctx, order_uid)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Server unavailible"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"order": order})
}

func GetAllOrders(c *gin.Context) {
	ctx := context.Background()

	orders, err := sharedService.GetAllOrders(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve orders"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"orders": orders})
}
