package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Product struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
	Qty   int     `json:"qty"`
}

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) Init(ctx context.Context) error {
	_, err := s.db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS products (
			id    INTEGER PRIMARY KEY AUTOINCREMENT,
			name  TEXT    NOT NULL,
			price REAL    NOT NULL DEFAULT 0,
			qty   INTEGER NOT NULL DEFAULT 0
		)
	`)
	return err
}

func (s *Store) Insert(ctx context.Context, p *Product) (int64, error) {
	res, err := s.db.ExecContext(ctx,
		`INSERT INTO products (name, price, qty) VALUES (?, ?, ?)`,
		p.Name, p.Price, p.Qty,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (s *Store) GetByID(ctx context.Context, id int) (*Product, error) {
	var p Product
	err := s.db.QueryRowContext(ctx,
		`SELECT id, name, price, qty FROM products WHERE id = ?`, id,
	).Scan(&p.ID, &p.Name, &p.Price, &p.Qty)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product %d not found", id)
		}
		return nil, err
	}
	return &p, nil
}

func (s *Store) List(ctx context.Context) ([]Product, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, name, price, qty FROM products ORDER BY id`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var p Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Qty); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, rows.Err()
}

func (s *Store) UpdateQty(ctx context.Context, id, delta int) error {
	res, err := s.db.ExecContext(ctx,
		`UPDATE products SET qty = qty + ? WHERE id = ?`, delta, id,
	)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("product %d not found", id)
	}
	return nil
}

func (s *Store) Delete(ctx context.Context, id int) error {
	res, err := s.db.ExecContext(ctx, `DELETE FROM products WHERE id = ?`, id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("product %d not found", id)
	}
	return nil
}

func (s *Store) Buy(ctx context.Context, userID, productID, qty int) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var currentQty int
	err = tx.QueryRow(
		`SELECT qty FROM products WHERE id = ? FOR UPDATE`, productID,
	).Scan(&currentQty)
	if err != nil {
		return err
	}
	if currentQty < qty {
		return fmt.Errorf("insufficient stock: have %d, need %d", currentQty, qty)
	}

	_, err = tx.Exec(
		`UPDATE products SET qty = qty - ? WHERE id = ?`, qty, productID,
	)
	if err != nil {
		return err
	}

	_, err = tx.Exec(
		`INSERT INTO orders (user_id, product_id, qty, created_at) VALUES (?, ?, ?, ?)`,
		userID, productID, qty, time.Now(),
	)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func main() {
	db, err := sql.Open("sqlite3", "file:shop.db?cache=shared")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	db.SetMaxOpenConns(1)

	ctx := context.Background()
	store := NewStore(db)

	if err := store.Init(ctx); err != nil {
		log.Fatal(err)
	}

	id1, _ := store.Insert(ctx, &Product{Name: "Widget", Price: 9.99, Qty: 100})
	id2, _ := store.Insert(ctx, &Product{Name: "Gadget", Price: 24.99, Qty: 50})
	fmt.Printf("Inserted: %d, %d\n", id1, id2)

	p, _ := store.GetByID(ctx, 1)
	fmt.Printf("Product 1: %+v\n", p)

	store.UpdateQty(ctx, 1, -5)

	all, _ := store.List(ctx)
	fmt.Println("All products:")
	for _, pr := range all {
		fmt.Printf("  %d: %s $%.2f (qty: %d)\n", pr.ID, pr.Name, pr.Price, pr.Qty)
	}

	store.Delete(ctx, 2)

	all, _ = store.List(ctx)
	fmt.Println("After delete:")
	for _, pr := range all {
		fmt.Printf("  %d: %s $%.2f (qty: %d)\n", pr.ID, pr.Name, pr.Price, pr.Qty)
	}
}
