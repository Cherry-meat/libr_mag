package handlers

import (
	"net/http"
	"library-management/internal/database"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetUsers(c *gin.Context) {
	rows, err := database.DB.Query(`
		SELECT id, email, first_name, last_name, role, ticket_linked, is_blocked, created_at
		FROM users
		ORDER BY created_at DESC
	`)
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get users"})
		return
	}
	defer rows.Close()

	var users []map[string]interface{}
	for rows.Next() {
		var id uuid.UUID
		var email, firstName, lastName, role string
		var ticketLinked, isBlocked bool
		var createdAt string
		
		rows.Scan(&id, &email, &firstName, &lastName, &role, &ticketLinked, &isBlocked, &createdAt)
		
		users = append(users, map[string]interface{}{
			"id": id,
			"email": email,
			"first_name": firstName,
			"last_name": lastName,
			"role": role,
			"ticket_linked": ticketLinked,
			"is_blocked": isBlocked,
			"created_at": createdAt,
		})
	}

	c.JSON(http.StatusOK, users)
}

func UpdateUserRole(c *gin.Context) {
	userID := c.Param("id")
	
	var req struct {
		Role string `json:"role" binding:"required"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := `UPDATE users SET role = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`
	result, err := database.DB.Exec(query, req.Role, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update role"})
		return
	}
	
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User role updated successfully"})
}

func AdminBlockUser(c *gin.Context) {
	userID := c.Param("id")
	
	var req struct {
		Blocked bool   `json:"blocked" binding:"required"`
		Reason  string `json:"reason"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := `UPDATE users SET is_blocked = $1, block_reason = $2, updated_at = CURRENT_TIMESTAMP WHERE id = $3`
	_, err := database.DB.Exec(query, req.Blocked, req.Reason, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User status updated successfully"})
}

func GetDebtors(c *gin.Context) {
	rows, err := database.DB.Query(`
		SELECT u.id, u.email, u.first_name, u.last_name, 
		       COUNT(bl.id) as overdue_books
		FROM users u
		JOIN book_loans bl ON u.id = bl.user_id
		WHERE bl.due_date < CURRENT_DATE AND bl.is_returned = false
		GROUP BY u.id
		ORDER BY overdue_books DESC
	`)
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get debtors"})
		return
	}
	defer rows.Close()

	var debtors []map[string]interface{}
	for rows.Next() {
		var id uuid.UUID
		var email, firstName, lastName string
		var overdueBooks int
		
		rows.Scan(&id, &email, &firstName, &lastName, &overdueBooks)
		
		debtors = append(debtors, map[string]interface{}{
			"id": id,
			"email": email,
			"first_name": firstName,
			"last_name": lastName,
			"overdue_books": overdueBooks,
		})
	}

	c.JSON(http.StatusOK, debtors)
}