package handlers

import (
    "net/http"
    "library-management/internal/database"
    "library-management/internal/models"
    "github.com/gin-gonic/gin"
)

func GetProfile(c *gin.Context) {
    userID := c.GetString("userID")
    
    var user models.User
    query := `SELECT id, email, first_name, last_name, phone, address, avatar_url, 
              role, ticket_number, ticket_linked, is_blocked, block_reason 
              FROM users WHERE id = $1`
    
    err := database.DB.QueryRow(query, userID).Scan(
        &user.ID, &user.Email, &user.FirstName, &user.LastName,
        &user.Phone, &user.Address, &user.AvatarURL, &user.Role,
        &user.TicketNumber, &user.TicketLinked, &user.IsBlocked, &user.BlockReason,
    )
    
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }

    c.JSON(http.StatusOK, user)
}

func UpdateProfile(c *gin.Context) {
    userID := c.GetString("userID")
    
    var updates models.User
    if err := c.ShouldBindJSON(&updates); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    query := `UPDATE users SET 
              first_name = COALESCE($1, first_name),
              last_name = COALESCE($2, last_name),
              phone = COALESCE($3, phone),
              address = COALESCE($4, address),
              updated_at = CURRENT_TIMESTAMP
              WHERE id = $5 RETURNING id`
    
    err := database.DB.QueryRow(query, 
        updates.FirstName, updates.LastName, 
        updates.Phone, updates.Address, userID).Scan(&userID)
    
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Profile updated successfully"})
}

func UploadAvatar(c *gin.Context) {
    userID := c.GetString("userID")
    
    file, err := c.FormFile("avatar")
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "No file provided"})
        return
    }

    avatarURL := "/uploads/" + file.Filename
    
    query := `UPDATE users SET avatar_url = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`
    _, err = database.DB.Exec(query, avatarURL, userID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update avatar"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"avatar_url": avatarURL})
}

func LinkTicket(c *gin.Context) {
    userID := c.GetString("userID")
    
    var req models.LinkTicketRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    query := `UPDATE users SET ticket_number = $1, ticket_linked = true WHERE id = $2`
    _, err := database.DB.Exec(query, req.TicketNumber, userID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to link ticket"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Ticket linked successfully"})
}

func BlockTicket(c *gin.Context) {
    userID := c.GetString("userID")
    
    var req struct {
        Reason string `json:"reason"`
    }
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    query := `UPDATE users SET is_blocked = true, block_reason = $1 WHERE id = $2`
    _, err := database.DB.Exec(query, req.Reason, userID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to block ticket"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Ticket blocked successfully"})
}