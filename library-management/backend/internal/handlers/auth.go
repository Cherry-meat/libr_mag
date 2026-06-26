package handlers

import (
    "database/sql"
    "net/http"
    "library-management/internal/database"
    "library-management/internal/models"
    "library-management/internal/utils"
    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
    "golang.org/x/crypto/bcrypt"
)

func Register(c *gin.Context) {
    var req models.RegisterRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
        return
    }

    query := `INSERT INTO users (id, email, password_hash, first_name, last_name, phone, role) 
              VALUES ($1, $2, $3, $4, $5, $6, $7)`
    
    userID := uuid.New()
    _, err = database.DB.Exec(query, userID, req.Email, string(hashedPassword), 
        req.FirstName, req.LastName, req.Phone, "user")
    
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
        return
    }

    c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully", "user_id": userID})
}

func Login(c *gin.Context) {
    var req models.LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    var user models.User
    query := `SELECT id, email, password_hash, role FROM users WHERE email = $1`
    err := database.DB.QueryRow(query, req.Email).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.Role)
    
    if err == sql.ErrNoRows {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
        return
    } else if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
        return
    }

    err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
        return
    }

    token, err := utils.GenerateToken(user.ID, user.Email, user.Role)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"token": token, "user": user})
}