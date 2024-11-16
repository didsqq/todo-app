package service

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/didsqq/todo-app"
	"github.com/didsqq/todo-app/pkg/repository"
)

const (
	salt      = "gafd21traeg"  // для усложнения хеширования
	signinKey = "fgafsgt43gsd" // ключ для подписи jwt токена
	tokenTTL  = 12 * time.Hour // время жизни токена(12 часов)
)

type tokenClaims struct {
	jwt.StandardClaims     // базовые стандартные поля jwt(время истечения действия токена(ExpiresAt), время выпуска(IssuedAt))
	UserId             int `json:"user_id"` // id для идентификации пользователя
}

type AuthService struct {
	repo repository.Authorization
}

func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) CreateUser(user todo.User) (int, error) {
	user.Password = s.generatePasswordHash(user.Password)
	return s.repo.CreateUser(user)
}

func (s *AuthService) GenerateToken(username, password string) (string, error) {
	user, err := s.repo.GetUser(username, s.generatePasswordHash(password))
	if err != nil {
		return "", err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(), // когда истечет токен(через 12 часов)
			IssuedAt:  time.Now().Unix(),               // когда был токен создан
		},
		user.Id,
	})

	return token.SignedString([]byte(signinKey))
}

func (s *AuthService) ParseToken(accesToken string) (int, error) {
	token, err := jwt.ParseWithClaims(accesToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) { // 3 аргумент анонимная функция,сначала работает ParseWithClaims после декодирования токена вызывается анонимная функция, interface{} — это пустой интерфейс, который может содержать значение любого типа
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok { // Проверка, соответствует ли метод подписи типу token.Method — это поле структуры jwt.Token
			return nil, errors.New("invalid signinKey method")
		}

		return []byte(signinKey), nil
	}) //ParseWithClaims завершает процесс
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return 0, errors.New("token claims are not of type *tokenClaims")
	}

	return claims.UserId, nil
}

func (s *AuthService) generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
