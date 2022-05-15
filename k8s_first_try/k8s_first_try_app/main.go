package main

import (
	"github.com/gin-gonic/gin"
	"k8s_first_try/routers"
	"log"
)

func main() {
	router := gin.Default()

	router.GET("/health", routers.HealthGET)
	router.GET("/", routers.HomeGET)
	router.GET("/:name", routers.HelloStudent)
	router.GET("/helloStudent/:name", routers.HelloStudent)

	err := router.Run(":8000")
	if err != nil {
		log.Fatalf("Error while run API: %v", err)
	}
}
