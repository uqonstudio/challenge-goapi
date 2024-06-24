package customer

import (
	"challenge-goapi/config"
	"challenge-goapi/entity"
	"challenge-goapi/lib"
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var db = config.ConnectDB()

// @Summary Get Customers
// @Description Retrieve a list of customers
// @Accept json
// @Produce json
// @Param name query string false "Filter by customer name"
// @Success 200
// @Failure 400
// @Tags Customers
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization" default(Bearer <Add access token here>)
// @Router /customers [get]
func GetCustomers(c *gin.Context) {
	var customers []entity.Customer
	name := c.Query("name")
	var rows *sql.Rows
	var err error

	query := "SELECT id, name, address, phoneNumber FROM ms_customer"

	if name == "" {
		rows, err = db.Query(query)
	} else {
		query += " WHERE name = $1;"
		rows, err = db.Query(query, name)
	}

	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	defer rows.Close()

	for rows.Next() {
		var cust entity.Customer
		err = rows.Scan(&cust.Id, &cust.Name, &cust.Address, &cust.PhoneNumber)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		customers = append(customers, cust)
	}

	if err = rows.Err(); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if len(customers) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "belum ada data customer"})
		return
	} else {
		c.JSON(200, gin.H{
			"message": "data customer",
			"data":    customers,
		})
	}
}

// @Summary Get Customer by ID
// @Description Retrieve a customer by ID
// @Accept json
// @Produce json
// @Param id path int true "Customer ID"
// @Success 200
// @Failure 400 {string} string "error"
// @Tags Customers
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization" default(Bearer <Add access token here>)
// @Router /customers/{id} [get]
func GetCustomer(c *gin.Context) {
	cid, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
	}

	query := "SELECT id, name, address, phoneNumber FROM ms_customer WHERE id = $1;"
	var cust entity.Customer
	err = db.QueryRow(query, cid).Scan(&cust.Id, &cust.Name, &cust.Address, &cust.PhoneNumber)
	if err != nil && err == sql.ErrNoRows {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "data customer",
		"data":    cust,
	})
}

// @Summary Create Customer
// @Description Create a new customer
// @Accept json
// @Produce json
// @Param customer body entity.Customer true "Customer data"
// @Success 200
// @Failure 400 {string} string "error"
// @Tags Customers
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization" default(Bearer <Add access token here>)
// @Router /customers [post]
func CreateCustomer(c *gin.Context) {
	var customer entity.Customer
	var err error
	if err = c.ShouldBindJSON(&customer); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// validasi customers
	if err = validateCustomer(&customer); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	query := "INSERT INTO ms_customer (name, address, phoneNumber) VALUES ($1, $2, $3) RETURNING id;"
	err = db.QueryRow(query, customer.Name, customer.Address, customer.PhoneNumber).Scan(&customer.Id)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "customer successfully created",
		"data":    customer,
	})

}

// @Summary Update Customer
// @Description Update an existing customer
// @Accept json
// @Produce json
// @Param id path int true "Customer ID"
// @Param customer body entity.Customer true "Customer data"
// @Success 200
// @Failure 400 {string} string "error"
// @Tags Customers
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization" default(Bearer <Add access token here>)
// @Router /customers/{id} [put]
func UpdateCustomer(c *gin.Context) {
	var err error
	var customer entity.Customer
	cid, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err = c.ShouldBind(&customer)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err = validateCustomer(&customer); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	query := "UPDATE ms_customer SET name = $1, address = $2, phonenumber = $3 WHERE id = $4;"
	_, err = db.Exec(query, customer.Name, customer.Address, customer.PhoneNumber, cid)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	customer.Id = cid
	c.JSON(http.StatusOK, gin.H{
		"message": "customer successfully updated",
		"data":    customer,
	})
}

// @Summary Delete Customer
// @Description Delete an existing customer
// @Accept json
// @Produce json
// @Param id path int true "Customer ID"
// @Success 200 {string} string "customer successfully deleted"
// @Failure 400 {string} string "error"
// @Tags Customers
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization" default(Bearer <Add access token here>)
// @Router /customers/{id} [delete]
func DeleteCustomer(c *gin.Context) {
	var err error
	cid, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	query := "DELETE FROM ms_customer WHERE id = $1;"
	_, err = db.Exec(query, cid)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "customer successfully deleted",
		"data":    "ok",
	})
}

func validateCustomer(customer *entity.Customer) error {
	var err error
	if err = lib.ValidateString(customer.Name); err != nil {
		return err
	}
	if err = lib.ValidateString(customer.Address); err != nil {
		return err
	}
	if err = lib.ValidatePhoneNumber(customer.PhoneNumber); err != nil {
		return err
	}
	return nil
}
