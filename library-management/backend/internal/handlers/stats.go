package handlers

import (
    "database/sql"
    "net/http"
    "library-management/internal/database"
    "github.com/gin-gonic/gin"
)

func GetUserStats(c *gin.Context) {
    userID := c.GetString("userID")
    
    var stats struct {
        ActiveLoans    int `json:"active_loans"`
        OverdueBooks   int `json:"overdue_books"`
        TotalBooksRead int `json:"total_books_read"`
    }

    err := database.DB.QueryRow(`
        SELECT 
            COUNT(CASE WHEN bl.is_returned = false THEN 1 END) as active_loans,
            COUNT(CASE WHEN bl.is_returned = false AND bl.due_date < CURRENT_DATE THEN 1 END) as overdue_books,
            COUNT(CASE WHEN bl.is_returned = true THEN 1 END) as total_books_read
        FROM users u
        LEFT JOIN book_loans bl ON u.id = bl.user_id
        WHERE u.id = $1
    `, userID).Scan(&stats.ActiveLoans, &stats.OverdueBooks, &stats.TotalBooksRead)
    
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get stats"})
        return
    }

    c.JSON(http.StatusOK, stats)
}

func GetCalendar(c *gin.Context) {
    userID := c.GetString("userID")
    
    rows, err := database.DB.Query(`
        SELECT loan_date, due_date, return_date, b.title, bl.is_returned
        FROM book_loans bl
        JOIN books b ON bl.book_id = b.id
        WHERE bl.user_id = $1
        ORDER BY loan_date DESC
    `, userID)
    
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get calendar data"})
        return
    }
    defer rows.Close()

    var events []map[string]interface{}
    for rows.Next() {
        var loanDate, dueDate string
        var returnDate sql.NullString
        var title string
        var isReturned bool
        
        err := rows.Scan(&loanDate, &dueDate, &returnDate, &title, &isReturned)
        if err != nil {
            continue
        }
        
        event := map[string]interface{}{
            "title": title,
            "loan_date": loanDate,
            "due_date": dueDate,
            "is_returned": isReturned,
        }
        
        if returnDate.Valid {
            event["return_date"] = returnDate.String
        }
        
        events = append(events, event)
    }

    c.JSON(http.StatusOK, events)
}