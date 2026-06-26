package main

import (
	"log"
	"library-management/internal/config"
	"library-management/internal/database"
	"library-management/internal/handlers"
	"library-management/internal/middleware"
	"library-management/internal/utils"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	utils.InitJWT(cfg.JWTSecret)

	err = database.InitDB(cfg)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Публичные маршруты
	auth := r.Group("/api/auth")
	{
		auth.POST("/register", handlers.Register)
		auth.POST("/login", handlers.Login)
	}

	// Защищенные маршруты
	protected := r.Group("/api")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.GET("/profile", handlers.GetProfile)
		protected.PUT("/profile", handlers.UpdateProfile)
		protected.POST("/profile/avatar", handlers.UploadAvatar)
		protected.POST("/profile/link-ticket", handlers.LinkTicket)
		protected.POST("/profile/block-ticket", handlers.BlockTicket)

		protected.GET("/books/history", handlers.GetBookHistory)
		protected.POST("/books/renew", handlers.RenewBook)

		protected.GET("/notifications", handlers.GetNotifications)
		protected.PUT("/notifications/:id", handlers.MarkNotificationRead)

		protected.GET("/stats", handlers.GetUserStats)
		protected.GET("/calendar", handlers.GetCalendar)
		protected.GET("/ticket", handlers.GetTicket)
	}

	// Админские маршруты
	admin := r.Group("/api/admin")
	admin.Use(middleware.AuthMiddleware(), middleware.AdminMiddleware())
	{
		admin.GET("/users", handlers.GetUsers)
		admin.PUT("/users/:id/role", handlers.UpdateUserRole)
		admin.POST("/users/:id/block", handlers.AdminBlockUser)
		admin.GET("/debtors", handlers.GetDebtors)
		admin.GET("/reports", handlers.GenerateReport)
		admin.POST("/books/issue", handlers.ConfirmBookIssue)
		admin.POST("/books/return", handlers.ConfirmBookReturn)
		admin.GET("/system-notifications", handlers.GetSystemNotifications)
		admin.POST("/system-notifications", handlers.CreateSystemNotification)
		
		// Добавляем статистику для админа
		admin.GET("/stats", func(c *gin.Context) {
			var stats struct {
				TotalUsers   int `json:"total_users"`
				TotalBooks   int `json:"total_books"`
				ActiveLoans  int `json:"active_loans"`
				OverdueLoans int `json:"overdue_loans"`
			}
			database.DB.QueryRow(`
				SELECT 
					(SELECT COUNT(*) FROM users) as total_users,
					(SELECT COUNT(*) FROM books) as total_books,
					(SELECT COUNT(*) FROM book_loans WHERE is_returned = false) as active_loans,
					(SELECT COUNT(*) FROM book_loans WHERE is_returned = false AND due_date < CURRENT_DATE) as overdue_loans
			`).Scan(&stats.TotalUsers, &stats.TotalBooks, &stats.ActiveLoans, &stats.OverdueLoans)
			c.JSON(200, stats)
		})
	}

	log.Println("Server starting on :8080")
	r.Run(":8080")
}