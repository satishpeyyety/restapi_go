package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
)

// Employee represents the employee model based on the table structure
type Employee struct {
	ID       uuid.UUID `json:"id" gorm:"type:char(36);primaryKey"`
	Name     string    `json:"name" gorm:"type:text;not null"`
	Position string    `json:"position" gorm:"type:text;not null"`
	Salary   float64   `json:"salary" gorm:"type:float;not null"`
}

var db *gorm.DB

// ConnectDB initializes the TiDB connection
func ConnectDB() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, proceeding with system environment variables")
	}

	// Use environment variables to form the TiDB connection string
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?tls=true&charset=utf8mb4&parseTime=True",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)
	fmt.Println(dsn)

	// Connect to the TiDB database using GORM
	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto-migrate the Employee model to create the table if it doesn't exist
	if err := db.AutoMigrate(&Employee{}); err != nil {
		log.Fatal("Failed to auto-migrate:", err)
	}
}

// CreateEmployee handles creating a new employee
func CreateEmployee(c *gin.Context) {
	var employee Employee

	// Bind incoming JSON to employee struct
	if err := c.ShouldBindJSON(&employee); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Assign a new UUID to the employee ID
	employee.ID = uuid.New()

	// Create a new employee record
	db.Create(&employee)
	c.JSON(http.StatusOK, employee)
}

// GetEmployee handles retrieving a specific employee by ID
func GetEmployee(c *gin.Context) {
	id := c.Param("id")
	var employee Employee

	// Find employee by ID
	if err := db.First(&employee, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Employee not found"})
		return
	}

	c.JSON(http.StatusOK, employee)
}

// GetEmployees handles listing all employees
func GetEmployees(c *gin.Context) {
	var employees []Employee
	db.Find(&employees)
	c.JSON(http.StatusOK, employees)
}

// UpdateEmployee handles updating an employee's data
func UpdateEmployee(c *gin.Context) {
	id := c.Param("id")
	var employee Employee

	// Find employee by ID
	if err := db.First(&employee, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Employee not found"})
		return
	}

	// Bind updated data to employee struct
	if err := c.ShouldBindJSON(&employee); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Save changes to the database
	db.Save(&employee)
	c.JSON(http.StatusOK, employee)
}

// DeleteEmployee handles deleting an employee by ID
func DeleteEmployee(c *gin.Context) {
	id := c.Param("id")
	var employee Employee

	// Find employee by ID
	if err := db.First(&employee, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Employee not found"})
		return
	}

	// Delete the employee record
	db.Delete(&employee)
	c.JSON(http.StatusOK, gin.H{"message": "Employee deleted"})
}

func main() {
	// Initialize TiDB connection
	ConnectDB()

	// Create a new Gin router
	r := gin.Default()

	// Define RESTful routes
	r.POST("/employees", CreateEmployee)
	r.GET("/employees/:id", GetEmployee)
	r.GET("/employees", GetEmployees)
	r.PUT("/employees/:id", UpdateEmployee)
	r.DELETE("/employees/:id", DeleteEmployee)

	// Start the server
	r.Run(":8080") // The API will run on http://localhost:8080
}
