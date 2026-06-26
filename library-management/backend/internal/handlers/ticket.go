package handlers

import (
    "net/http"
    "library-management/internal/database"
    "github.com/gin-gonic/gin"
)

func GetTicket(c *gin.Context) {
    userID := c.GetString("userID")
    
    var stats struct {
        TotalBooksRead int `json:"total_books_read"`
        ActiveLoans    int `json:"active_loans"`
        OverdueBooks   int `json:"overdue_books"`
    }

    err := database.DB.QueryRow(`
        SELECT 
            COUNT(CASE WHEN bl.is_returned = true THEN 1 END) as total_books_read,
            COUNT(CASE WHEN bl.is_returned = false THEN 1 END) as active_loans,
            COUNT(CASE WHEN bl.is_returned = false AND bl.due_date < CURRENT_DATE THEN 1 END) as overdue_books
        FROM users u
        LEFT JOIN book_loans bl ON u.id = bl.user_id
        WHERE u.id = $1
    `, userID).Scan(&stats.TotalBooksRead, &stats.ActiveLoans, &stats.OverdueBooks)

    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get ticket stats"})
        return
    }

    c.JSON(http.StatusOK, stats)
}