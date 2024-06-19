package bill

import (
	"challenge-goapi/entity"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

// var db = config.ConnectDB()
// type Billing struct {
// 	bill        entity.Bill
// 	employee    entity.Employee
// 	customer    entity.Customer
// 	billDetails []entity.BillDetails
// }

func GetBills(c *gin.Context) {
	startDate := c.Query("startDate")
	endDate := c.Query("endDate")
	productName := c.Query("productName")
	// var transactions []Transaction
	var bills []entity.Transaction
	var rows *sql.Rows
	var err error

	// query := "SELECT billid, entrydate, finishdate, employee, customer, totalbill FROM tx_bill"
	// query += " WHERE 1 = 1"
	query := "SELECT b.billid, b.entrydate, b.finishdate, b.employee, b.customer, b.totalbill, e.name AS employee_name, e.phonenumber AS employee_phone, e.address AS employee_address, c.name AS customer_name, c.phonenumber AS customer_phone, c.address AS customer_address"
	query += " FROM tx_bill b"
	query += " JOIN ms_employee e ON b.employee = e.id"
	query += " JOIN ms_customer c ON b.customer = c.id"
	query += " WHERE 1 = 1"

	if startDate == "" && endDate == "" && productName == "" {
		rows, err = db.Query(query)
	}
	if startDate != "" && endDate == "" && productName == "" {
		query += " AND entrydate >= $1 AND entrydate <= $1"
		rows, err = db.Query(query, startDate)
	}
	if startDate == "" && endDate != "" && productName == "" {
		query += " AND finishdate >= $1 AND finishdate <= $1"
		rows, err = db.Query(query, endDate)
	}
	if startDate != "" && endDate != "" && productName == "" {
		query += " AND entrydate >= $1 AND finishdate <= $2"
		rows, err = db.Query(query, startDate, endDate)
	}
	if startDate != "" && endDate != "" && productName != "" {
		query += " AND entrydate >= $1 AND finishdate <= $2"
		query += " AND EXISTS (SELECT 1 FROM tx_billdetails d JOIN ms_products e ON d.product = e.id WHERE d.billid = b.billid AND name ILIKE '%' || $3 || '%')"
		rows, err = db.Query(query, startDate, endDate, productName)
	}
	if startDate != "" && endDate == "" && productName != "" {
		query += " AND entrydate >= $1 AND entrydate <= $1"
		query += " AND EXISTS (SELECT 1 FROM tx_billdetails d JOIN ms_products e ON d.product = e.id WHERE d.billid = b.billid AND name ILIKE '%' || $2 || '%')"
		rows, err = db.Query(query, startDate, productName)
	}
	if startDate == "" && endDate != "" && productName != "" {
		query += " AND finishdate >= $1 AND finishdate <= $1"
		query += " AND EXISTS (SELECT 1 FROM tx_billdetails d JOIN ms_products e ON d.product = e.id WHERE d.billid = b.billid AND name ILIKE '%' || $2 || '%')"
		rows, err = db.Query(query, endDate, productName)
	}
	if startDate == "" && endDate == "" && productName != "" {
		query += " AND EXISTS (SELECT 1 FROM tx_billdetails d JOIN ms_products e ON d.product = e.id WHERE d.billid = b.billid AND name ILIKE '%' || $1 || '%')"
		rows, err = db.Query(query, productName)
	}

	// rows, err = db.Query(query, startDate, endDate, productName)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var bill = entity.Transaction{}
		var employee = entity.Employee{}
		var customer = entity.Customer{}
		err = rows.Scan(
			&bill.BillId, &bill.EntryDate, &bill.FinishDate, &employee.Id, &customer.Id, &bill.TotalBill,
			&employee.Name, &employee.PhoneNumber, &employee.Address,
			&customer.Name, &customer.PhoneNumber, &customer.Address)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		bill.Employee = employee
		bill.Customer = customer

		// Get bill details for each bill
		billDetails, err := getBillDetails(bill.BillId)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		bill.BillDetails = billDetails

		bills = append(bills, bill)
	}
	if err = rows.Err(); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// get bill details	from database based on billId

	if len(bills) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "tidak ada data bill"})
		return
	} else {
		c.JSON(200, gin.H{
			"message": "list data transactions",
			"data":    bills,
		})
	}
}

func GetBill(c *gin.Context) {

	idBill := c.Param("id_bill")
	var transaction entity.Transaction
	var bill entity.Bill
	var err error

	query := "SELECT billid, entrydate, finishdate, employee, customer, totalbill FROM tx_bill WHERE billid = $1"
	err = db.QueryRow(query, idBill).Scan(&bill.BillId, &bill.EntryDate, &bill.FinishDate, &bill.Employee, &bill.Customer, &bill.TotalBill)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "tidak ada data bill"})
		return
	} else if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	transaction.ID = bill.Id
	transaction.BillId = bill.BillId
	transaction.EntryDate = bill.EntryDate
	transaction.FinishDate = bill.FinishDate
	transaction.TotalBill = bill.TotalBill

	// Fetch bill details
	transaction.BillDetails, err = getBillDetails(bill.BillId)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Fetch employee and customer details
	query = "SELECT id, name, phonenumber, address FROM ms_employee WHERE id = $1"
	err = db.QueryRow(query, bill.Employee).Scan(&transaction.Employee.Id, &transaction.Employee.Name, &transaction.Employee.PhoneNumber, &transaction.Employee.Address)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	query = "SELECT id, name, phonenumber, address FROM ms_customer WHERE id = $1"
	err = db.QueryRow(query, bill.Customer).Scan(&transaction.Customer.Id, &transaction.Customer.Name, &transaction.Customer.PhoneNumber, &transaction.Customer.Address)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message": "data bill",
		"data":    transaction,
	})

}

// get bill details from database based on billId
func getBillDetails(billId string) ([]entity.BillDetails, error) {
	var billDetails []entity.BillDetails
	query := "SELECT billid, product, productprice, qty FROM tx_billdetails WHERE billid = $1"
	rows, err := db.Query(query, billId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var billDetail entity.BillDetails
		err = rows.Scan(&billDetail.BillId, &billDetail.Product, &billDetail.ProductPrice, &billDetail.Qty)
		if err != nil {
			return nil, err
		}
		billDetails = append(billDetails, billDetail)
	}

	return billDetails, nil
}
