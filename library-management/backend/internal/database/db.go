package database

import (
    "database/sql"
    "fmt"
    "library-management/internal/config"
    _ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB(cfg *config.Config) error {
    connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)
    
    var err error
    DB, err = sql.Open("postgres", connStr)
    if err != nil {
        return fmt.Errorf("failed to open database: %w", err)
    }

    err = DB.Ping()
    if err != nil {
        return fmt.Errorf("failed to ping database: %w", err)
    }

    // Сначала создаем расширение
    _, err = DB.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`)
    if err != nil {
        return fmt.Errorf("failed to create extension: %w", err)
    }

    err = createTables()
    if err != nil {
        return fmt.Errorf("failed to create tables: %w", err)
    }

    return nil
}

func createTables() error {
    // Создаем таблицы в правильном порядке
    queries := []string{
        // Таблица пользователей
        `CREATE TABLE IF NOT EXISTS users (
            id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
            email VARCHAR(255) UNIQUE NOT NULL,
            password_hash VARCHAR(255) NOT NULL,
            first_name VARCHAR(100),
            last_name VARCHAR(100),
            phone VARCHAR(20),
            address TEXT,
            avatar_url TEXT,
            role VARCHAR(50) DEFAULT 'user',
            ticket_number VARCHAR(50) UNIQUE,
            ticket_linked BOOLEAN DEFAULT false,
            is_blocked BOOLEAN DEFAULT false,
            block_reason TEXT,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        )`,
        
        // Таблица книг
        `CREATE TABLE IF NOT EXISTS books (
            id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
            title VARCHAR(255) NOT NULL,
            author VARCHAR(255),
            isbn VARCHAR(20),
            publisher VARCHAR(255),
            year INTEGER,
            total_copies INTEGER DEFAULT 1,
            available_copies INTEGER DEFAULT 1,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        )`,
        
        // Таблица выдачи книг (должна создаваться после users и books)
        `CREATE TABLE IF NOT EXISTS book_loans (
            id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
            user_id UUID REFERENCES users(id) ON DELETE CASCADE,
            book_id UUID REFERENCES books(id) ON DELETE CASCADE,
            loan_date DATE NOT NULL,
            due_date DATE NOT NULL,
            return_date DATE,
            is_returned BOOLEAN DEFAULT false,
            is_renewed BOOLEAN DEFAULT false,
            renewed_at TIMESTAMP,
            status VARCHAR(50) DEFAULT 'active',
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        )`,
        
        // Таблица уведомлений
        `CREATE TABLE IF NOT EXISTS notifications (
            id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
            user_id UUID REFERENCES users(id) ON DELETE CASCADE,
            title VARCHAR(255) NOT NULL,
            message TEXT NOT NULL,
            type VARCHAR(50),
            is_read BOOLEAN DEFAULT false,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        )`,
        
        // Таблица системных уведомлений
        `CREATE TABLE IF NOT EXISTS system_notifications (
            id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
            title VARCHAR(255) NOT NULL,
            message TEXT NOT NULL,
            priority VARCHAR(50) DEFAULT 'normal',
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        )`,
    }

    for i, query := range queries {
        _, err := DB.Exec(query)
        if err != nil {
            return fmt.Errorf("failed to execute query #%d: %s, error: %w", i+1, query, err)
        }
    }

    return nil
}