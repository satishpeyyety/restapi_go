// api_test.go
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

// TestResult struct to store results of each test case
type TestResult struct {
	Name   string
	Passed bool
	Error  string
}

var testResults []TestResult

// Test setup function to initialize the Gin router
func setupRouter() *gin.Engine {
	r := gin.Default()
	r.POST("/employees", CreateEmployee)
	r.GET("/employees/:id", GetEmployee)
	r.GET("/employees", GetEmployees)
	r.PUT("/employees/:id", UpdateEmployee)
	r.DELETE("/employees/:id", DeleteEmployee)
	return r
}

// Helper function to record test results
func recordResult(name string, passed bool, errMsg string) {
	testResults = append(testResults, TestResult{Name: name, Passed: passed, Error: errMsg})
}

// Test cases for the APIs
func TestCreateEmployee(t *testing.T) {
	ConnectDB()
	r := setupRouter()

	employee := Employee{
		Name:     "John Doe",
		Position: "Developer",
		Salary:   60000,
	}
	jsonData, _ := json.Marshal(employee)

	req, _ := http.NewRequest("POST", "/employees", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	passed := w.Code == http.StatusOK
	recordResult("TestCreateEmployee", passed, "")
	if !passed {
		t.Errorf("Expected status code 200, got %v", w.Code)
	}
}

func TestGetEmployee(t *testing.T) {
	ConnectDB()
	r := setupRouter()

	employee := Employee{
		ID:       uuid.New(),
		Name:     "Jane Doe",
		Position: "Manager",
		Salary:   80000,
	}
	db.Create(&employee)

	req, _ := http.NewRequest("GET", fmt.Sprintf("/employees/%s", employee.ID.String()), nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	passed := w.Code == http.StatusOK
	recordResult("TestGetEmployee", passed, "")
	if !passed {
		t.Errorf("Expected status code 200, got %v", w.Code)
	}
}

func TestGetEmployees(t *testing.T) {
	ConnectDB()
	r := setupRouter()

	req, _ := http.NewRequest("GET", "/employees", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	passed := w.Code == http.StatusOK
	recordResult("TestGetEmployees", passed, "")
	if !passed {
		t.Errorf("Expected status code 200, got %v", w.Code)
	}
}

func TestUpdateEmployee(t *testing.T) {
	ConnectDB()
	r := setupRouter()

	employee := Employee{
		ID:       uuid.New(),
		Name:     "Jake Doe",
		Position: "Analyst",
		Salary:   70000,
	}
	db.Create(&employee)

	employee.Name = "Jake Updated"
	jsonData, _ := json.Marshal(employee)

	req, _ := http.NewRequest("PUT", fmt.Sprintf("/employees/%s", employee.ID.String()), bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	passed := w.Code == http.StatusOK
	recordResult("TestUpdateEmployee", passed, "")
	if !passed {
		t.Errorf("Expected status code 200, got %v", w.Code)
	}
}

func TestDeleteEmployee(t *testing.T) {
	ConnectDB()
	r := setupRouter()

	employee := Employee{
		ID:       uuid.New(),
		Name:     "Mark Doe",
		Position: "Designer",
		Salary:   55000,
	}
	db.Create(&employee)

	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/employees/%s", employee.ID.String()), nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	passed := w.Code == http.StatusOK
	recordResult("TestDeleteEmployee", passed, "")
	if !passed {
		t.Errorf("Expected status code 200, got %v", w.Code)
	}
}

// Print summary of test results
func printSummary() {
	fmt.Println("\nTest Summary:")
	for _, result := range testResults {
		if result.Passed {
			fmt.Printf("✓ %s passed\n", result.Name)
		} else {
			fmt.Printf("✗ %s failed: %s\n", result.Name, result.Error)
		}
	}
}

func TestMain(m *testing.M) {
	// Run the tests
	code := m.Run()
	printSummary() // Print the summary of test results
	os.Exit(code)
}
