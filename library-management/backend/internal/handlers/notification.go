package handlers

import (
	"net/http"
	"library-management/internal/database"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetNotifications(c *gin.Context) {
	userID := c.GetString("userID")
	
	rows, err := database.DB.Query(`
		SELECT id, title, message, type, is_read, created_at 
		FROM notifications 
		WHERE user_id = $1 
		ORDER BY created_at DESC
	`, userID)
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get notifications"})
		return
	}
	defer rows.Close()

	var notifications []map[string]interface{}
	for rows.Next() {
		var id uuid.UUID
		var title, message, notificationType string
		var isRead bool
		var createdAt string
		
		rows.Scan(&id, &title, &message, &notificationType, &isRead, &createdAt)
		
		notifications = append(notifications, map[string]interface{}{
			"id": id,
			"title": title,
			"message": message,
			"type": notificationType,
			"is_read": isRead,
			"created_at": createdAt,
		})
	}

	c.JSON(http.StatusOK, notifications)
}

func MarkNotificationRead(c *gin.Context) {
	notificationID := c.Param("id")
	userID := c.GetString("userID")
	
	query := `UPDATE notifications SET is_read = true WHERE id = $1 AND user_id = $2`
	result, err := database.DB.Exec(query, notificationID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update notification"})
		return
	}
	
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Notification not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Notification marked as read"})
}

func CreateSystemNotification(c *gin.Context) {
	var req struct {
		Title    string `json:"title" binding:"required"`
		Message  string `json:"message" binding:"required"`
		Priority string `json:"priority"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Priority == "" {
		req.Priority = "normal"
	}

	query := `INSERT INTO system_notifications (id, title, message, priority) VALUES ($1, $2, $3, $4)`
	id := uuid.New()
	_, err := database.DB.Exec(query, id, req.Title, req.Message, req.Priority)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create notification"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "System notification created"})
}

func GetSystemNotifications(c *gin.Context) {
	rows, err := database.DB.Query(`
		SELECT id, title, message, priority, created_at 
		FROM system_notifications 
		ORDER BY created_at DESC
	`)
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get system notifications"})
		return
	}
	defer rows.Close()

	var notifications []map[string]interface{}
	for rows.Next() {
		var id uuid.UUID
		var title, message, priority string
		var createdAt string
		
		rows.Scan(&id, &title, &message, &priority, &createdAt)
		
		notifications = append(notifications, map[string]interface{}{
			"id": id,
			"title": title,
			"message": message,
			"priority": priority,
			"created_at": createdAt,
		})
	}

	c.JSON(http.StatusOK, notifications)
}