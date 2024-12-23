package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/didsqq/todo-app"
	"github.com/didsqq/todo-app/pkg/handler"
	"github.com/didsqq/todo-app/pkg/repository"
	"github.com/didsqq/todo-app/pkg/service"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter)) //устанавливается формат вывода логов. здесь используется json формат
	if err := initConfig(); err != nil {
		logrus.Fatalf("error initializing configs: %s", err.Error())
	}

	if err := godotenv.Load(); err != nil { // используется библиотека godotenv для загрузки переменных окружения из файла .env
		logrus.Fatalf("error initializing configs: %s", err.Error())
	}

	db, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("db.host"), // читает значения из конфигурационного файла
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
		Password: os.Getenv("DB_PASSWORD"), // читает пароль из переменных окружениях
	})
	if err != nil {
		logrus.Fatalf("error initializing db: %s", err.Error())
	}

	repos := repository.NewRepository(db)    // слой для работы с бд
	services := service.NewService(repos)    // логика, используются методы repository
	handlers := handler.NewHandler(services) // http-обработчики используют методы service

	srv := new(todo.Server) // создается сервер
	go func() {             // запускаем в отдельной горутине, чтобы запуск сервера не блокировал основной поток
		if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil { // запускается на порту из конф файла, возвращаются маршруты для http-обработчиков
			logrus.Fatalf("error occured while running http server: %s", err.Error())
		}
	}()

	logrus.Print("TodoApp Started")

	quit := make(chan os.Signal, 1)                      // создается канал для получения сигналов завершения
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT) // подписываемся на сигналы SIGTERM, SIGINT обычно это ctrl+c
	<-quit                                               // блокирует выполениение, пока не будет получен сигнал завершения

	logrus.Print("TodoApp Shutting Down")

	if err := srv.Shutdown(context.Background()); err != nil {
		logrus.Errorf("error occured on server shutting down: %s", err.Error())
	}

	if err := db.Close(); err != nil {
		logrus.Errorf("error occured on db connection close: %s", err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
