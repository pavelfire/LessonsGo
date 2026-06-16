package domain

import (
	"context"
	"time"
)

type Repository interface {
	// Users
	CreateUser(ctx context.Context, email, passwordHash string) (User, error)
	GetUserByEmail(ctx context.Context, email string) (User, error)
	GetUserByID(ctx context.Context, id int64) (User, error)

	// Categories
	CreateCategory(ctx context.Context, name string) (Category, error)
	GetCategory(ctx context.Context, id int64) (Category, error)
	ListCategories(ctx context.Context) ([]Category, error)
	UpdateCategory(ctx context.Context, id int64, name string) (Category, error)
	DeleteCategory(ctx context.Context, id int64) error
	CategoryBookCount(ctx context.Context, categoryID int64) (int, error)

	// Books
	CreateBook(ctx context.Context, book Book) (Book, error)
	GetBook(ctx context.Context, id int64) (Book, error)
	ListBooks(ctx context.Context, filter BookFilter) ([]Book, error)
	CountBooks(ctx context.Context, filter BookFilter) (int, error)
	UpdateBook(ctx context.Context, id int64, title string, year int, author string, priceUSD float64, categoryID int64) (Book, error)
	DeleteBook(ctx context.Context, id int64) error

	// Cart
	ListCartItems(ctx context.Context, userID int64) ([]CartItem, error)
	AddCartItem(ctx context.Context, userID, bookID int64, expiresAt time.Time) error
	RemoveCartItem(ctx context.Context, userID, bookID int64) error
	ClearCart(ctx context.Context, userID int64) error
	ReleaseExpiredCartItems(ctx context.Context) (int, error)
	CheckoutCart(ctx context.Context, userID int64) error
}
