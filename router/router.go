package router

import (
	"labs/l0/controllers"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	useRedis := (os.Getenv("REDIS")) == "true"
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	if err := controllers.InitControllers(useRedis, redisAddr); err != nil {
		log.Fatal("Failed to init controllers: ", err)
	}

	health := r.Group("/health")
	{
		health.GET("/ping", controllers.Ping)
	}

	orders := r.Group("/orders")
	{
		orders.GET("/:uid", controllers.GetOrderByUID)
		orders.GET("/all", controllers.GetAllOrders)
	}

	return r
}
