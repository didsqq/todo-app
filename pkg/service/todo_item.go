package service

import (
	"github.com/didsqq/todo-app"
	"github.com/didsqq/todo-app/pkg/repository"
)

type TodoItemService struct {
	repo     repository.TodoItem
	listRepo repository.TodoList
}

func NewTodoItemService(repo repository.TodoItem, listRepo repository.TodoList) *TodoItemService {
	return &TodoItemService{repo: repo, listRepo: listRepo}
}

func (s *TodoItemService) Create(userId int, listId int, item todo.TodoItem) (int, error) {
	_, err := s.listRepo.GetById(userId, listId)
	if err != nil {
		return 0, err
	}
	return s.repo.Create(userId, listId, item)
}

func (s *TodoItemService) GetAll(userId int, listId int) ([]todo.TodoItem, error) {
	_, err := s.listRepo.GetById(userId, listId)
	if err != nil {
		return nil, err
	}
	return s.repo.GetAll(listId)
}

func (s *TodoItemService) GetById(userId int, listId int, itemId int) (todo.TodoItem, error) {
	_, err := s.listRepo.GetById(userId, listId)
	if err != nil {
		return todo.TodoItem{}, err
	}
	return s.repo.GetById(listId, itemId)
}

func (s *TodoItemService) Delete(userId int, listId, itemId int) error {
	_, err := s.listRepo.GetById(userId, listId)
	if err != nil {
		return err
	}
	return s.repo.Delete(listId, itemId)
}

func (s *TodoItemService) Update(userId int, listId int, itemId int, input todo.UpdateItemInput) error {
	if err := input.Validate(); err != nil {
		return err
	}
	_, err := s.listRepo.GetById(userId, listId)
	if err != nil {
		return err
	}
	return s.repo.Update(listId, itemId, input)
}
