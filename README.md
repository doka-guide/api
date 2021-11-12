# Бэкенд Доки

Дока — это добрая энциклопедия для веб-разработчиков. Наша цель — сделать документацию по веб-разработке практичной, понятной и не унылой.

Этот репозиторий содержит микросервис (API) для работы с формами на сайте Доки. Данные сохраняются в PostgreSQL.

## Как устроен проект

Проект собирается с помощью GitHub Actions или локально. Для сборки и запуска микросервиса используется команда (необходимо предварительно установить компилятор Go на компьютер и настроить окружение):

```bash
go run main.go
```

Артефакт сборки попадает на сервер. Для настройки работы сервера используются данные из файла `.env`, которые записаны в формате ключ-значение. Для работы микросервиса необходимо заполнить все поля. Список полей:

```bash
# Режим работы приложения (прод (PRODUCTION) или отладка (DEBUG))
MODE=

# Настройки приложения
APP_HOST=
APP_NAME=

# Настройка соединения с почтовым сервером
MAIL_TYPE=
MAIL_HOST=
MAIL_USER=
MAIL_PASS=

# Настройки загрузки файлов пользователей через форму
UPLOAD_FOLDER=
UPLOAD_MAX_SIZE=

# Пользователь по умолчанию
USER_NAME=
USER_MAIL=
USER_PASS=

# Доступ к PostgreSQL
API_SECRET=
DB_HOST=
DB_DRIVER=
DB_USER=
DB_PASSWORD=
DB_NAME=
DB_PORT=
```

## Формат запросов и ответов

Для того, чтобы отправить форму или запросить данные из БД необходимо войти с помощью учётных данных пользователя. Только авторизованные пользователи могут работать с формами. При это под пользователем понимается сервисный пользователь. Механизм уникальных пользователей позволяет разделять формы на группы. Для отправки формы на сайт или получения данных необходимо выполнить два шага (вместо `localhost:8080` необходимо использовать адрес и порт, на которых будет работать микросервис):

1. Авторизация

```bash
$ curl -X POST \
  -H "Content-Type: application/json" \
  -d '{"email":"<email>", "password":"<пароль>"}' \
  localhost:8080/login
```

Ответ содержит ключ авторизации, который используется как токен для доступа к микросервису.

2. Работа с API

Получить список всех форм из базы данных

```bash
$ curl -X GET \
  -H "Content-Type: : application/json" \
  -H "Authorization: <ключ авторизации>" \
  -H "Content-Type: application/json" \
  localhost:8080/form
```

Ответ содержит данные форм в формате JSON.

Отправка данных формы:

```bash
$ curl -X POST \
  -H "Accept: application/json" \
  -H "Authorization: <ключ авторизации>" \
  -H "Content-Type: application/json" \
  -d '{<Данные формы>}' \
  localhost:8080/form
```

Перед отправкой данные необходимо преобразовать в формат JSON, сериализовать и подставить вместо `<Данные формы>`.