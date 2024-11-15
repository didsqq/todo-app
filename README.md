go get -u github.com/gin-gonic/gin 
настраивает маршруты (routes) для HTTP-запросов с использованием фреймворка Gin

go get -u github.com/spf13/viper
для работы с файлами конфигураций

docker run --name=todo-db -e POSTGRES_PASSWORD='qwerty' -p 5436:5432 -d --rm postgres
    Пояснение по каждому флагу и параметру:
    docker run — запускает новый контейнер.
    --name=todo-db — задает имя контейнера todo-db. Это имя поможет вам легко идентифицировать контейнер при его управлении (например, при остановке или просмотре логов).
    -e POSTGRES_PASSWORD='qwerty' — устанавливает переменную окружения POSTGRES_PASSWORD со значением 'qwerty'. Эта переменная определяет пароль для учетной записи postgres по умолчанию в PostgreSQL. Вы должны всегда использовать сложные пароли, особенно в продакшн-среде, чтобы повысить безопасность.
    -p 5436:5432 — перенаправляет порт 5432 контейнера на локальный порт 5436. Это значит, что вы сможете получить доступ к PostgreSQL на вашем компьютере через localhost:5436.
    5436 — локальный порт, который будет открыт на вашем компьютере.
    5432 — порт по умолчанию для PostgreSQL в контейнере.
    -d — запускает контейнер в фоновом режиме (detached mode), что позволяет продолжать использовать терминал после запуска контейнера.
    --rm — автоматически удаляет контейнер после его остановки, что удобно для тестов и разработки, но не рекомендуется для долгосрочного хранения данных.
    postgres — это имя образа Docker для PostgreSQL, который Docker загрузит и запустит. Если образ ещё не был загружен, Docker сначала скачает его из Docker Hub.

go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
export PATH=$PATH:$HOME/go/bin
для работы с миграциями

docker exec -it <container_name_or_id> bash - вход в контейнер
psql -U postgres - подключение к postgres
\dt - вывод таблицы

migrate create -ext sql -dir ./schema -seq init 
    migrate create: Это команда для создания нового миграционного файла.
    -ext sql: Указывает расширение файлов миграции. В данном случае используется расширение .sql, что означает, что миграционные файлы будут содержать SQL-запросы
    -dir ./schema: Указывает директорию, в которой будут сохраняться миграционные файлы.
    -seq: Это флаг, который указывает, что миграции будут использовать последовательные номера.
    init: Это название миграции, которую вы хотите создать.
migrate -path ./schema -database 'postgres://postgres:qwerty@localhost:5436/postgres?sslmode=disable' up - применить миграцию

go get -u github.com/jmoiron/sqlx - для работы с бд
github.com/lib/pq — это драйвер PostgreSQL для Go, который необходим для взаимодействия с базой данных PostgreSQL через стандартную библиотеку database/sql

go get -u github.com/joho/godotenv - для работы с паролями

go get -u github.com/sirupsen/logrus - работа с логами

go get -u github.com/dgrijalva/jwt-go - работа с jwt