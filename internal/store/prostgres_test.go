package store

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/mrckurz/CI-CD-MCM/internal/model"
)

func newMockPostgresStore(t *testing.T) (*PostgresStore, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}

	cleanup := func() {
		_ = db.Close()
	}

	return &PostgresStore{DB: db}, mock, cleanup
}

func TestPostgresEnsureTable(t *testing.T) {
	s, mock, cleanup := newMockPostgresStore(t)
	defer cleanup()

	mock.ExpectExec("CREATE TABLE IF NOT EXISTS products").
		WillReturnResult(sqlmock.NewResult(0, 0))

	if err := s.EnsureTable(); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestPostgresGetAll(t *testing.T) {
	s, mock, cleanup := newMockPostgresStore(t)
	defer cleanup()

	rows := sqlmock.NewRows([]string{"id", "name", "price"}).
		AddRow(1, "Widget", 9.99)

	mock.ExpectQuery("SELECT id, name, price FROM products ORDER BY id").
		WillReturnRows(rows)

	products, err := s.GetAll()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(products) != 1 {
		t.Fatalf("expected 1 product, got %d", len(products))
	}

	if products[0].Name != "Widget" {
		t.Fatalf("expected Widget, got %s", products[0].Name)
	}
}

func TestPostgresGetByID(t *testing.T) {
	s, mock, cleanup := newMockPostgresStore(t)
	defer cleanup()

	rows := sqlmock.NewRows([]string{"id", "name", "price"}).
		AddRow(1, "Widget", 9.99)

	mock.ExpectQuery("SELECT id, name, price FROM products WHERE id = \\$1").
		WithArgs(1).
		WillReturnRows(rows)

	product, err := s.GetByID(1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if product.ID != 1 {
		t.Fatalf("expected ID 1, got %d", product.ID)
	}
}

func TestPostgresCreate(t *testing.T) {
	s, mock, cleanup := newMockPostgresStore(t)
	defer cleanup()

	rows := sqlmock.NewRows([]string{"id"}).AddRow(1)

	mock.ExpectQuery("INSERT INTO products").
		WithArgs("Widget", 9.99).
		WillReturnRows(rows)

	product, err := s.Create(model.Product{
		Name:  "Widget",
		Price: 9.99,
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if product.ID != 1 {
		t.Fatalf("expected ID 1, got %d", product.ID)
	}
}

func TestPostgresUpdate(t *testing.T) {
	s, mock, cleanup := newMockPostgresStore(t)
	defer cleanup()

	mock.ExpectExec("UPDATE products SET").
		WithArgs("Updated", 19.99, 1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	product, err := s.Update(1, model.Product{
		Name:  "Updated",
		Price: 19.99,
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if product.ID != 1 {
		t.Fatalf("expected ID 1, got %d", product.ID)
	}
}

func TestPostgresDelete(t *testing.T) {
	s, mock, cleanup := newMockPostgresStore(t)
	defer cleanup()

	mock.ExpectExec("DELETE FROM products WHERE id = \\$1").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := s.Delete(1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}
