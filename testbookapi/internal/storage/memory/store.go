package memory

import (
	"context"
	"errors"
	"strings"
	"sync"
	"time"

	"lessonsgo/testbookapi/internal/domain"
)

type Store struct {
	mu         sync.Mutex
	nextUserID int64
	nextCatID  int64
	nextBookID int64
	nextCartID int64

	users      map[int64]domain.User
	usersByEmail map[string]int64
	categories map[int64]domain.Category
	books      map[int64]domain.Book
	cartItems  map[int64]domain.CartItem
	cartByUser map[int64]map[int64]int64 // userID -> bookID -> cartItemID
}

func New() *Store {
	return &Store{
		nextUserID:   1,
		nextCatID:    1,
		nextBookID:   1,
		nextCartID:   1,
		users:        make(map[int64]domain.User),
		usersByEmail: make(map[string]int64),
		categories:   make(map[int64]domain.Category),
		books:        make(map[int64]domain.Book),
		cartItems:    make(map[int64]domain.CartItem),
		cartByUser:   make(map[int64]map[int64]int64),
	}
}

func (s *Store) CreateUser(ctx context.Context, email, passwordHash string) (domain.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	email = strings.ToLower(strings.TrimSpace(email))
	if _, exists := s.usersByEmail[email]; exists {
		return domain.User{}, domain.ErrEmailTaken
	}

	user := domain.User{
		ID:           s.nextUserID,
		Email:        email,
		PasswordHash: passwordHash,
		CreatedAt:    time.Now().UTC(),
	}
	s.nextUserID++
	s.users[user.ID] = user
	s.usersByEmail[email] = user.ID
	return user, nil
}

func (s *Store) GetUserByEmail(ctx context.Context, email string) (domain.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	id, ok := s.usersByEmail[strings.ToLower(strings.TrimSpace(email))]
	if !ok {
		return domain.User{}, domain.ErrNotFound
	}
	return s.users[id], nil
}

func (s *Store) GetUserByID(ctx context.Context, id int64) (domain.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	user, ok := s.users[id]
	if !ok {
		return domain.User{}, domain.ErrNotFound
	}
	return user, nil
}

func (s *Store) CreateCategory(ctx context.Context, name string) (domain.Category, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, c := range s.categories {
		if strings.EqualFold(c.Name, name) {
			return domain.Category{}, domain.ErrInvalidInput
		}
	}

	cat := domain.Category{ID: s.nextCatID, Name: name}
	s.nextCatID++
	s.categories[cat.ID] = cat
	return cat, nil
}

func (s *Store) GetCategory(ctx context.Context, id int64) (domain.Category, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	cat, ok := s.categories[id]
	if !ok {
		return domain.Category{}, domain.ErrNotFound
	}
	return cat, nil
}

func (s *Store) ListCategories(ctx context.Context) ([]domain.Category, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	out := make([]domain.Category, 0, len(s.categories))
	for _, c := range s.categories {
		out = append(out, c)
	}
	return out, nil
}

func (s *Store) UpdateCategory(ctx context.Context, id int64, name string) (domain.Category, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	cat, ok := s.categories[id]
	if !ok {
		return domain.Category{}, domain.ErrNotFound
	}
	cat.Name = name
	s.categories[id] = cat
	return cat, nil
}

func (s *Store) DeleteCategory(ctx context.Context, id int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.categories[id]; !ok {
		return domain.ErrNotFound
	}
	delete(s.categories, id)
	return nil
}

func (s *Store) CategoryBookCount(ctx context.Context, categoryID int64) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	count := 0
	for _, b := range s.books {
		if b.CategoryID == categoryID {
			count++
		}
	}
	return count, nil
}

func (s *Store) CreateBook(ctx context.Context, book domain.Book) (domain.Book, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.categories[book.CategoryID]; !ok {
		return domain.Book{}, domain.ErrNotFound
	}

	book.ID = s.nextBookID
	book.CreatedAt = time.Now().UTC()
	s.nextBookID++
	s.books[book.ID] = book
	return book, nil
}

func (s *Store) GetBook(ctx context.Context, id int64) (domain.Book, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	book, ok := s.books[id]
	if !ok {
		return domain.Book{}, domain.ErrNotFound
	}
	return book, nil
}

func (s *Store) ListBooks(ctx context.Context, filter domain.BookFilter) ([]domain.Book, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.releaseExpiredLocked(time.Now().UTC())

	books := s.filteredBooksLocked(filter)
	if filter.Offset >= len(books) {
		return []domain.Book{}, nil
	}
	end := filter.Offset + filter.Limit
	if end > len(books) {
		end = len(books)
	}
	return books[filter.Offset:end], nil
}

func (s *Store) CountBooks(ctx context.Context, filter domain.BookFilter) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.releaseExpiredLocked(time.Now().UTC())
	return len(s.filteredBooksLocked(filter)), nil
}

func (s *Store) filteredBooksLocked(filter domain.BookFilter) []domain.Book {
	categorySet := make(map[int64]struct{}, len(filter.CategoryIDs))
	for _, id := range filter.CategoryIDs {
		categorySet[id] = struct{}{}
	}

	out := make([]domain.Book, 0)
	for _, b := range s.books {
		if b.Stock <= 0 {
			continue
		}
		if len(categorySet) > 0 {
			if _, ok := categorySet[b.CategoryID]; !ok {
				continue
			}
		}
		out = append(out, b)
	}
	return out
}

func (s *Store) UpdateBook(ctx context.Context, id int64, title string, year int, author string, priceUSD float64, categoryID int64) (domain.Book, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	book, ok := s.books[id]
	if !ok {
		return domain.Book{}, domain.ErrNotFound
	}
	if _, ok := s.categories[categoryID]; !ok {
		return domain.Book{}, domain.ErrNotFound
	}

	book.Title = title
	book.Year = year
	book.Author = author
	book.PriceUSD = priceUSD
	book.CategoryID = categoryID
	s.books[id] = book
	return book, nil
}

func (s *Store) DeleteBook(ctx context.Context, id int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.books[id]; !ok {
		return domain.ErrNotFound
	}
	delete(s.books, id)
	return nil
}

func (s *Store) ListCartItems(ctx context.Context, userID int64) ([]domain.CartItem, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.releaseExpiredLocked(time.Now().UTC())

	bookMap := s.cartByUser[userID]
	out := make([]domain.CartItem, 0, len(bookMap))
	for _, itemID := range bookMap {
		out = append(out, s.cartItems[itemID])
	}
	return out, nil
}

func (s *Store) AddCartItem(ctx context.Context, userID, bookID int64, expiresAt time.Time) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.releaseExpiredLocked(time.Now().UTC())

	book, ok := s.books[bookID]
	if !ok {
		return domain.ErrNotFound
	}
	if book.Stock <= 0 {
		return domain.ErrOutOfStock
	}

	if s.cartByUser[userID] == nil {
		s.cartByUser[userID] = make(map[int64]int64)
	}
	if _, exists := s.cartByUser[userID][bookID]; exists {
		return domain.ErrAlreadyInCart
	}

	book.Stock--
	s.books[bookID] = book

	now := time.Now().UTC()
	item := domain.CartItem{
		ID:        s.nextCartID,
		UserID:    userID,
		BookID:    bookID,
		AddedAt:   now,
		ExpiresAt: expiresAt,
	}
	s.nextCartID++
	s.cartItems[item.ID] = item
	s.cartByUser[userID][bookID] = item.ID
	return nil
}

func (s *Store) RemoveCartItem(ctx context.Context, userID, bookID int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.removeCartItemLocked(userID, bookID)
}

func (s *Store) ClearCart(ctx context.Context, userID int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	bookMap := s.cartByUser[userID]
	for bookID := range bookMap {
		if err := s.removeCartItemLocked(userID, bookID); err != nil && !errors.Is(err, domain.ErrNotFound) {
			return err
		}
	}
	return nil
}

func (s *Store) ReleaseExpiredCartItems(ctx context.Context) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.releaseExpiredLocked(time.Now().UTC()), nil
}

func (s *Store) CheckoutCart(ctx context.Context, userID int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.releaseExpiredLocked(time.Now().UTC())

	bookMap := s.cartByUser[userID]
	if len(bookMap) == 0 {
		return domain.ErrCartEmpty
	}

	for bookID := range bookMap {
		itemID := bookMap[bookID]
		delete(s.cartItems, itemID)
	}
	delete(s.cartByUser, userID)
	return nil
}

func (s *Store) removeCartItemLocked(userID, bookID int64) error {
	bookMap := s.cartByUser[userID]
	if bookMap == nil {
		return domain.ErrNotFound
	}
	itemID, ok := bookMap[bookID]
	if !ok {
		return domain.ErrNotFound
	}

	delete(s.cartItems, itemID)
	delete(bookMap, bookID)
	if len(bookMap) == 0 {
		delete(s.cartByUser, userID)
	}

	if book, ok := s.books[bookID]; ok {
		book.Stock++
		s.books[bookID] = book
	}
	return nil
}

func (s *Store) releaseExpiredLocked(now time.Time) int {
	released := 0
	for id, item := range s.cartItems {
		if item.ExpiresAt.After(now) {
			continue
		}
		if err := s.removeCartItemLocked(item.UserID, item.BookID); err == nil {
			released++
		} else {
			delete(s.cartItems, id)
		}
	}
	return released
}

// SetCartItemExpiry adjusts expiry for tests.
func (s *Store) SetCartItemExpiry(itemID int64, expiresAt time.Time) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	item, ok := s.cartItems[itemID]
	if !ok {
		return domain.ErrNotFound
	}
	item.ExpiresAt = expiresAt
	s.cartItems[itemID] = item
	return nil
}

func (s *Store) SetAdmin(userID int64, isAdmin bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	user, ok := s.users[userID]
	if !ok {
		return domain.ErrNotFound
	}
	user.IsAdmin = isAdmin
	s.users[userID] = user
	return nil
}
