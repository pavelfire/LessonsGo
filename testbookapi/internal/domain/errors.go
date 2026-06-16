package domain

import "errors"

var (
	ErrNotFound            = errors.New("not found")
	ErrUnauthorized        = errors.New("unauthorized")
	ErrForbidden           = errors.New("forbidden")
	ErrInvalidInput        = errors.New("invalid input")
	ErrEmailTaken          = errors.New("email already registered")
	ErrInvalidCredentials  = errors.New("invalid email or password")
	ErrOutOfStock          = errors.New("book is out of stock")
	ErrAlreadyInCart       = errors.New("book is already in the cart")
	ErrCategoryInUse       = errors.New("category is in use by books")
	ErrStockNotEditable    = errors.New("stock cannot be edited after creation")
	ErrCartEmpty           = errors.New("cart is empty")
	ErrBookNotPurchasable  = errors.New("book is not available for purchase")
)
