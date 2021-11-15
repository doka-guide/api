package controllers

import (
	"net/http"
	"os"

	"github.com/doka-guide/api/api/responses"
)

// API точка входа
func (server *Server) Home(w http.ResponseWriter, r *http.Request) {
	responses.JSON(w, http.StatusOK, os.Getenv("APP_NAME"))
}

// OptionsHome – Используется для подготовки соединения
func (server *Server) OptionsHome(w http.ResponseWriter, r *http.Request) {
	responses.JSON(w, http.StatusOK, []byte("Предварительный запрос был обработан"))
}
