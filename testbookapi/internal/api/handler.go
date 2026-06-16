package api

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"lessonsgo/testbookapi/internal/auth"
	"lessonsgo/testbookapi/internal/book"
	"lessonsgo/testbookapi/internal/cart"
	"lessonsgo/testbookapi/internal/category"
	"lessonsgo/testbookapi/internal/domain"
)

type Handler struct {
	authSvc     *auth.Service
	categorySvc *category.Service
	bookSvc     *book.Service
	cartSvc     *cart.Service
	logger      *slog.Logger
}

func NewHandler(
	authSvc *auth.Service,
	categorySvc *category.Service,
	bookSvc *book.Service,
	cartSvc *cart.Service,
	logger *slog.Logger,
) *Handler {
	return &Handler{
		authSvc:     authSvc,
		categorySvc: categorySvc,
		bookSvc:     bookSvc,
		cartSvc:     cartSvc,
		logger:      logger,
	}
}

func (h *Handler) Register(mux *http.ServeMux) {
	mux.HandleFunc("POST /auth/register", h.register)
	mux.HandleFunc("POST /auth/login", h.login)

	mux.HandleFunc("GET /categories", h.listCategories)
	mux.HandleFunc("GET /categories/{id}", h.getCategory)
	mux.Handle("POST /categories", h.requireAdmin(http.HandlerFunc(h.createCategory)))
	mux.Handle("PUT /categories/{id}", h.requireAdmin(http.HandlerFunc(h.updateCategory)))
	mux.Handle("DELETE /categories/{id}", h.requireAdmin(http.HandlerFunc(h.deleteCategory)))

	mux.HandleFunc("GET /books", h.listBooks)
	mux.HandleFunc("GET /books/{id}", h.getBook)
	mux.Handle("POST /books", h.requireAdmin(http.HandlerFunc(h.createBook)))
	mux.Handle("PUT /books/{id}", h.requireAdmin(http.HandlerFunc(h.updateBook)))
	mux.Handle("DELETE /books/{id}", h.requireAdmin(http.HandlerFunc(h.deleteBook)))

	mux.Handle("GET /cart", h.requireAuth(http.HandlerFunc(h.getCart)))
	mux.Handle("POST /cart/items", h.requireAuth(http.HandlerFunc(h.addCartItem)))
	mux.Handle("DELETE /cart/items/{book_id}", h.requireAuth(http.HandlerFunc(h.removeCartItem)))
	mux.Handle("POST /cart/checkout", h.requireAuth(http.HandlerFunc(h.checkout)))
}

func (h *Handler) register(w http.ResponseWriter, r *http.Request) {
	var req auth.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	resp, err := h.authSvc.Register(r.Context(), req)
	if err != nil {
		h.handleError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, resp)
}

func (h *Handler) login(w http.ResponseWriter, r *http.Request) {
	var req auth.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	resp, err := h.authSvc.Login(r.Context(), req)
	if err != nil {
		h.handleError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, resp)
}

func (h *Handler) listCategories(w http.ResponseWriter, r *http.Request) {
	cats, err := h.categorySvc.List(r.Context())
	if err != nil {
		h.handleError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"categories": cats})
}

func (h *Handler) getCategory(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid category id")
		return
	}

	cat, err := h.categorySvc.Get(r.Context(), id)
	if err != nil {
		h.handleError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, cat)
}

func (h *Handler) createCategory(w http.ResponseWriter, r *http.Request) {
	var req category.CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	cat, err := h.categorySvc.Create(r.Context(), req)
	if err != nil {
		h.handleError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, cat)
}

func (h *Handler) updateCategory(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid category id")
		return
	}

	var req category.UpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	cat, err := h.categorySvc.Update(r.Context(), id, req)
	if err != nil {
		h.handleError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, cat)
}

func (h *Handler) deleteCategory(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid category id")
		return
	}

	if err := h.categorySvc.Delete(r.Context(), id); err != nil {
		h.handleError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) listBooks(w http.ResponseWriter, r *http.Request) {
	categoryIDs, err := parseCategoryFilter(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid category filter")
		return
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	result, err := h.bookSvc.List(r.Context(), categoryIDs, limit, offset)
	if err != nil {
		h.handleError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func (h *Handler) getBook(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid book id")
		return
	}

	b, err := h.bookSvc.Get(r.Context(), id)
	if err != nil {
		h.handleError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, b)
}

func (h *Handler) createBook(w http.ResponseWriter, r *http.Request) {
	var req book.CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	b, err := h.bookSvc.Create(r.Context(), req)
	if err != nil {
		h.handleError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, b)
}

func (h *Handler) updateBook(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid book id")
		return
	}

	var req book.UpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	b, err := h.bookSvc.Update(r.Context(), id, req)
	if err != nil {
		h.handleError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, b)
}

func (h *Handler) deleteBook(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid book id")
		return
	}

	if err := h.bookSvc.Delete(r.Context(), id); err != nil {
		h.handleError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) getCart(w http.ResponseWriter, r *http.Request) {
	userID := userIDFromContext(r.Context())
	resp, err := h.cartSvc.List(r.Context(), userID)
	if err != nil {
		h.handleError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, resp)
}

func (h *Handler) addCartItem(w http.ResponseWriter, r *http.Request) {
	var body struct {
		BookID int64 `json:"book_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	userID := userIDFromContext(r.Context())
	if err := h.cartSvc.Add(r.Context(), userID, body.BookID); err != nil {
		h.handleError(w, err)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) removeCartItem(w http.ResponseWriter, r *http.Request) {
	bookID, err := parseID(r.PathValue("book_id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid book id")
		return
	}

	userID := userIDFromContext(r.Context())
	if err := h.cartSvc.Remove(r.Context(), userID, bookID); err != nil {
		h.handleError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) checkout(w http.ResponseWriter, r *http.Request) {
	userID := userIDFromContext(r.Context())
	if err := h.cartSvc.Checkout(r.Context(), userID); err != nil {
		h.handleError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "completed"})
}

func (h *Handler) requireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, err := h.parseBearer(r)
		if err != nil {
			h.handleError(w, err)
			return
		}
		ctx := context.WithValue(r.Context(), userIDKey, claims.UserID)
		ctx = context.WithValue(ctx, isAdminKey, claims.IsAdmin)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (h *Handler) requireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, err := h.parseBearer(r)
		if err != nil {
			h.handleError(w, err)
			return
		}
		if !claims.IsAdmin {
			writeError(w, http.StatusForbidden, domain.ErrForbidden.Error())
			return
		}
		ctx := context.WithValue(r.Context(), userIDKey, claims.UserID)
		ctx = context.WithValue(ctx, isAdminKey, claims.IsAdmin)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (h *Handler) parseBearer(r *http.Request) (auth.Claims, error) {
	header := r.Header.Get("Authorization")
	if header == "" {
		return auth.Claims{}, domain.ErrUnauthorized
	}
	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") || parts[1] == "" {
		return auth.Claims{}, domain.ErrUnauthorized
	}
	return h.authSvc.ParseToken(parts[1])
}

func (h *Handler) handleError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrNotFound):
		writeError(w, http.StatusNotFound, err.Error())
	case errors.Is(err, domain.ErrUnauthorized):
		writeError(w, http.StatusUnauthorized, err.Error())
	case errors.Is(err, domain.ErrForbidden):
		writeError(w, http.StatusForbidden, err.Error())
	case errors.Is(err, domain.ErrInvalidInput),
		errors.Is(err, domain.ErrEmailTaken),
		errors.Is(err, domain.ErrInvalidCredentials),
		errors.Is(err, domain.ErrOutOfStock),
		errors.Is(err, domain.ErrAlreadyInCart),
		errors.Is(err, domain.ErrCategoryInUse),
		errors.Is(err, domain.ErrStockNotEditable),
		errors.Is(err, domain.ErrCartEmpty),
		errors.Is(err, domain.ErrBookNotPurchasable):
		writeError(w, http.StatusBadRequest, err.Error())
	default:
		h.logger.Error("request failed", "error", err)
		writeError(w, http.StatusInternalServerError, "internal server error")
	}
}

func parseID(raw string) (int64, error) {
	return strconv.ParseInt(raw, 10, 64)
}

func parseCategoryFilter(r *http.Request) ([]int64, error) {
	values := r.URL.Query()["category"]
	if len(values) == 0 {
		values = r.URL.Query()["categories"]
	}
	if len(values) == 0 {
		return nil, nil
	}

	ids := make([]int64, 0, len(values))
	for _, v := range values {
		for _, part := range strings.Split(v, ",") {
			part = strings.TrimSpace(part)
			if part == "" {
				continue
			}
			id, err := strconv.ParseInt(part, 10, 64)
			if err != nil {
				return nil, err
			}
			ids = append(ids, id)
		}
	}
	return ids, nil
}

type contextKey string

const (
	userIDKey  contextKey = "user_id"
	isAdminKey contextKey = "is_admin"
)

func userIDFromContext(ctx context.Context) int64 {
	id, _ := ctx.Value(userIDKey).(int64)
	return id
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

// StartCartCleanup runs a background job that releases expired cart reservations.
func StartCartCleanup(ctx context.Context, cartSvc *cart.Service, logger *slog.Logger, interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				released, err := cartSvc.ReleaseExpired(ctx)
				if err != nil {
					logger.Error("cart cleanup failed", "error", err)
					continue
				}
				if released > 0 {
					logger.Info("released expired cart items", "count", released)
				}
			}
		}
	}()
}
