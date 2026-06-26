package handlers

import (
	"encoding/csv"
	"net/http"
	"strconv"
	"library-management/internal/database"
	"github.com/gin-gonic/gin"
)

func GenerateReport(c *gin.Context) {
	reportType := c.Query("type")
	if reportType == "" {
		reportType = "loans"
	}
	
	switch reportType {
	case "loans":
		generateLoansReport(c)
	case "users":
		generateUsersReport(c)
	case "books":
		generateBooksReport(c)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid report type"})
	}
}

func generateLoansReport(c *gin.Context) {
	rows, err := database.DB.Query(`
		SELECT 
			u.first_name || ' ' || u.last_name as user_name,
			u.email,
			b.title as book_title,
			b.author,
			bl.loan_date,
			bl.due_date,
			CASE 
				WHEN bl.is_returned THEN 'Returned'
				WHEN bl.due_date < CURRENT_DATE THEN 'Overdue'
				ELSE 'Active'
			END as status
		FROM book_loans bl
		JOIN users u ON bl.user_id = u.id
		JOIN books b ON bl.book_id = b.id
		ORDER BY bl.loan_date DESC
	`)
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate report"})
		return
	}
	defer rows.Close()

	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", "attachment;filename=loans_report.csv")
	
	writer := csv.NewWriter(c.Writer)
	defer writer.Flush()
	
	writer.Write([]string{"User", "Email", "Book", "Author", "Loan Date", "Due Date", "Status"})
	
	for rows.Next() {
		var userName, email, bookTitle, author, loanDate, dueDate, status string
		rows.Scan(&userName, &email, &bookTitle, &author, &loanDate, &dueDate, &status)
		writer.Write([]string{userName, email, bookTitle, author, loanDate, dueDate, status})
	}
}

func generateUsersReport(c *gin.Context) {
	rows, err := database.DB.Query(`
		SELECT 
			u.email,
			u.first_name || ' ' || u.last_name as full_name,
			u.role,
			u.ticket_linked,
			u.is_blocked,
			COUNT(bl.id) as total_loans,
			COUNT(CASE WHEN bl.is_returned = true THEN 1 END) as returned_loans,
			COUNT(CASE WHEN bl.is_returned = false AND bl.due_date < CURRENT_DATE THEN 1 END) as overdue_loans
		FROM users u
		LEFT JOIN book_loans bl ON u.id = bl.user_id
		GROUP BY u.id
		ORDER BY total_loans DESC
	`)
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate report"})
		return
	}
	defer rows.Close()

	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", "attachment;filename=users_report.csv")
	
	writer := csv.NewWriter(c.Writer)
	defer writer.Flush()
	
	writer.Write([]string{"Email", "Full Name", "Role", "Ticket Linked", "Blocked", "Total Loans", "Returned", "Overdue"})
	
	for rows.Next() {
		var email, fullName, role string
		var ticketLinked, isBlocked bool
		var totalLoans, returnedLoans, overdueLoans int
		
		rows.Scan(&email, &fullName, &role, &ticketLinked, &isBlocked, &totalLoans, &returnedLoans, &overdueLoans)
		
		ticketStr := "No"
		if ticketLinked {
			ticketStr = "Yes"
		}
		blockedStr := "No"
		if isBlocked {
			blockedStr = "Yes"
		}
		
		writer.Write([]string{
			email, fullName, role, ticketStr, blockedStr,
			strconv.Itoa(totalLoans), strconv.Itoa(returnedLoans), strconv.Itoa(overdueLoans),
		})
	}
}

func generateBooksReport(c *gin.Context) {
	rows, err := database.DB.Query(`
		SELECT 
			b.title,
			b.author,
			b.isbn,
			b.total_copies,
			b.available_copies,
			COUNT(bl.id) as total_loans,
			COUNT(CASE WHEN bl.is_returned = true THEN 1 END) as times_returned
		FROM books b
		LEFT JOIN book_loans bl ON b.id = bl.book_id
		GROUP BY b.id
		ORDER BY total_loans DESC
	`)
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate report"})
		return
	}
	defer rows.Close()

	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", "attachment;filename=books_report.csv")
	
	writer := csv.NewWriter(c.Writer)
	defer writer.Flush()
	
	writer.Write([]string{"Title", "Author", "ISBN", "Total Copies", "Available", "Total Loans", "Times Returned"})
	
	for rows.Next() {
		var title, author, isbn string
		var totalCopies, availableCopies, totalLoans, timesReturned int
		
		rows.Scan(&title, &author, &isbn, &totalCopies, &availableCopies, &totalLoans, &timesReturned)
		
		writer.Write([]string{
			title, author, isbn,
			strconv.Itoa(totalCopies), strconv.Itoa(availableCopies),
			strconv.Itoa(totalLoans), strconv.Itoa(timesReturned),
		})
	}
}