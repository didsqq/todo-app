package repository

import (
	"fmt"

	"github.com/didsqq/todo-app"
	"github.com/jmoiron/sqlx"
)

type TodoItemPostgres struct {
	db *sqlx.DB
}

func NewTodoItemPostgres(db *sqlx.DB) *TodoItemPostgres {
	return &TodoItemPostgres{db: db}
}

func (r *TodoItemPostgres) Create(userId int, listId int, item todo.TodoItem) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}
	var id int
	createQuery := fmt.Sprintf("INSERT INTO %s (title, description, list_id) VALUES ($1, $2, $3) RETURNING id", todoItemsTable)
	row := tx.QueryRow(createQuery, item.Title, item.Description, listId)

	if err := row.Scan(&id); err != nil {
		tx.Rollback()
		return 0, err
	}

	return id, tx.Commit()
}

func (r *TodoItemPostgres) GetAll(listId int) ([]todo.TodoItem, error) {
	var items []todo.TodoItem
	getAllQuery := fmt.Sprintf("SELECT * FROM %s WHERE list_id=$1", todoItemsTable)
	err := r.db.Select(&items, getAllQuery, listId)
	return items, err
}

func (r *TodoItemPostgres) GetById(listId int, itemId int) (todo.TodoItem, error) {
	var item todo.TodoItem
	getByIdQuery := fmt.Sprintf("SELECT * FROM %s WHERE list_id=$1 AND id=$2", todoItemsTable)
	err := r.db.Get(&item, getByIdQuery, listId, itemId)
	return item, err
}
