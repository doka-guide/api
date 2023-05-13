// Package controllers - пакет для обработки данных запросов
package controllers

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"

	// Драйвер для работы с PostgreSQL
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/doka-guide/api/api/auth"
	"github.com/doka-guide/api/api/models"
	"github.com/doka-guide/api/api/responses"
)

// GetUserIDByToken — проверка авторизации пользователей
func GetUserIDByToken(w http.ResponseWriter, r *http.Request) uint64 {
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		fmt.Printf("Ошибка авторизации: Пользователя таким токеном нет в списке")
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return 0
	}

	return uid
}

// CheckPermission — проверяет, есть ли права на осуществление запроса
func CheckPermission(db *gorm.DB, id uint64, permName string) bool {
	if id == 0 {
		return false
	} else {
		users := []models.User{}
		db.Debug().Raw("SELECT users.* FROM permissions JOIN group_permissions ON group_permissions.perms_id = permissions.id JOIN grouped_users ON grouped_users.group_id = group_permissions.group_id JOIN users ON users.id = grouped_users.user_id WHERE permissions.name = ? AND users.id = ?", permName, id).Scan(&users)
		if len(users) > 0 {
			fmt.Printf("Пользователь с id = %d имеет права на операцию '%s'", id, permName)
		} else {
			fmt.Printf("Ошибка авторизации: Пользователь с id = %d не имеет права на операцию '%s'", id, permName)
		}
		return len(users) > 0
	}
}

// Server — Объект Сервер
type Server struct {
	DB     *gorm.DB
	Router *mux.Router
}

// Initialize — Инициализация сервера
func (server *Server) Initialize(Dbdriver, DBUser, DBPassword, DBPort, DBHost, DBName string) {
	var err error
	if Dbdriver == "postgres" {
		DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", DBHost, DBPort, DBUser, DBName, DBPassword)
		server.DB, err = gorm.Open(Dbdriver, DBURL)
		if err != nil {
			fmt.Printf("Не могу подсоединиться к базе данных, используя драйвер %s", Dbdriver)
			log.Fatal("Ошибка:", err)
		} else {
			fmt.Printf("База данных %s подключена\n", Dbdriver)
		}
	}

	// Миграция базы данных
	server.DB.Debug().AutoMigrate(&models.User{})
	server.Router = mux.NewRouter()
	server.initializeRoutes()
}

// Run — Запуск сервера
func (server *Server) Run(addr string) {
	fmt.Println("Запустился на хосте", addr)
	log.Fatal(http.ListenAndServe(addr, server.Router))
}
