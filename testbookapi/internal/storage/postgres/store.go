package postgres

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"strings"
	"time"

	"lessonsgo/testbookapi/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

type Store struct {
	pool *pgxpool.Pool
}

func New(ctx context.Context, databaseURL string) (*Store, error) {
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return nil, fmt.Errorf("connect to database: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}

	store := &Store{pool: pool}
	if err := store.migrate(ctx); err != nil {
		pool.Close()
		return nil, err
	}

	return store, nil
}

func (s *Store) Close() {
	s.pool.Close()
}

func (s *Store) migrate(ctx context.Context) error {
	sqlBytes, err := migrationsFS.ReadFile("migrations/001_init.sql")
	if err != nil {
		return fmt.Errorf("read migration: %w", err)
	}

	_, err = s.pool.Exec(ctx, string(sqlBytes))
	return err
}

func (s *Store) CreateUser(ctx context.Context, email, passwordHash string) (domain.User, error) {
	row := s.pool.QueryRow(ctx, `
		INSERT INTO users (email, password_hash)
		VALUES ($1, $2)
		RETURNING id, email, password_hash, is_admin, created_at
	`, email, passwordHash)

	user, err := scanUser(row)
	if err != nil {
		if isUniqueViolation(err) {
			return domain.User{}, domain.ErrEmailTaken
		}
		return domain.User{}, err
	}
	return user, nil
}

func (s *Store) GetUserByEmail(ctx context.Context, email string) (domain.User, error) {
	row := s.pool.QueryRow(ctx, `
		SELECT id, email, password_hash, is_admin, created_at
		FROM users WHERE email = $1
	`, email)
	return scanUser(row)
}

func (s *Store) GetUserByID(ctx context.Context, id int64) (domain.User, error) {
	row := s.pool.QueryRow(ctx, `
		SELECT id, email, password_hash, is_admin, created_at
		FROM users WHERE id = $1
	`, id)
	return scanUser(row)
}

func (s *Store) CreateCategory(ctx context.Context, name string) (domain.Category, error) {
	row := s.pool.QueryRow(ctx, `
		INSERT INTO categories (name) VALUES ($1)
		RETURNING id, name
	`, name)

	var cat domain.Category
	if err := row.Scan(&cat.ID, &cat.Name); err != nil {
		if isUniqueViolation(err) {
			return domain.Category{}, domain.ErrInvalidInput
		}
		return domain.Category{}, err
	}
	return cat, nil
}

func (s *Store) GetCategory(ctx context.Context, id int64) (domain.Category, error) {
	row := s.pool.QueryRow(ctx, `SELECT id, name FROM categories WHERE id = $1`, id)
	var cat domain.Category
	if err := row.Scan(&cat.ID, &cat.Name); err != nil {
		return domain.Category{}, mapNotFound(err)
	}
	return cat, nil
}

func (s *Store) ListCategories(ctx context.Context) ([]domain.Category, error) {
	rows, err := s.pool.Query(ctx, `SELECT id, name FROM categories ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []domain.Category
	for rows.Next() {
		var cat domain.Category
		if err := rows.Scan(&cat.ID, &cat.Name); err != nil {
			return nil, err
		}
		out = append(out, cat)
	}
	return out, rows.Err()
}

func (s *Store) UpdateCategory(ctx context.Context, id int64, name string) (domain.Category, error) {
	row := s.pool.QueryRow(ctx, `
		UPDATE categories SET name = $2 WHERE id = $1
		RETURNING id, name
	`, id, name)

	var cat domain.Category
	if err := row.Scan(&cat.ID, &cat.Name); err != nil {
		return domain.Category{}, mapNotFound(err)
	}
	return cat, nil
}

func (s *Store) DeleteCategory(ctx context.Context, id int64) error {
	tag, err := s.pool.Exec(ctx, `DELETE FROM categories WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (s *Store) CategoryBookCount(ctx context.Context, categoryID int64) (int, error) {
	row := s.pool.QueryRow(ctx, `SELECT COUNT(*) FROM books WHERE category_id = $1`, categoryID)
	var count int
	if err := row.Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

func (s *Store) CreateBook(ctx context.Context, book domain.Book) (domain.Book, error) {
	row := s.pool.QueryRow(ctx, `
		INSERT INTO books (title, year_published, author, price_usd, category_id, stock)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, title, year_published, author, price_usd, category_id, stock, created_at
	`, book.Title, book.Year, book.Author, book.PriceUSD, book.CategoryID, book.Stock)

	return scanBook(row)
}

func (s *Store) GetBook(ctx context.Context, id int64) (domain.Book, error) {
	row := s.pool.QueryRow(ctx, `
		SELECT id, title, year_published, author, price_usd, category_id, stock, created_at
		FROM books WHERE id = $1
	`, id)
	return scanBook(row)
}

func (s *Store) ListBooks(ctx context.Context, filter domain.BookFilter) ([]domain.Book, error) {
	query, args := buildBooksQuery(filter, false)
	rows, err := s.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []domain.Book
	for rows.Next() {
		book, err := scanBookRow(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, book)
	}
	return out, rows.Err()
}

func (s *Store) CountBooks(ctx context.Context, filter domain.BookFilter) (int, error) {
	query, args := buildBooksQuery(filter, true)
	row := s.pool.QueryRow(ctx, query, args...)
	var count int
	if err := row.Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

func (s *Store) UpdateBook(ctx context.Context, id int64, title string, year int, author string, priceUSD float64, categoryID int64) (domain.Book, error) {
	row := s.pool.QueryRow(ctx, `
		UPDATE books
		SET title = $2, year_published = $3, author = $4, price_usd = $5, category_id = $6
		WHERE id = $1
		RETURNING id, title, year_published, author, price_usd, category_id, stock, created_at
	`, id, title, year, author, priceUSD, categoryID)
	book, err := scanBook(row)
	if err != nil {
		return domain.Book{}, mapNotFound(err)
	}
	return book, nil
}

func (s *Store) DeleteBook(ctx context.Context, id int64) error {
	tag, err := s.pool.Exec(ctx, `DELETE FROM books WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (s *Store) ListCartItems(ctx context.Context, userID int64) ([]domain.CartItem, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT id, user_id, book_id, added_at, expires_at
		FROM cart_items
		WHERE user_id = $1
		ORDER BY added_at
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []domain.CartItem
	for rows.Next() {
		var item domain.CartItem
		if err := rows.Scan(&item.ID, &item.UserID, &item.BookID, &item.AddedAt, &item.ExpiresAt); err != nil {
			return nil, err
		}
		out = append(out, item)
	}
	return out, rows.Err()
}

func (s *Store) AddCartItem(ctx context.Context, userID, bookID int64, expiresAt time.Time) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var stock int
	err = tx.QueryRow(ctx, `
		UPDATE books SET stock = stock - 1
		WHERE id = $1 AND stock > 0
		RETURNING stock
	`, bookID).Scan(&stock)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ErrOutOfStock
		}
		return err
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO cart_items (user_id, book_id, expires_at)
		VALUES ($1, $2, $3)
	`, userID, bookID, expiresAt)
	if err != nil {
		if isUniqueViolation(err) {
			return domain.ErrAlreadyInCart
		}
		return err
	}

	return tx.Commit(ctx)
}

func (s *Store) RemoveCartItem(ctx context.Context, userID, bookID int64) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	tag, err := tx.Exec(ctx, `
		DELETE FROM cart_items WHERE user_id = $1 AND book_id = $2
	`, userID, bookID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}

	_, err = tx.Exec(ctx, `UPDATE books SET stock = stock + 1 WHERE id = $1`, bookID)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (s *Store) ClearCart(ctx context.Context, userID int64) error {
	items, err := s.ListCartItems(ctx, userID)
	if err != nil {
		return err
	}
	for _, item := range items {
		if err := s.RemoveCartItem(ctx, userID, item.BookID); err != nil {
			return err
		}
	}
	return nil
}

func (s *Store) ReleaseExpiredCartItems(ctx context.Context) (int, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx)

	rows, err := tx.Query(ctx, `
		SELECT user_id, book_id
		FROM cart_items
		WHERE expires_at <= NOW()
		FOR UPDATE
	`)
	if err != nil {
		return 0, err
	}

	type pair struct {
		userID int64
		bookID int64
	}
	var expired []pair
	for rows.Next() {
		var p pair
		if err := rows.Scan(&p.userID, &p.bookID); err != nil {
			rows.Close()
			return 0, err
		}
		expired = append(expired, p)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return 0, err
	}

	for _, p := range expired {
		tag, err := tx.Exec(ctx, `
			DELETE FROM cart_items WHERE user_id = $1 AND book_id = $2
		`, p.userID, p.bookID)
		if err != nil {
			return 0, err
		}
		if tag.RowsAffected() == 0 {
			continue
		}
		if _, err := tx.Exec(ctx, `UPDATE books SET stock = stock + 1 WHERE id = $1`, p.bookID); err != nil {
			return 0, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return 0, err
	}
	return len(expired), nil
}

func (s *Store) CheckoutCart(ctx context.Context, userID int64) error {
	tag, err := s.pool.Exec(ctx, `DELETE FROM cart_items WHERE user_id = $1`, userID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrCartEmpty
	}
	return nil
}

func buildBooksQuery(filter domain.BookFilter, countOnly bool) (string, []any) {
	var b strings.Builder
	args := make([]any, 0, 4)

	if countOnly {
		b.WriteString("SELECT COUNT(*) FROM books WHERE stock > 0")
	} else {
		b.WriteString(`
			SELECT id, title, year_published, author, price_usd, category_id, stock, created_at
			FROM books WHERE stock > 0
		`)
	}

	if len(filter.CategoryIDs) > 0 {
		args = append(args, filter.CategoryIDs)
		fmt.Fprintf(&b, " AND category_id = ANY($%d)", len(args))
	}

	if !countOnly {
		b.WriteString(" ORDER BY id")
		if filter.Limit > 0 {
			args = append(args, filter.Limit)
			fmt.Fprintf(&b, " LIMIT $%d", len(args))
		}
		if filter.Offset > 0 {
			args = append(args, filter.Offset)
			fmt.Fprintf(&b, " OFFSET $%d", len(args))
		}
	}

	return b.String(), args
}

func scanUser(row pgx.Row) (domain.User, error) {
	var user domain.User
	if err := row.Scan(&user.ID, &user.Email, &user.PasswordHash, &user.IsAdmin, &user.CreatedAt); err != nil {
		return domain.User{}, mapNotFound(err)
	}
	return user, nil
}

func scanBook(row pgx.Row) (domain.Book, error) {
	var book domain.Book
	if err := row.Scan(&book.ID, &book.Title, &book.Year, &book.Author, &book.PriceUSD, &book.CategoryID, &book.Stock, &book.CreatedAt); err != nil {
		return domain.Book{}, mapNotFound(err)
	}
	return book, nil
}

type bookScanner interface {
	Scan(dest ...any) error
}

func scanBookRow(row bookScanner) (domain.Book, error) {
	var book domain.Book
	if err := row.Scan(&book.ID, &book.Title, &book.Year, &book.Author, &book.PriceUSD, &book.CategoryID, &book.Stock, &book.CreatedAt); err != nil {
		return domain.Book{}, err
	}
	return book, nil
}

func mapNotFound(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.ErrNotFound
	}
	return err
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}
