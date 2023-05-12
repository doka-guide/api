// Package controllers - пакет для обработки данных запросов
package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"

	// Драйвер для работы с PostgreSQL
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/doka-guide/api/api/models"
)

// Объект Сервер
type Server struct {
	DB     *gorm.DB
	Router *mux.Router
}

// Инициализация сервера
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

// Запуск сервера
func (server *Server) Run(addr string) {
	fmt.Println("Запустился на хосте", addr)
	log.Fatal(http.ListenAndServe(addr, server.Router))
}
