package repository

import (
	"errors"
	"fmt"
	"strings"

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

func (r *TodoItemPostgres) Delete(listId int, itemId int) error {
	deleteItemQuery := fmt.Sprintf("DELETE FROM %s WHERE list_id=$1 AND id=$2", todoItemsTable)
	row, err := r.db.Exec(deleteItemQuery, listId, itemId)
	if err != nil {
		return err
	}
	result, err := row.RowsAffected()
	if result == 0 {
		err = errors.New("item does not exist")
	}
	return err
}

func (r *TodoItemPostgres) Update(listId int, itemId int, input todo.UpdateItemInput) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if input.Title != nil {
		setValues = append(setValues, fmt.Sprintf("title=$%d", argId))
		args = append(args, *input.Title)
		argId++
	}

	if input.Description != nil {
		setValues = append(setValues, fmt.Sprintf("description=$%d", argId))
		args = append(args, *input.Description)
		argId++
	}

	if input.Done != nil {
		setValues = append(setValues, fmt.Sprintf("done=$%d", argId))
		args = append(args, *input.Done)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf("UPDATE %s SET %s WHERE list_id=$%d AND id=$%d", todoItemsTable, setQuery, argId, argId+1)

	args = append(args, listId, itemId)

	update, err := r.db.Exec(query, args...)
	if err != nil {
		return err
	}

	result, err := update.RowsAffected()
	if result == 0 {
		err = errors.New("item does not exist")
	}

	return err
}
