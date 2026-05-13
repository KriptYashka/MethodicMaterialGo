package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Product struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
	Qty   int     `json:"qty"`
}

// ----- Repository layer -----

type ProductRepository interface {
	List(ctx context.Context) ([]Product, error)
	GetByID(ctx context.Context, id int) (*Product, error)
	Create(ctx context.Context, p *Product) (int64, error)
	Update(ctx context.Context, p *Product) error
	Delete(ctx context.Context, id int) error
}

type SQLiteRepo struct {
	db *sql.DB
}

func NewSQLiteRepo(db *sql.DB) *SQLiteRepo {
	return &SQLiteRepo{db: db}
}

func (r *SQLiteRepo) Init(ctx context.Context) error {
	_, err := r.db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS products (
			id    INTEGER PRIMARY KEY AUTOINCREMENT,
			name  TEXT    NOT NULL,
			price REAL    NOT NULL DEFAULT 0,
			qty   INTEGER NOT NULL DEFAULT 0
		)
	`)
	return err
}

func (r *SQLiteRepo) List(ctx context.Context) ([]Product, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, name, price, qty FROM products ORDER BY id`)
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

func (r *SQLiteRepo) GetByID(ctx context.Context, id int) (*Product, error) {
	var p Product
	err := r.db.QueryRowContext(ctx,
		`SELECT id, name, price, qty FROM products WHERE id = ?`, id,
	).Scan(&p.ID, &p.Name, &p.Price, &p.Qty)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *SQLiteRepo) Create(ctx context.Context, p *Product) (int64, error) {
	res, err := r.db.ExecContext(ctx,
		`INSERT INTO products (name, price, qty) VALUES (?, ?, ?)`,
		p.Name, p.Price, p.Qty,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *SQLiteRepo) Update(ctx context.Context, p *Product) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE products SET name=?, price=?, qty=? WHERE id=?`,
		p.Name, p.Price, p.Qty, p.ID,
	)
	return err
}

func (r *SQLiteRepo) Delete(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM products WHERE id=?`, id)
	return err
}

// ----- Service layer -----

type ProductService struct {
	repo ProductRepository
}

func NewProductService(repo ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) List(ctx context.Context) ([]Product, error) {
	products, err := s.repo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("list: %w", err)
	}
	if products == nil {
		products = []Product{}
	}
	return products, nil
}

func (s *ProductService) Get(ctx context.Context, id int) (*Product, error) {
	p, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product %d not found", id)
		}
		return nil, fmt.Errorf("get: %w", err)
	}
	return p, nil
}

func (s *ProductService) Create(ctx context.Context, p *Product) (*Product, error) {
	id, err := s.repo.Create(ctx, p)
	if err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}
	p.ID = int(id)
	return p, nil
}

func (s *ProductService) Update(ctx context.Context, p *Product) error {
	if err := s.repo.Update(ctx, p); err != nil {
		return fmt.Errorf("update: %w", err)
	}
	return nil
}

func (s *ProductService) Delete(ctx context.Context, id int) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete: %w", err)
	}
	return nil
}

// ----- Handler layer -----

type ProductHandler struct {
	svc *ProductService
}

func NewProductHandler(svc *ProductService) *ProductHandler {
	return &ProductHandler{svc: svc}
}

func (h *ProductHandler) List(w http.ResponseWriter, r *http.Request) {
	products, err := h.svc.List(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, products)
}

func (h *ProductHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	p, err := h.svc.Get(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, p)
}

func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	var p Product
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	created, err := h.svc.Create(r.Context(), &p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusCreated, created)
}

func (h *ProductHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	var p Product
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	p.ID = id
	if err := h.svc.Update(r.Context(), &p); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *ProductHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	if err := h.svc.Delete(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// ----- helpers -----

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start))
	})
}

// ----- server setup -----

func setupServer(dbPath string) (*ProductHandler, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}

	repo := NewSQLiteRepo(db)
	if err := repo.Init(context.Background()); err != nil {
		return nil, fmt.Errorf("init db: %w", err)
	}
	svc := NewProductService(repo)
	h := NewProductHandler(svc)
	return h, nil
}

func main() {
	h, err := setupServer("shop.db")
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /products", h.List)
	mux.HandleFunc("GET /products/{id}", h.Get)
	mux.HandleFunc("POST /products", h.Create)
	mux.HandleFunc("PUT /products/{id}", h.Update)
	mux.HandleFunc("DELETE /products/{id}", h.Delete)

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      loggingMiddleware(corsMiddleware(mux)),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Println("Server on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	srv.Shutdown(ctx)
}
