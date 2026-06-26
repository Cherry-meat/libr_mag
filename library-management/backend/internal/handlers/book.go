package handlers

import (
	"database/sql"
	"net/http"
	"time"
	"library-management/internal/database"
	"library-management/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetBookHistory(c *gin.Context) {
	userID := c.GetString("userID")
	
	rows, err := database.DB.Query(`
		SELECT bl.id, bl.book_id, b.title, b.author, bl.loan_date, bl.due_date, 
		       bl.return_date, bl.is_returned, bl.is_renewed, bl.status
		FROM book_loans bl
		JOIN books b ON bl.book_id = b.id
		WHERE bl.user_id = $1
		ORDER BY bl.loan_date DESC
	`, userID)
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get book history"})
		return
	}
	defer rows.Close()

	var loans []map[string]interface{}
	for rows.Next() {
		var loan models.BookLoan
		var title, author string
		var returnDate sql.NullTime
		
		rows.Scan(&loan.ID, &loan.BookID, &title, &author, &loan.LoanDate, 
			&loan.DueDate, &returnDate, &loan.IsReturned, &loan.IsRenewed, &loan.Status)
		
		if returnDate.Valid {
			loan.ReturnDate = &returnDate.Time
		}
		
		loans = append(loans, map[string]interface{}{
			"id": loan.ID,
			"book_id": loan.BookID,
			"title": title,
			"author": author,
			"loan_date": loan.LoanDate,
			"due_date": loan.DueDate,
			"return_date": loan.ReturnDate,
			"is_returned": loan.IsReturned,
			"is_renewed": loan.IsRenewed,
			"status": loan.Status,
		})
	}

	c.JSON(http.StatusOK, loans)
}

func RenewBook(c *gin.Context) {
	userID := c.GetString("userID")
	
	var req models.RenewBookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	loanID, err := uuid.Parse(req.LoanID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid loan ID"})
		return
	}

	var dueDate time.Time
	query := `SELECT due_date FROM book_loans WHERE id = $1 AND user_id = $2 AND is_returned = false`
	err = database.DB.QueryRow(query, loanID, userID).Scan(&dueDate)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Loan not found or already returned"})
		return
	}

	newDueDate := dueDate.AddDate(0, 0, 14)
	
	updateQuery := `UPDATE book_loans 
	                SET due_date = $1, is_renewed = true, renewed_at = CURRENT_TIMESTAMP 
	                WHERE id = $2`
	_, err = database.DB.Exec(updateQuery, newDueDate, loanID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to renew book"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Book renewed successfully",
		"new_due_date": newDueDate,
	})
}

func ConfirmBookIssue(c *gin.Context) {
	var req struct {
		UserID string `json:"user_id" binding:"required"`
		BookID string `json:"book_id" binding:"required"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := uuid.Parse(req.UserID)
	bookID, _ := uuid.Parse(req.BookID)

	var available int
	checkQuery := `SELECT available_copies FROM books WHERE id = $1`
	err := database.DB.QueryRow(checkQuery, bookID).Scan(&available)
	if err != nil || available <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Book not available"})
		return
	}

	tx, _ := database.DB.Begin()
	defer tx.Rollback()

	loanID := uuid.New()
	loanDate := time.Now()
	dueDate := loanDate.AddDate(0, 0, 21)
	
	insertQuery := `INSERT INTO book_loans (id, user_id, book_id, loan_date, due_date, status) 
	                VALUES ($1, $2, $3, $4, $5, $6)`
	tx.Exec(insertQuery, loanID, userID, bookID, loanDate, dueDate, "active")

	updateQuery := `UPDATE books SET available_copies = available_copies - 1 WHERE id = $1`
	tx.Exec(updateQuery, bookID)

	tx.Commit()
	c.JSON(http.StatusOK, gin.H{"message": "Book issued successfully"})
}

func ConfirmBookReturn(c *gin.Context) {
	var req struct {
		LoanID string `json:"loan_id" binding:"required"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	loanID, _ := uuid.Parse(req.LoanID)

	tx, _ := database.DB.Begin()
	defer tx.Rollback()

	var bookID uuid.UUID
	getQuery := `SELECT book_id FROM book_loans WHERE id = $1 AND is_returned = false`
	tx.QueryRow(getQuery, loanID).Scan(&bookID)

	now := time.Now()
	updateQuery := `UPDATE book_loans SET return_date = $1, is_returned = true, status = 'returned' WHERE id = $2`
	tx.Exec(updateQuery, now, loanID)

	updateBookQuery := `UPDATE books SET available_copies = available_copies + 1 WHERE id = $1`
	tx.Exec(updateBookQuery, bookID)

	tx.Commit()
	c.JSON(http.StatusOK, gin.H{"message": "Book returned successfully"})
}