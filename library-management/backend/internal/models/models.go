package models

import (
    "time"
    "github.com/google/uuid"
)

type User struct {
    ID           uuid.UUID `json:"id"`
    Email        string    `json:"email"`
    PasswordHash string    `json:"-"`
    FirstName    string    `json:"first_name"`
    LastName     string    `json:"last_name"`
    Phone        string    `json:"phone"`
    Address      string    `json:"address"`
    AvatarURL    string    `json:"avatar_url"`
    Role         string    `json:"role"`
    TicketNumber string    `json:"ticket_number"`
    TicketLinked bool      `json:"ticket_linked"`
    IsBlocked    bool      `json:"is_blocked"`
    BlockReason  string    `json:"block_reason"`
    CreatedAt    time.Time `json:"created_at"`
    UpdatedAt    time.Time `json:"updated_at"`
}

type Book struct {
    ID             uuid.UUID `json:"id"`
    Title          string    `json:"title"`
    Author         string    `json:"author"`
    ISBN           string    `json:"isbn"`
    Publisher      string    `json:"publisher"`
    Year           int       `json:"year"`
    TotalCopies    int       `json:"total_copies"`
    AvailableCopies int      `json:"available_copies"`
    CreatedAt      time.Time `json:"created_at"`
}

type BookLoan struct {
    ID         uuid.UUID  `json:"id"`
    UserID     uuid.UUID  `json:"user_id"`
    BookID     uuid.UUID  `json:"book_id"`
    LoanDate   time.Time  `json:"loan_date"`
    DueDate    time.Time  `json:"due_date"`
    ReturnDate *time.Time `json:"return_date,omitempty"`
    IsReturned bool       `json:"is_returned"`
    IsRenewed  bool       `json:"is_renewed"`
    RenewedAt  *time.Time `json:"renewed_at,omitempty"`
    Status     string     `json:"status"`
    CreatedAt  time.Time  `json:"created_at"`
}

type Notification struct {
    ID        uuid.UUID `json:"id"`
    UserID    uuid.UUID `json:"user_id"`
    Title     string    `json:"title"`
    Message   string    `json:"message"`
    Type      string    `json:"type"`
    IsRead    bool      `json:"is_read"`
    CreatedAt time.Time `json:"created_at"`
}

type SystemNotification struct {
    ID        uuid.UUID `json:"id"`
    Title     string    `json:"title"`
    Message   string    `json:"message"`
    Priority  string    `json:"priority"`
    CreatedAt time.Time `json:"created_at"`
}

type LoginRequest struct {
    Email    string `json:"email" binding:"required"`
    Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
    Email     string `json:"email" binding:"required,email"`
    Password  string `json:"password" binding:"required,min=6"`
    FirstName string `json:"first_name" binding:"required"`
    LastName  string `json:"last_name" binding:"required"`
    Phone     string `json:"phone"`
}

type LinkTicketRequest struct {
    TicketNumber string `json:"ticket_number" binding:"required"`
}

type RenewBookRequest struct {
    LoanID string `json:"loan_id" binding:"required"`
}