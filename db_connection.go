package main

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Todo struct {
	ID          int
	UID         string
	Task        string
	Status      bool
	CreatedAt   time.Time
	CompletedAt *time.Time
}

type Store struct {
	pool *pgxpool.Pool
}

// Initialize pgx connection pool
func NewStore(ctx context.Context, connStr string) (*Store, error) {
	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		return nil, fmt.Errorf("pgxpool.New: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping failed: %w", err)
	}
	return &Store{pool: pool}, nil
}

func (s *Store) Close() {
	s.pool.Close()
}

func (s *Store) AddTodo(ctx context.Context, task string) error {
	const q = `
	INSERT INTO todos (task, status)
	VALUES ($1, false)
	`
	_, err := s.pool.Exec(ctx, q, task)
	return err
}

func (s *Store) ListTodos(ctx context.Context, doneFilter *bool) ([]Todo, error) {
	base := `SELECT id, uid, task, status, created_at, completed_at FROM todos`
	var (
		rows pgx.Rows
		err  error
	)

	if doneFilter != nil {
		q := base + ` WHERE status = $1 ORDER BY created_at`
		rows, err = s.pool.Query(ctx, q, *doneFilter)
	} else {
		q := base + ` ORDER BY created_at`
		rows, err = s.pool.Query(ctx, q)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []Todo
	for rows.Next() {
		var t Todo
		if err := rows.Scan(&t.ID, &t.UID, &t.Task, &t.Status, &t.CreatedAt, &t.CompletedAt); err != nil {
			return nil, err
		}
		todos = append(todos, t)
	}
	return todos, rows.Err()
}

func (s *Store) MarkDone(ctx context.Context, uid string) error {
	const q = `
	UPDATE todos
	SET status = true, completed_at = now()
	WHERE uid = $1 AND status = false
	`
	ct, err := s.pool.Exec(ctx, q, uid)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return fmt.Errorf("no matching todo or already done")
	}
	return nil
}

func (s *Store) MarkAllDone(ctx context.Context) (int64, error) {
	const q = `
    UPDATE todos
    SET status = true, completed_at = now()
    WHERE status = false
    `
	ct, err := s.pool.Exec(ctx, q)
	if err != nil {
		return 0, err
	}
	return ct.RowsAffected(), nil
}

func (s *Store) DeleteTodo(ctx context.Context, uid string) error {
	const q = `DELETE FROM todos WHERE uid = $1`
	_, err := s.pool.Exec(ctx, q, uid)
	return err
}

func (s *Store) DeleteAllTodos(ctx context.Context) (int64, error) {
	const q = `DELETE FROM todos`
	ct, err := s.pool.Exec(ctx, q)
	if err != nil {
		return 0, err
	}
	return ct.RowsAffected(), nil
}

// UnmarkDone sets status=false for a specific todo by UID.
func (s *Store) UnmarkDone(ctx context.Context, uid string) error {
	const q = `
UPDATE todos
SET status = false, completed_at = NULL
WHERE uid = $1 AND status = true
`
	ct, err := s.pool.Exec(ctx, q, uid)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return fmt.Errorf("no matching todo or already undone")
	}
	return nil
}

// UnmarkAllDone sets status=false and resets completed_at for all done todos.
// Returns the number of rows affected.
func (s *Store) UnmarkAllDone(ctx context.Context) (int64, error) {
	const q = `
UPDATE todos
SET status = false, completed_at = NULL
WHERE status = true
`
	ct, err := s.pool.Exec(ctx, q)
	if err != nil {
		return 0, err
	}
	return ct.RowsAffected(), nil
}
