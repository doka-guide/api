package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/postgres" //postgres database driver

	"../models"
)

// Server object
type Server struct {
	DB     *gorm.DB
	Router *mux.Router
}

// Initialize Server
func (server *Server) Initialize(Dbdriver, DbUser, DbPassword, DbPort, DbHost, DbName string) {

	var err error
	if Dbdriver == "postgres" {
		DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", DbHost, DbPort, DbUser, DbName, DbPassword)
		server.DB, err = gorm.Open(Dbdriver, DBURL)
		if err != nil {
			fmt.Printf("Не могу подсоединиться к базе данных, используя драйвер %s", Dbdriver)
			log.Fatal("Ошибка:", err)
		} else {
			fmt.Printf("База данных %s подключена", Dbdriver)
		}
	}

	server.DB.Debug().AutoMigrate(&models.User{}) //database migration

	server.Router = mux.NewRouter()

	server.initializeRoutes()
}

// Run Server
func (server *Server) Run(addr string) {
	fmt.Println("Запустился на хосте %s", addr)
	log.Fatal(http.ListenAndServe(addr, server.Router))
}