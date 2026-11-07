package model

import "time"

// Status represents the completion state of a Todo.
type Status string

const (
	StatusPending    Status = "pending"
	StatusInProgress Status = "in_progress"
	StatusDone       Status = "done"
)

// Todo is the core domain model.
type Todo struct {
	ID          int64     `json:"id"           db:"id"`
	Title       string    `json:"title"        db:"title"`
	Description string    `json:"description"  db:"description"`
	Status      Status    `json:"status"       db:"status"`
	CreatedAt   time.Time `json:"created_at"   db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"   db:"updated_at"`
}

// CreateTodoRequest is the payload for creating a new Todo.
type CreateTodoRequest struct {
	Title       string `json:"title"       validate:"required,min=1,max=255"`
	Description string `json:"description" validate:"max=1000"`
	Status      Status `json:"status"      validate:"omitempty,oneof=pending in_progress done"`
}

// UpdateTodoRequest is the payload for updating an existing Todo.
type UpdateTodoRequest struct {
	Title       *string `json:"title"       validate:"omitempty,min=1,max=255"`
	Description *string `json:"description" validate:"omitempty,max=1000"`
	Status      *Status `json:"status"      validate:"omitempty,oneof=pending in_progress done"`
}
