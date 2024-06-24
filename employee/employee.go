package employee

import (
	"challenge-goapi/config"
	"challenge-goapi/entity"
	"challenge-goapi/lib"
	"net/http"
	"strconv"

	"database/sql"

	"github.com/gin-gonic/gin"
)

var db = config.ConnectDB()

// @Summary Get employee by ID
// @Description Retrieve employee data by ID
// @ID id
// @Produce json
// @Param id path int true "Employee ID"
// @Success 200 {object} entity.Employee "Employee data"
// @Failure 400
// @Tags Employees
// @Security Bearer
// @Router /employees/{id} [get]
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization" default(Bearer <Add access token here>)
func GetEmployee(c *gin.Context) {
	// Tulis kode kamu disini
	var err error
	var employee entity.Employee
	uid, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	query := "SELECT id, name, email, address, phoneNumber, department FROM ms_employee WHERE id = $1;"
	err = db.QueryRow(query, uid).Scan(&employee.Id, &employee.Name, &employee.Email, &employee.Address, &employee.PhoneNumber, &employee.Department)
	if err != nil && err == sql.ErrNoRows {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "data employee",
		"data":    employee,
	})
}

// @Summary Create a new employee
// @Description Add a new employee to the database
// @Accept json
// @Produce json
// @Param employee body entity.Employee true "Employee data"
// @Success 200 {array} entity.Employee "Employee data"
// @Failure 400
// @Tags Employees
// @Router /employees [post]
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization" default(Bearer <Add access token here>)
func CreateEmployee(c *gin.Context) {
	// Tulis kode kamu disini
	var err error
	var employee entity.Employee
	err = c.ShouldBind(&employee)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err = validateEmployee(&employee); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	employee.Password = lib.HashMD5(employee.Password)

	query := "INSERT INTO ms_employee (name, email, password, address, phonenumber, department) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id;"
	err = db.QueryRow(query, employee.Name, employee.Email, employee.Password, employee.Address, employee.PhoneNumber, employee.Department).Scan(&employee.Id)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "employee successfully created",
		"data":    employee,
	})
}

// @Summary Update an existing employee
// @Description Update the details of an existing employee in the database
// @Accept json
// @Produce json
// @Param id path int true "Employee ID"
// @Param employee body entity.Employee true "Employee data"
// @Success 200 {object} entity.Employee "Employee data"
// @Failure 400
// @Tags Employees
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization" default(Bearer <Add access token here>)
// @Router /employees/{id} [put]
func UpdateEmployee(c *gin.Context) {
	// Tulis kode kamu disini
	var employee entity.Employee
	uid, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	err = c.ShouldBind(&employee)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	// validate input parameters
	if err = validateEmployee(&employee); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	query := "UPDATE ms_employee SET name = $1, email = $2, address = $3, phonenumber = $4, department = $5, password = $6 WHERE id = $7;"
	_, err = db.Exec(query, employee.Name, employee.Email, employee.Address, employee.PhoneNumber, employee.Department, employee.Password, uid)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "employee successfully updated",
		"data":    employee,
	})
}

// @Summary Delete an employee
// @Description Remove an employee from the database
// @Accept json
// @Produce json
// @Param id path int true "Employee ID"
// @Success 200 {object} entity.Employee "Employee data"
// @Failure 400
// @Tags Employees
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization" default(Bearer <Add access token here>)
// @Router /employees/{id} [delete]
func DeleteEmployee(c *gin.Context) {
	// Tulis kode kamu disini
	uid, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	query := "DELETE FROM ms_employee WHERE id = $1;"
	_, err = db.Exec(query, uid)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "employee successfully deleted",
		"data":    "ok",
	})
}

// @Summary List all employees
// @Description Retrieve a list of all employees
// @Accept json
// @Produce json
// @Param name query string false "Filter by employee name"
// @Success 200 {object} entity.Employee "Employee data"
// @Failure 400
// @Tags Employees
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization" default(Bearer <Add access token here>)
// @Router /employees [get]
func GetEmployees(c *gin.Context) {
	var employees []entity.Employee
	name := c.Query("name")
	var rows *sql.Rows
	var err error

	query := "SELECT id, name, email, address, phoneNumber, department FROM ms_employee"
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
		var emp entity.Employee
		err = rows.Scan(&emp.Id, &emp.Name, &emp.Email, &emp.Address, &emp.PhoneNumber, &emp.Department)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		employees = append(employees, emp)
	}
	if err = rows.Err(); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if len(employees) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "employee.name tidak ditemukan"})
		return
	} else {
		c.JSON(200, gin.H{
			"message": "data employee",
			"data":    employees,
		})
	}

}

// @Summary Login an employee
// @Description Authenticate an employee and generate a token exp:. email : pegi@gmail.com, password: @Emka123
// @Accept json
// @Produce json
// @Param employee body entity.Login true "Login data"
// @Success 200 {object} entity.Employee "token"
// @Failure 400
// @Tags Auth Token
// @Router /login [post]
func Login(c *gin.Context) {
	var err error
	var employee entity.Employee
	if err = c.ShouldBindJSON(&employee); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	employee.Password = lib.HashMD5(employee.Password)

	query := "SELECT id, name, email, address, phoneNumber, department FROM ms_employee WHERE email = $1 AND password = $2; "
	err = db.QueryRow(query, &employee.Email, &employee.Password).Scan(&employee.Id, &employee.Name, &employee.Email, &employee.Address, &employee.PhoneNumber, &employee.Department)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// generate token
	token, err := lib.GenerateToken(&employee)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"token": token, "user": employee})
}

func validateEmployee(emp *entity.Employee) error {
	var err error
	if err = lib.ValidateString(emp.Name); err != nil {
		return err
	}
	if err = lib.ValidateEmail(emp.Email); err != nil {
		return err
	}
	if err = lib.ValidateString(emp.Address); err != nil {
		return err
	}
	if err = lib.ValidatePhoneNumber(emp.PhoneNumber); err != nil {
		return err
	}
	if err = lib.ValidateDepartment(emp.Department); err != nil {
		return err
	}
	if err = lib.ValidatePassword(emp.Password); err != nil {
		return nil
	}
	return nil
}
