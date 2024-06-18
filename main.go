package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"

	"challenge-goapi/config"
	"challenge-goapi/customer"
	"challenge-goapi/employee"
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
		employeeGroup := api.Group("/employees")
		{
			employeeGroup.GET("/", employee.GetEmployees)
			employeeGroup.GET("/:id", employee.GetEmployee)
			employeeGroup.POST("/", employee.CreateEmployee)
			employeeGroup.PUT("/:id", employee.UpdateEmployee)
			employeeGroup.DELETE("/:id", employee.DeleteEmployee)
		}
		customerGroup := api.Group("/customers")
		{
			customerGroup.GET("/", customer.GetCustomers)
			customerGroup.GET("/:id", customer.GetCustomer)
			customerGroup.POST("/", customer.CreateCustomer)
			customerGroup.PUT("/:id", customer.UpdateCustomer)
			customerGroup.DELETE("/:id", customer.DeleteCustomer)
		}

	}
	router.POST("/login", employee.Login)
	defer db.Close()
	router.Run(":8080")
}
