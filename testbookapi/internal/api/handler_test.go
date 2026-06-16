package api_test

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"sync"
	"testing"
	"time"

	"lessonsgo/testbookapi/internal/api"
	"lessonsgo/testbookapi/internal/auth"
	"lessonsgo/testbookapi/internal/book"
	"lessonsgo/testbookapi/internal/cart"
	"lessonsgo/testbookapi/internal/category"
	"lessonsgo/testbookapi/internal/storage/memory"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestServer(t *testing.T) (*httptest.Server, *memory.Store) {
	t.Helper()

	store := memory.New()
	authSvc := auth.NewService(store, "test-secret", time.Hour)
	categorySvc := category.NewService(store)
	bookSvc := book.NewService(store)
	cartSvc := cart.NewService(store)
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
	handler := api.NewHandler(authSvc, categorySvc, bookSvc, cartSvc, logger)

	mux := http.NewServeMux()
	handler.Register(mux)
	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)
	return server, store
}

func registerUser(t *testing.T, server *httptest.Server, email, password string) string {
	t.Helper()
	resp := postJSON(t, server.URL+"/auth/register", map[string]string{
		"email":    email,
		"password": password,
	}, "")
	require.Equal(t, http.StatusCreated, resp.StatusCode)
	defer resp.Body.Close()

	var body map[string]string
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&body))
	return body["token"]
}

func loginUser(t *testing.T, server *httptest.Server, email, password string) string {
	t.Helper()
	resp := postJSON(t, server.URL+"/auth/login", map[string]string{
		"email":    email,
		"password": password,
	}, "")
	require.Equal(t, http.StatusOK, resp.StatusCode)
	defer resp.Body.Close()

	var body map[string]string
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&body))
	return body["token"]
}

func postJSON(t *testing.T, url string, body any, token string) *http.Response {
	t.Helper()
	payload, err := json.Marshal(body)
	require.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(payload))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	return resp
}

func authRequest(t *testing.T, method, url, token string, body any) *http.Response {
	t.Helper()
	var reader io.Reader
	if body != nil {
		payload, err := json.Marshal(body)
		require.NoError(t, err)
		reader = bytes.NewReader(payload)
	}

	req, err := http.NewRequest(method, url, reader)
	require.NoError(t, err)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	return resp
}

func TestBookShopAPI(t *testing.T) {
	server, store := setupTestServer(t)

	userToken := registerUser(t, server, "user@example.com", "password123")
	adminToken := registerUser(t, server, "admin@example.com", "password123")
	require.NoError(t, store.SetAdmin(2, true))
	adminToken = loginUser(t, server, "admin@example.com", "password123")

	t.Run("anonymous can list books with category filter", func(t *testing.T) {
		cat1 := createCategory(t, server, adminToken, "Fiction")
		cat2 := createCategory(t, server, adminToken, "Science")
		book1 := createBook(t, server, adminToken, "Book A", cat1.ID, 2)
		createBook(t, server, adminToken, "Book B", cat2.ID, 1)

		resp, err := http.Get(server.URL + "/books?category=" + itoa(cat1.ID) + "&category=" + itoa(cat2.ID))
		require.NoError(t, err)
		defer resp.Body.Close()
		require.Equal(t, http.StatusOK, resp.StatusCode)

		var result book.ListResult
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
		assert.Equal(t, 2, result.Total)
		assert.Len(t, result.Books, 2)

		soldOut := createBook(t, server, adminToken, "Sold Out", cat1.ID, 0)
		resp2, err := http.Get(server.URL + "/books")
		require.NoError(t, err)
		defer resp2.Body.Close()
		require.NoError(t, json.NewDecoder(resp2.Body).Decode(&result))
		for _, b := range result.Books {
			assert.NotEqual(t, soldOut.ID, b.ID)
		}

		getResp, err := http.Get(server.URL + "/books/" + itoa(book1.ID))
		require.NoError(t, err)
		defer getResp.Body.Close()
		assert.Equal(t, http.StatusOK, getResp.StatusCode)
	})

	t.Run("admin CRUD and stock is immutable", func(t *testing.T) {
		cat := createCategory(t, server, adminToken, "History")
		created := createBook(t, server, adminToken, "History Book", cat.ID, 3)

		updateResp := authRequest(t, http.MethodPut, server.URL+"/books/"+itoa(created.ID), adminToken, map[string]any{
			"title":          "Updated Title",
			"year_published": 2020,
			"author":         "Author",
			"price_usd":      19.99,
			"category_id":    cat.ID,
			"stock":          99,
		})
		require.Equal(t, http.StatusBadRequest, updateResp.StatusCode)
		updateResp.Body.Close()

		updateResp = authRequest(t, http.MethodPut, server.URL+"/books/"+itoa(created.ID), adminToken, map[string]any{
			"title":          "Updated Title",
			"year_published": 2020,
			"author":         "Author",
			"price_usd":      19.99,
			"category_id":    cat.ID,
		})
		require.Equal(t, http.StatusOK, updateResp.StatusCode)
		updateResp.Body.Close()

		forbiddenResp := authRequest(t, http.MethodPost, server.URL+"/books", userToken, map[string]any{
			"title": "Nope", "year_published": 2021, "author": "X", "price_usd": 1, "category_id": cat.ID, "stock": 1,
		})
		require.Equal(t, http.StatusForbidden, forbiddenResp.StatusCode)
		forbiddenResp.Body.Close()
	})

	t.Run("cart checkout reduces stock and prevents race", func(t *testing.T) {
		cat := createCategory(t, server, adminToken, "Thriller")
		onlyOne := createBook(t, server, adminToken, "Last Copy", cat.ID, 1)

		user2Token := registerUser(t, server, "buyer2@example.com", "password123")

		add1 := postJSON(t, server.URL+"/cart/items", map[string]int64{"book_id": onlyOne.ID}, userToken)
		require.Equal(t, http.StatusCreated, add1.StatusCode)
		add1.Body.Close()

		add2 := postJSON(t, server.URL+"/cart/items", map[string]int64{"book_id": onlyOne.ID}, user2Token)
		require.Equal(t, http.StatusBadRequest, add2.StatusCode)
		add2.Body.Close()

		checkout := postJSON(t, server.URL+"/cart/checkout", nil, userToken)
		require.Equal(t, http.StatusOK, checkout.StatusCode)
		checkout.Body.Close()

		listResp, err := http.Get(server.URL + "/books")
		require.NoError(t, err)
		defer listResp.Body.Close()
		var result book.ListResult
		require.NoError(t, json.NewDecoder(listResp.Body).Decode(&result))
		for _, b := range result.Books {
			assert.NotEqual(t, onlyOne.ID, b.ID)
		}
	})

	t.Run("concurrent add to cart for last item", func(t *testing.T) {
		cat := createCategory(t, server, adminToken, "Concurrency")
		bookItem := createBook(t, server, adminToken, "Race Book", cat.ID, 1)

		tokenA := registerUser(t, server, "race-a@example.com", "password123")
		tokenB := registerUser(t, server, "race-b@example.com", "password123")

		var wg sync.WaitGroup
		results := make(chan int, 2)
		for _, token := range []string{tokenA, tokenB} {
			wg.Add(1)
			go func(tok string) {
				defer wg.Done()
				resp := postJSON(t, server.URL+"/cart/items", map[string]int64{"book_id": bookItem.ID}, tok)
				results <- resp.StatusCode
				resp.Body.Close()
			}(token)
		}
		wg.Wait()
		close(results)

		var codes []int
		for code := range results {
			codes = append(codes, code)
		}
		assert.Contains(t, codes, http.StatusCreated)
		assert.Contains(t, codes, http.StatusBadRequest)
	})

	t.Run("expired cart releases stock", func(t *testing.T) {
		cat := createCategory(t, server, adminToken, "Expiry")
		bookItem := createBook(t, server, adminToken, "Expiring", cat.ID, 1)

		addResp := postJSON(t, server.URL+"/cart/items", map[string]int64{"book_id": bookItem.ID}, userToken)
		require.Equal(t, http.StatusCreated, addResp.StatusCode)
		addResp.Body.Close()

		items, err := store.ListCartItems(t.Context(), 1)
		require.NoError(t, err)
		require.Len(t, items, 1)
		require.NoError(t, store.SetCartItemExpiry(items[0].ID, time.Now().UTC().Add(-time.Minute)))

		released, err := store.ReleaseExpiredCartItems(t.Context())
		require.NoError(t, err)
		assert.Equal(t, 1, released)

		getResp, err := http.Get(server.URL + "/books/" + itoa(bookItem.ID))
		require.NoError(t, err)
		defer getResp.Body.Close()
		assert.Equal(t, http.StatusOK, getResp.StatusCode)
	})
}

func createCategory(t *testing.T, server *httptest.Server, token, name string) struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
} {
	t.Helper()
	resp := postJSON(t, server.URL+"/categories", map[string]string{"name": name}, token)
	require.Equal(t, http.StatusCreated, resp.StatusCode)
	defer resp.Body.Close()

	var cat struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	}
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&cat))
	return cat
}

func createBook(t *testing.T, server *httptest.Server, token, title string, categoryID int64, stock int) domainBook {
	t.Helper()
	resp := postJSON(t, server.URL+"/books", map[string]any{
		"title":          title,
		"year_published": 2024,
		"author":         "Test Author",
		"price_usd":      12.5,
		"category_id":    categoryID,
		"stock":          stock,
	}, token)
	require.Equal(t, http.StatusCreated, resp.StatusCode)
	defer resp.Body.Close()

	var b domainBook
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&b))
	return b
}

type domainBook struct {
	ID         int64   `json:"id"`
	Title      string  `json:"title"`
	CategoryID int64   `json:"category_id"`
	Stock      int     `json:"stock"`
	PriceUSD   float64 `json:"price_usd"`
}

func itoa(v int64) string {
	return strconv.FormatInt(v, 10)
}
