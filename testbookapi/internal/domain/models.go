package domain

import "time"

type User struct {
	ID           int64
	Email        string
	PasswordHash string
	IsAdmin      bool
	CreatedAt    time.Time
}

type Category struct {
	ID   int64
	Name string
}

type Book struct {
	ID        int64
	Title     string
	Year      int
	Author    string
	PriceUSD  float64
	CategoryID int64
	Stock     int
	CreatedAt time.Time
}

type CartItem struct {
	ID        int64
	UserID    int64
	BookID    int64
	AddedAt   time.Time
	ExpiresAt time.Time
}

type BookFilter struct {
	CategoryIDs []int64
	Limit       int
	Offset      int
}
