package store

import (
	"testing"

	"github.com/mrckurz/CI-CD-MCM/internal/model"
)

func TestCreateAndGet(t *testing.T) {
	s := NewMemoryStore()

	p := model.Product{Name: "Widget", Price: 9.99}
	created := s.Create(p)

	got, err := s.GetByID(created.ID)
	if err != nil {
		t.Fatalf("expected product, got error: %v", err)
	}

	if got.Name != "Widget" || got.Price != 9.99 {
		t.Errorf("unexpected product: %+v", got)
	}
}

func TestGetAllEmpty(t *testing.T) {
	s := NewMemoryStore()
	products := s.GetAll()
	if len(products) != 0 {
		t.Errorf("expected 0 products, got %d", len(products))
	}
}

func TestGetByIDInvalid(t *testing.T) {
	s := NewMemoryStore()

	_, err := s.GetByID(999)
	if err != ErrNotFound {
		t.Error("expected ErrNotFound")
	}
}

func TestUpdateExisting(t *testing.T) {
	s := NewMemoryStore()

	created := s.Create(model.Product{Name: "Old", Price: 1.00})

	updated, err := s.Update(created.ID, model.Product{Name: "New", Price: 2.00})
	if err != nil {
		t.Fatalf("expected update to succeed: %v", err)
	}

	if updated.Name != "New" || updated.Price != 2.00 {
		t.Errorf("unexpected updated product: %+v", updated)
	}
}

func TestUpdateNonExistent(t *testing.T) {
	s := NewMemoryStore()

	_, err := s.Update(999, model.Product{Name: "New", Price: 2.00})
	if err != ErrNotFound {
		t.Error("expected ErrNotFound")
	}
}

func TestDeleteExisting(t *testing.T) {
	s := NewMemoryStore()

	created := s.Create(model.Product{Name: "Widget", Price: 9.99})

	err := s.Delete(created.ID)
	if err != nil {
		t.Fatalf("expected delete to succeed: %v", err)
	}

	_, err = s.GetByID(created.ID)
	if err != ErrNotFound {
		t.Error("expected product to be deleted")
	}
}

func TestDeleteNonExistent(t *testing.T) {
	s := NewMemoryStore()
	err := s.Delete(999)
	if err != ErrNotFound {
		t.Error("expected ErrNotFound when deleting non-existent product")
	}
}
