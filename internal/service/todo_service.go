package service

import (
	"context"
	"fmt"

	"github.com/yourusername/todo-microservice/internal/model"
	"github.com/yourusername/todo-microservice/internal/repository"
)

// TodoService contains the business logic for Todos.
type TodoService interface {
	CreateTodo(ctx context.Context, req *model.CreateTodoRequest) (*model.Todo, error)
	GetAllTodos(ctx context.Context) ([]*model.Todo, error)
	GetTodoByID(ctx context.Context, id int64) (*model.Todo, error)
	UpdateTodo(ctx context.Context, id int64, req *model.UpdateTodoRequest) (*model.Todo, error)
	DeleteTodo(ctx context.Context, id int64) error
}

type todoService struct {
	repo repository.TodoRepository
}

// NewTodoService creates a new TodoService backed by the given repository.
func NewTodoService(repo repository.TodoRepository) TodoService {
	return &todoService{repo: repo}
}

func (s *todoService) CreateTodo(ctx context.Context, req *model.CreateTodoRequest) (*model.Todo, error) {
	todo, err := s.repo.Create(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("service create todo: %w", err)
	}
	return todo, nil
}

func (s *todoService) GetAllTodos(ctx context.Context) ([]*model.Todo, error) {
	todos, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("service get all todos: %w", err)
	}
	return todos, nil
}

func (s *todoService) GetTodoByID(ctx context.Context, id int64) (*model.Todo, error) {
	todo, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("service get todo by id: %w", err)
	}
	return todo, nil
}

func (s *todoService) UpdateTodo(ctx context.Context, id int64, req *model.UpdateTodoRequest) (*model.Todo, error) {
	todo, err := s.repo.Update(ctx, id, req)
	if err != nil {
		return nil, fmt.Errorf("service update todo: %w", err)
	}
	return todo, nil
}

func (s *todoService) DeleteTodo(ctx context.Context, id int64) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("service delete todo: %w", err)
	}
	return nil
}
