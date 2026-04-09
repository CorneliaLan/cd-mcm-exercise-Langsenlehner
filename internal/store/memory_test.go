package store

import (
	"testing"

	"github.com/mrckurz/CI-CD-MCM/internal/model"
)

func TestCreateAndGet(t *testing.T) {
	s := NewMemoryStore()

	created := s.Create(model.Product{
		Name:  "Keyboard",
		Price: 49.99,
	})

	if created.ID == 0 {
		t.Error("expected created product to have a non-zero ID")
	}

	got, err := s.GetByID(created.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if got.ID != created.ID {
		t.Errorf("expected ID %d, got %d", created.ID, got.ID)
	}

	if got.Name != created.Name {
		t.Errorf("expected Name %q, got %q", created.Name, got.Name)
	}

	if got.Price != created.Price {
		t.Errorf("expected Price %v, got %v", created.Price, got.Price)
	}
}

func TestGetAllEmpty(t *testing.T) {
	s := NewMemoryStore()

	products := s.GetAll()
	if len(products) != 0 {
		t.Errorf("expected 0 products, got %d", len(products))
	}
}

func TestDeleteNonExistent(t *testing.T) {
	s := NewMemoryStore()

	err := s.Delete(999)
	if err != ErrNotFound {
		t.Errorf("expected ErrNotFound when deleting non-existent product, got %v", err)
	}
}

func TestUpdateExistingProduct(t *testing.T) {
	s := NewMemoryStore()

	created := s.Create(model.Product{
		Name:  "Mouse",
		Price: 19.99,
	})

	updated, err := s.Update(created.ID, model.Product{
		Name:  "Gaming Mouse",
		Price: 29.99,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if updated.ID != created.ID {
		t.Errorf("expected ID %d, got %d", created.ID, updated.ID)
	}

	if updated.Name != "Gaming Mouse" {
		t.Errorf("expected Name %q, got %q", "Gaming Mouse", updated.Name)
	}

	if updated.Price != 29.99 {
		t.Errorf("expected Price %v, got %v", 29.99, updated.Price)
	}

	got, err := s.GetByID(created.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if got.Name != "Gaming Mouse" {
		t.Errorf("expected Name %q, got %q", "Gaming Mouse", got.Name)
	}

	if got.Price != 29.99 {
		t.Errorf("expected Price %v, got %v", 29.99, got.Price)
	}
}

func TestDeleteExistingProduct(t *testing.T) {
	s := NewMemoryStore()

	created := s.Create(model.Product{
		Name:  "Monitor",
		Price: 199.99,
	})

	err := s.Delete(created.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	_, err = s.GetByID(created.ID)
	if err != ErrNotFound {
		t.Errorf("expected ErrNotFound after delete, got %v", err)
	}
}

func TestGetByIDInvalidID(t *testing.T) {
	s := NewMemoryStore()

	_, err := s.GetByID(999)
	if err != ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}