package product

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

// @Summary Get all products
// @Description Retrieve a list of all products
// @Produce json
// @Param name query string false "Product name"
// @Security ApiKeyAuth
// @Success 200
// @Failure 400 {string} string "error"
// @Tags Products
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization" default(Bearer <Add access token here>)
// @Router /products [get]
func GetProducts(c *gin.Context) {
	var products []entity.Product
	var rows *sql.Rows
	var err error
	name := c.Query("name")

	query := "SELECT id, name, unit, price FROM ms_products"
	if name == "" {
		rows, err = db.Query(query)
	} else {
		query += " WHERE ILIKE '%' || $1 || '%'"
	}

	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	defer rows.Close()

	for rows.Next() {
		var prd entity.Product
		err = rows.Scan(&prd.Id, &prd.Name, &prd.Unit, &prd.Price)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		products = append(products, prd)
	}

	if err = rows.Err(); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if len(products) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "product.name tidak ditemukan"})
		return
	} else {
		c.JSON(200, gin.H{
			"message": "data product",
			"data":    products,
		})
	}

}

// @Summary Get product by id
// @Description Retrieve a product by id
// @Produce json
// @Param id path int true "Product id"
// @Success 200
// @Failure 400
// @Tags Products
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization" default(Bearer <Add access token here>)
// @Router /products/{id} [get]
func GetProduct(c *gin.Context) {
	var err error
	pid, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	var product entity.Product
	query := "SELECT id, name, price, unit FROM ms_products WHERE id = $1"
	err = db.QueryRow(query, pid).Scan(&product.Id, &product.Name, &product.Price, &product.Unit)
	if err != nil && err == sql.ErrNoRows {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "data product",
		"data":    product,
	})
}

// @Summary Create product
// @Description Create a new product
// @Accept  json
// @Produce  json
// @Param product body entity.Product true "Product data"
// @Success 200
// @Failure 400
// @Tags Products
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization" default(Bearer <Add access token here>)
// @Router /products [post]
func CreateProduct(c *gin.Context) {
	var product entity.Product
	var err error
	if err = c.ShouldBind(&product); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if err = ValidateProduct(&product); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	query := "INSERT INTO ms_products(name, price, unit) VALUES($1, $2, $3) RETURNING id"
	err = db.QueryRow(query, product.Name, product.Price, product.Unit).Scan(&product.Id)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "product successfully created",
		"data":    product,
	})
}

// @Summary Update product
// @Description Update a product
// @Accept  json
// @Produce  json
// @Param id path int true "Product id"
// @Param product body entity.Product true "Product data"
// @Success 200
// @Failure 400
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization" default(Bearer <Add access token here>)
// @Tags Products
// @Router /products/{id} [put]
func UpdateProduct(c *gin.Context) {
	var err error
	var product entity.Product
	pid, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if err = c.ShouldBind(&product); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if err = ValidateProduct(&product); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	query := "UPDATE ms_products SET name = $1, price = $2, unit = $3 WHERE id = $4 "
	_, err = db.Exec(query, product.Name, product.Price, product.Unit, pid)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "product successfully updated",
		"data":    "ok",
	})
}

// @Summary Delete product
// @Description Delete a product
// @Produce json
// @Param id path int true "Product id"
// @Success 200
// @Failure 400
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization" default(Bearer <Add access token here>)
// @Tags Products
// @Router /products/{id} [delete]
func DeleteProduct(c *gin.Context) {
	var err error
	pid, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	query := "DELETE FROM ms_products WHERE id = $1"
	_, err = db.Exec(query, pid)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "product successfully deleted",
		"data":    "ok",
	})
}

func ValidateProduct(product *entity.Product) error {
	var err error
	if err = lib.ValidateString(product.Name); err != nil {
		return err
	}
	if err = lib.ValidateUnit(product.Unit); err != nil {
		return err
	}
	if err = lib.ValidatePrice(product.Price); err != nil {
		return err
	}
	return nil
}
