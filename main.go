package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"

	"challenge-goapi/bill"
	"challenge-goapi/config"
	"challenge-goapi/customer"
	"challenge-goapi/employee"
	"challenge-goapi/middleware"
	"challenge-goapi/product"

	docs "challenge-goapi/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var db = config.ConnectDB()

// @title Laundry API Documentation
// @version 1.0
// @description A documentation how to access the api's routes on your laundry application
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url https://uqonstd.xyz/
// @contact.email sdesain25@gmail.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {
	// Tulis kode kamu disini
	router := gin.Default()
	docs.SwaggerInfo.BasePath = "/api/v1"
	// Generate Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.Use(middleware.LoggerMiddleware)
	api := router.Group("/api/v1")
	{
		api.POST("/login", employee.Login)
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
		productGroup := api.Group("/products")
		{
			productGroup.GET("/", product.GetProducts)
			productGroup.GET("/:id", product.GetProduct)
			productGroup.POST("/", product.CreateProduct)
			productGroup.PUT("/:id", product.UpdateProduct)
			productGroup.DELETE("/:id", product.DeleteProduct)
		}
		billGroup := api.Group("/transactions")
		{
			billGroup.GET("/", bill.GetBills)
			billGroup.GET("/:id_bill", bill.GetBill)
			billGroup.POST("/", bill.CreateBill)
		}

	}
	defer db.Close()
	router.Run(":8080")
}
