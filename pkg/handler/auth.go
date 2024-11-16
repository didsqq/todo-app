package handler

import (
	"net/http"

	"github.com/didsqq/todo-app"
	"github.com/gin-gonic/gin"
)

func (h *Handler) signUp(c *gin.Context) {
	var input todo.User

	if err := c.BindJSON(&input); err != nil { // gin автоматически связывает json данные из тела запроса с переменной input
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.Authorization.CreateUser(input) // метод для создания пользователя
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

type singInInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *Handler) signIn(c *gin.Context) {
	var input singInInput

	if err := c.BindJSON(&input); err != nil { // gin автоматически связывает json данные из тела запроса с переменной input
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	//logrus.Infof("username:%s pass:%s", input.Username, input.Password)
	token, err := h.services.Authorization.GenerateToken(input.Username, input.Password) // метод для создания пользователя
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"token": token,
	})
}
