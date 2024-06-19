package bill

import (
	"challenge-goapi/config"
	"challenge-goapi/entity"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type CartDetails struct {
	BillId    string `json:"billId"`
	ProductId int    `json:"productId"`
	Qty       int    `json:"qty"`
	Price     int    `json:"price"`
	SubTotal  int    `json:"subTotal"`
}

type Cart struct {
	BillId      string        `json:"billId"`
	EntryDate   time.Time     `json:"entryDate"`
	FinishDate  time.Time     `json:"finishDate"`
	EmployeeId  int           `json:"employee"`
	CustomerId  int           `json:"customerId"`
	BillDetails []CartDetails `json:"billDetails"`
	TotalBill   int           `json:"totalBill"`
}

var db = config.ConnectDB()

func CreateBill(c *gin.Context) {
	var bill Cart
	user, err := validateToken(c)
	if err != nil {
		c.JSON(401, gin.H{"error": err.Error()})
		return
	}
	bill, err = performTransaction(user, c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "bill successfully created",
		"data":    bill,
	})
}

func validateToken(c *gin.Context) (entity.Employee, error) {
	authHeader := c.GetHeader("Authorization")
	tokenString := authHeader[7:] // remove "Bearer" from
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte("@Enigma2024"), nil
	})
	if err != nil {
		return entity.Employee{}, fmt.Errorf("invalid or expired token")
	}
	claims := token.Claims.(jwt.MapClaims)
	userId := int(claims["uid"].(float64))

	return entity.Employee{
		Id: userId,
	}, nil
}

func performTransaction(user entity.Employee, c *gin.Context) (Cart, error) {
	var err error
	var billId string
	var cart Cart

	if err = c.ShouldBind(&cart); err != nil {
		return cart, err
	}

	fmt.Printf("\nCart : %v\n", cart)

	cart.EmployeeId = user.Id // employee Id

	// Perform validation on cart data
	for _, cartDetails := range cart.BillDetails {
		// Check if product exists
		productExists, err := productExists(cartDetails.ProductId)
		if err != nil {
			return cart, err
		}
		if !productExists {
			return cart, fmt.Errorf("product with ID %d does not exist", cartDetails.ProductId)
		}

		// Validate quantity
		if cartDetails.Qty <= 0 {
			return cart, fmt.Errorf("quantity must be greater than 0")
		}
	}

	// transaction begin
	tx, err := db.Begin()
	if err != nil {
		return cart, err
	}
	defer func() {
		if err != nil {
			fmt.Println("transaction rollback")
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	// count bill return bill number
	billNumber, err := billCount(tx)
	if err != nil {
		return cart, err
	}
	timenow := time.Now()
	cart.EntryDate = timenow
	// create billId
	billId = fmt.Sprintf("EL%d%d-%d%d%d%d%d-%d", cart.EmployeeId, cart.CustomerId, timenow.Year(), timenow.YearDay(), timenow.Hour(), timenow.Minute(), timenow.Second(), billNumber)
	cart.BillId = billId

	// get prices for each item in the cart
	for idx, cartDetails := range cart.BillDetails {
		var subtotal = 0
		price, err := getProduct(cartDetails.ProductId, tx)
		if err != nil {
			return cart, err
		}
		subtotal = price * cartDetails.Qty
		cart.BillDetails[idx].BillId = billId
		cart.BillDetails[idx].Price = price
		cart.BillDetails[idx].SubTotal = subtotal

		cart.TotalBill += subtotal
	}

	// insert bill
	bill := entity.Bill{
		BillId:     cart.BillId,
		EntryDate:  cart.EntryDate,
		FinishDate: cart.FinishDate,
		Employee:   cart.EmployeeId,
		Customer:   cart.CustomerId,
		TotalBill:  cart.TotalBill,
	}
	if err = insertBill(&bill, tx); err != nil {
		return cart, err
	}

	// insert bill details
	for _, cartDetails := range cart.BillDetails {

		billDetails := entity.BillDetails{
			BillId:       cartDetails.BillId,
			Product:      cartDetails.ProductId,
			ProductPrice: cartDetails.Price,
			Qty:          cartDetails.Qty,
		}

		if err = insertBillDetails(&billDetails, tx); err != nil {
			return cart, err
		}
	}

	return cart, nil

}

func billCount(tx *sql.Tx) (int, error) {
	countRows := "SELECT COUNT(DISTINCT billid) + 1 AS billnumber FROM tx_bill;"
	billNumber := 0
	err := tx.QueryRow(countRows).Scan(&billNumber)
	if err != nil {
		return 0, err
	}
	return billNumber, nil
}

func insertBill(bill *entity.Bill, tx *sql.Tx) error {
	insertBillRequest := "INSERT INTO tx_bill (billid, entrydate, finishdate, employee, customer, totalbill ) VALUES ($1, $2, $3, $4, $5, $6);"
	_, err := tx.Exec(insertBillRequest, bill.BillId, bill.EntryDate, bill.FinishDate, bill.Employee, bill.Customer, bill.TotalBill)
	if err != nil {
		return err
	} else {
		return nil
	}
}

func insertBillDetails(billDetails *entity.BillDetails, tx *sql.Tx) error {
	query := "INSERT INTO tx_billdetails (billid, product, qty, productprice ) VALUES($1, $2, $3, $4);"
	_, err := tx.Exec(query, billDetails.BillId, billDetails.Product, billDetails.Qty, billDetails.ProductPrice)
	if err != nil {
		return err
	} else {
		return nil
	}
}

func getProduct(productId int, tx *sql.Tx) (int, error) {
	var price int
	query := "SELECT price FROM ms_products WHERE id = $1;"
	err := tx.QueryRow(query, productId).Scan(&price)
	if err != nil && err == sql.ErrNoRows {
		return 0, err
	}
	// fmt.Printf("\nPrice : %d\n", price)
	return price, nil
}

func productExists(productId int) (bool, error) {
	query := "SELECT COUNT(*) FROM ms_products WHERE id = $1;"
	var count int
	err := db.QueryRow(query, productId).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
