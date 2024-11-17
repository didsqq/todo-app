package repository

import (
	"errors"
	"fmt"

	"github.com/didsqq/todo-app"
	"github.com/jmoiron/sqlx"
)

type TodoListPostgres struct {
	db *sqlx.DB
}

func NewTodoListService(db *sqlx.DB) *TodoListPostgres {
	return &TodoListPostgres{db: db}
}

func (r *TodoListPostgres) Create(userId int, list todo.TodoList) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var id int
	createListQuery := fmt.Sprintf("INSERT INTO %s (title, description, user_id) VALUES ($1, $2, $3) RETURNING id", todoListsTable)
	row := tx.QueryRow(createListQuery, list.Title, list.Description, userId)
	if err := row.Scan(&id); err != nil {
		tx.Rollback()
		return 0, err
	}

	return id, tx.Commit()
}

func (r *TodoListPostgres) GetAll(userId int) ([]todo.TodoList, error) {
	var lists []todo.TodoList
	getListsQuery := fmt.Sprintf("SELECT * FROM %s WHERE user_id=$1", todoListsTable)
	err := r.db.Select(&lists, getListsQuery, userId)
	return lists, err
}

func (r *TodoListPostgres) GetById(userId int, listId int) (todo.TodoList, error) {
	var list todo.TodoList
	getListByIdQuery := fmt.Sprintf("SELECT * FROM %s WHERE id=$1 AND user_id=$2", todoListsTable)
	err := r.db.Get(&list, getListByIdQuery, listId, userId)
	return list, err
}

func (r *TodoListPostgres) Delete(userId int, listId int) error {
	deleteListQuery := fmt.Sprintf("DELETE FROM %s WHERE id=$1 AND user_id=$2", todoListsTable)
	del, err := r.db.Exec(deleteListQuery, listId, userId)
	if err != nil {
		return err
	}
	result, err := del.RowsAffected() // метод возвращает кол-во строк затронутых sql
	if result == 0 {
		err = errors.New("list does not exist")
	}
	return err
}
