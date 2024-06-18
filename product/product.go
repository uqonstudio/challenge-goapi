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

func GetProducts(c *gin.Context) {
	var products []entity.Product
	var rows *sql.Rows
	var err error
	name := c.Query("name")

	query := "SELECT id, name, unit, price FROM ms_product"
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

func GetProduct(c *gin.Context) {
	var err error
	pid, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	var product entity.Product
	query := "SELECT id, name, price, unit FROM ms_product WHERE id = $1"
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
	if err = lib.ValidateString(product.Unit); err != nil {
		return err
	}
	if err = lib.ValidatePrice(product.Price); err != nil {
		return err
	}
	return nil
}
