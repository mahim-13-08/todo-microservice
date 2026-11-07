package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/yourusername/todo-microservice/internal/model"
	_ "github.com/mattn/go-sqlite3"
)

// ErrNotFound is returned when a Todo cannot be found.
var ErrNotFound = errors.New("todo not found")

// TodoRepository defines the interface for data access.
type TodoRepository interface {
	Create(ctx context.Context, req *model.CreateTodoRequest) (*model.Todo, error)
	GetAll(ctx context.Context) ([]*model.Todo, error)
	GetByID(ctx context.Context, id int64) (*model.Todo, error)
	Update(ctx context.Context, id int64, req *model.UpdateTodoRequest) (*model.Todo, error)
	Delete(ctx context.Context, id int64) error
}

type sqliteRepository struct {
	db *sql.DB
}

// NewSQLiteRepository creates a new SQLite-backed TodoRepository.
func NewSQLiteRepository(dsn string) (TodoRepository, error) {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping db: %w", err)
	}
	repo := &sqliteRepository{db: db}
	if err := repo.migrate(); err != nil {
		return nil, fmt.Errorf("migrate: %w", err)
	}
	return repo, nil
}

func (r *sqliteRepository) migrate() error {
	const schema = `
	CREATE TABLE IF NOT EXISTS todos (
		id          INTEGER PRIMARY KEY AUTOINCREMENT,
		title       TEXT    NOT NULL,
		description TEXT    NOT NULL DEFAULT '',
		status      TEXT    NOT NULL DEFAULT 'pending',
		created_at  DATETIME NOT NULL,
		updated_at  DATETIME NOT NULL
	);`
	_, err := r.db.Exec(schema)
	return err
}

func (r *sqliteRepository) Create(ctx context.Context, req *model.CreateTodoRequest) (*model.Todo, error) {
	status := req.Status
	if status == "" {
		status = model.StatusPending
	}
	now := time.Now().UTC()
	res, err := r.db.ExecContext(ctx,
		`INSERT INTO todos (title, description, status, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?)`,
		req.Title, req.Description, status, now, now,
	)
	if err != nil {
		return nil, fmt.Errorf("insert todo: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("last insert id: %w", err)
	}
	return r.GetByID(ctx, id)
}

func (r *sqliteRepository) GetAll(ctx context.Context) ([]*model.Todo, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, title, description, status, created_at, updated_at FROM todos ORDER BY id DESC`)
	if err != nil {
		return nil, fmt.Errorf("query todos: %w", err)
	}
	defer rows.Close()

	var todos []*model.Todo
	for rows.Next() {
		t := &model.Todo{}
		if err := rows.Scan(&t.ID, &t.Title, &t.Description, &t.Status, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan todo: %w", err)
		}
		todos = append(todos, t)
	}
	return todos, rows.Err()
}

func (r *sqliteRepository) GetByID(ctx context.Context, id int64) (*model.Todo, error) {
	t := &model.Todo{}
	err := r.db.QueryRowContext(ctx,
		`SELECT id, title, description, status, created_at, updated_at FROM todos WHERE id = ?`, id,
	).Scan(&t.ID, &t.Title, &t.Description, &t.Status, &t.CreatedAt, &t.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("query todo by id: %w", err)
	}
	return t, nil
}

func (r *sqliteRepository) Update(ctx context.Context, id int64, req *model.UpdateTodoRequest) (*model.Todo, error) {
	existing, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Title != nil {
		existing.Title = *req.Title
	}
	if req.Description != nil {
		existing.Description = *req.Description
	}
	if req.Status != nil {
		existing.Status = *req.Status
	}
	existing.UpdatedAt = time.Now().UTC()

	_, err = r.db.ExecContext(ctx,
		`UPDATE todos SET title=?, description=?, status=?, updated_at=? WHERE id=?`,
		existing.Title, existing.Description, existing.Status, existing.UpdatedAt, id,
	)
	if err != nil {
		return nil, fmt.Errorf("update todo: %w", err)
	}
	return existing, nil
}

func (r *sqliteRepository) Delete(ctx context.Context, id int64) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM todos WHERE id=?`, id)
	if err != nil {
		return fmt.Errorf("delete todo: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}
