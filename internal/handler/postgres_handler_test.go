package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gorilla/mux"
	"github.com/mrckurz/CI-CD-MCM/internal/store"
)

func setupPostgresRouter(t *testing.T) (*mux.Router, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}

	pgStore := &store.PostgresStore{DB: db}
	h := NewPostgresHandler(pgStore)

	r := mux.NewRouter()
	h.RegisterRoutes(r)

	cleanup := func() {
		_ = db.Close()
	}

	return r, mock, cleanup
}

func TestPostgresHealthEndpoint(t *testing.T) {
	r, _, cleanup := setupPostgresRouter(t)
	defer cleanup()

	req := httptest.NewRequest("GET", "/health", nil)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
}

func TestPostgresGetProductsEmpty(t *testing.T) {
	r, mock, cleanup := setupPostgresRouter(t)
	defer cleanup()

	rows := sqlmock.NewRows([]string{"id", "name", "price"})

	mock.ExpectQuery("SELECT id, name, price FROM products ORDER BY id").
		WillReturnRows(rows)

	req := httptest.NewRequest("GET", "/products", nil)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
}

func TestPostgresGetProduct(t *testing.T) {
	r, mock, cleanup := setupPostgresRouter(t)
	defer cleanup()

	rows := sqlmock.NewRows([]string{"id", "name", "price"}).
		AddRow(1, "Widget", 9.99)

	mock.ExpectQuery("SELECT id, name, price FROM products WHERE id = \\$1").
		WithArgs(1).
		WillReturnRows(rows)

	req := httptest.NewRequest("GET", "/products/1", nil)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
}

func TestPostgresGetProductNotFound(t *testing.T) {
	r, mock, cleanup := setupPostgresRouter(t)
	defer cleanup()

	rows := sqlmock.NewRows([]string{"id", "name", "price"})

	mock.ExpectQuery("SELECT id, name, price FROM products WHERE id = \\$1").
		WithArgs(999).
		WillReturnRows(rows)

	req := httptest.NewRequest("GET", "/products/999", nil)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", rr.Code)
	}
}

func TestPostgresCreateProduct(t *testing.T) {
	r, mock, cleanup := setupPostgresRouter(t)
	defer cleanup()

	rows := sqlmock.NewRows([]string{"id"}).AddRow(1)

	mock.ExpectQuery("INSERT INTO products").
		WithArgs("Widget", 9.99).
		WillReturnRows(rows)

	req := httptest.NewRequest("POST", "/products", strings.NewReader(`{"name":"Widget","price":9.99}`))
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d", rr.Code)
	}
}

func TestPostgresCreateProductInvalidJSON(t *testing.T) {
	r, _, cleanup := setupPostgresRouter(t)
	defer cleanup()

	req := httptest.NewRequest("POST", "/products", strings.NewReader(`invalid-json`))
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rr.Code)
	}
}

func TestPostgresCreateProductInvalidProduct(t *testing.T) {
	r, _, cleanup := setupPostgresRouter(t)
	defer cleanup()

	req := httptest.NewRequest("POST", "/products", strings.NewReader(`{"name":"","price":-1}`))
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rr.Code)
	}
}

func TestPostgresUpdateProduct(t *testing.T) {
	r, mock, cleanup := setupPostgresRouter(t)
	defer cleanup()

	mock.ExpectExec("UPDATE products SET").
		WithArgs("Updated", 19.99, 1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	req := httptest.NewRequest("PUT", "/products/1", strings.NewReader(`{"name":"Updated","price":19.99}`))
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
}

func TestPostgresUpdateProductInvalidJSON(t *testing.T) {
	r, _, cleanup := setupPostgresRouter(t)
	defer cleanup()

	req := httptest.NewRequest("PUT", "/products/1", strings.NewReader(`invalid-json`))
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rr.Code)
	}
}

func TestPostgresUpdateProductNotFound(t *testing.T) {
	r, mock, cleanup := setupPostgresRouter(t)
	defer cleanup()

	mock.ExpectExec("UPDATE products SET").
		WithArgs("Updated", 19.99, 999).
		WillReturnResult(sqlmock.NewResult(0, 0))

	req := httptest.NewRequest("PUT", "/products/999", strings.NewReader(`{"name":"Updated","price":19.99}`))
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", rr.Code)
	}
}

func TestPostgresDeleteProduct(t *testing.T) {
	r, mock, cleanup := setupPostgresRouter(t)
	defer cleanup()

	mock.ExpectExec("DELETE FROM products WHERE id = \\$1").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	req := httptest.NewRequest("DELETE", "/products/1", nil)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
}

func TestPostgresDeleteProductNotFound(t *testing.T) {
	r, mock, cleanup := setupPostgresRouter(t)
	defer cleanup()

	mock.ExpectExec("DELETE FROM products WHERE id = \\$1").
		WithArgs(999).
		WillReturnResult(sqlmock.NewResult(0, 0))

	req := httptest.NewRequest("DELETE", "/products/999", nil)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", rr.Code)
	}
}
