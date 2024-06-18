package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"

	"challenge-goapi/config"
	"challenge-goapi/middleware"
)

var db = config.ConnectDB()

func main() {
	// Tulis kode kamu disini
	router := gin.Default()
	router.Use(middleware.LoggerMiddleware)
	api := router.Group("/api")
	{
		api.Use(middleware.AuthMiddleware)
		employee := api.Group("/employee")
		{
			employee.GET("/", func(c *gin.Context) {
				fmt.Println("employ page")
			})
		}

	}
	defer db.Close()
	router.Run(":8080")
}
