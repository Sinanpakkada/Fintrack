package main

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Transaction represents a financial transaction
type Transaction struct {
	ID          string    `json:"id"`
	Amount      float64   `json:"amount"`
	Category    string    `json:"category"`
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
	Type        string    `json:"type"` // "income" or "expense"
}

// MonthlySummary represents a monthly financial summary
type MonthlySummary struct {
	Month          string  `json:"month"`
	Year           int     `json:"year"`
	TotalIncome    float64 `json:"totalIncome"`
	TotalExpenses  float64 `json:"totalExpenses"`
	NetAmount      float64 `json:"netAmount"`
}

var transactions []Transaction

func main() {
	r := gin.Default()

	// Enable CORS
	// Update your CORS configuration in main.go
r.Use(cors.New(cors.Config{
    AllowOrigins:     []string{"http://localhost:5173"}, // This should match your frontend URL
    AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
    AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
    ExposeHeaders:    []string{"Content-Length"},
    AllowCredentials: true,
    MaxAge:           12 * time.Hour,
}))

	// Add some test data
	addTestData()

	// Routes
	r.GET("/api/transactions", getTransactions)
	r.POST("/api/transactions", addTransaction)
	r.DELETE("/api/transactions/:id", deleteTransaction)
	r.GET("/api/summary", getMonthlySummary)

	r.Run(":8080") // Listen on port 8080
}

func addTestData() {
	transactions = append(transactions, Transaction{
		ID:          "1",
		Amount:      1500,
		Category:    "Salary",
		Description: "Monthly salary",
		Date:        time.Now().AddDate(0, 0, -5),
		Type:        "income",
	})
	transactions = append(transactions, Transaction{
		ID:          "2",
		Amount:      45.99,
		Category:    "Groceries",
		Description: "Weekly shopping",
		Date:        time.Now().AddDate(0, 0, -3),
		Type:        "expense",
	})
	transactions = append(transactions, Transaction{
		ID:          "3",
		Amount:      25.00,
		Category:    "Entertainment",
		Description: "Movie tickets",
		Date:        time.Now().AddDate(0, 0, -1),
		Type:        "expense",
	})
}

func getTransactions(c *gin.Context) {
	c.JSON(http.StatusOK, transactions)
}

func addTransaction(c *gin.Context) {
	var newTransaction Transaction
	if err := c.ShouldBindJSON(&newTransaction); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate a simple ID based on timestamp
	newTransaction.ID = time.Now().Format("20060102150405")
	
	transactions = append(transactions, newTransaction)
	c.JSON(http.StatusCreated, newTransaction)
}

func deleteTransaction(c *gin.Context) {
	id := c.Param("id")
	
	for i, transaction := range transactions {
		if transaction.ID == id {
			// Remove the transaction
			transactions = append(transactions[:i], transactions[i+1:]...)
			c.JSON(http.StatusOK, gin.H{"message": "Transaction deleted"})
			return
		}
	}
	
	c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
}

func getMonthlySummary(c *gin.Context) {
	// Group transactions by month and year
	summaryMap := make(map[string]*MonthlySummary)
	
	for _, transaction := range transactions {
		month := transaction.Date.Format("January")
		year := transaction.Date.Year()
		key := month + "-" + string(year)
		
		summary, exists := summaryMap[key]
		if !exists {
			summary = &MonthlySummary{
				Month: month,
				Year:  year,
			}
			summaryMap[key] = summary
		}
		
		if transaction.Type == "income" {
			summary.TotalIncome += transaction.Amount
		} else if transaction.Type == "expense" {
			summary.TotalExpenses += transaction.Amount
		}
		
		summary.NetAmount = summary.TotalIncome - summary.TotalExpenses
	}
	
	// Convert map to slice
	summaries := make([]MonthlySummary, 0, len(summaryMap))
	for _, summary := range summaryMap {
		summaries = append(summaries, *summary)
	}
	
	c.JSON(http.StatusOK, summaries)
}