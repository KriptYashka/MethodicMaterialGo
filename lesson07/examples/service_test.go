package main

import (
	"context"
	"testing"
)

type mockRepo struct {
	products []Product
	getFn    func(int) (*Product, error)
	createFn func(*Product) (int64, error)
	err      error
}

func (m *mockRepo) List(ctx context.Context) ([]Product, error) {
	return m.products, m.err
}

func (m *mockRepo) GetByID(ctx context.Context, id int) (*Product, error) {
	if m.getFn != nil {
		return m.getFn(id)
	}
	for _, p := range m.products {
		if p.ID == id {
			return &p, nil
		}
	}
	return nil, m.err
}

func (m *mockRepo) Create(ctx context.Context, p *Product) (int64, error) {
	if m.createFn != nil {
		return m.createFn(p)
	}
	p.ID = len(m.products) + 1
	m.products = append(m.products, *p)
	return int64(p.ID), nil
}

func (m *mockRepo) Update(ctx context.Context, p *Product) error {
	return m.err
}

func (m *mockRepo) Delete(ctx context.Context, id int) error {
	return m.err
}

func TestListProducts(t *testing.T) {
	repo := &mockRepo{
		products: []Product{
			{ID: 1, Name: "Widget", Price: 9.99, Qty: 10},
			{ID: 2, Name: "Gadget", Price: 24.99, Qty: 5},
		},
	}
	svc := NewProductService(repo)

	products, err := svc.List(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(products) != 2 {
		t.Fatalf("expected 2, got %d", len(products))
	}
}

func TestListProductsEmpty(t *testing.T) {
	repo := &mockRepo{products: []Product{}}
	svc := NewProductService(repo)

	products, err := svc.List(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if products == nil {
		t.Fatal("expected empty slice, got nil")
	}
	if len(products) != 0 {
		t.Fatalf("expected 0, got %d", len(products))
	}
}

func TestGetProduct(t *testing.T) {
	repo := &mockRepo{
		getFn: func(id int) (*Product, error) {
			return &Product{ID: id, Name: "Test", Price: 10, Qty: 1}, nil
		},
	}
	svc := NewProductService(repo)

	p, err := svc.Get(context.Background(), 42)
	if err != nil {
		t.Fatal(err)
	}
	if p.Name != "Test" {
		t.Fatalf("expected Test, got %s", p.Name)
	}
}

func TestCreateProduct(t *testing.T) {
	repo := &mockRepo{}
	svc := NewProductService(repo)

	p, err := svc.Create(context.Background(), &Product{Name: "New", Price: 5, Qty: 3})
	if err != nil {
		t.Fatal(err)
	}
	if p.ID != 1 {
		t.Fatalf("expected ID 1, got %d", p.ID)
	}
}
