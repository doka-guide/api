// Package controllers - пакет для обработки данных запросов
package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/doka-guide/api/api/models"
	"github.com/doka-guide/api/api/responses"
	"github.com/doka-guide/api/api/utils/formaterror"
	"github.com/gorilla/mux"
)

// CreateUser - Создание пользователя
func (server *Server) CreateUser(w http.ResponseWriter, r *http.Request) {
	// Проверка авторизации
	if !CheckPermission(server.DB, GetUserIDByToken(w, r), "USER-POST") {
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
	}
	user := models.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	user.Prepare()
	err = user.Validate("")
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	userCreated, err := user.SaveUser(server.DB)

	if err != nil {

		formattedError := formaterror.FormatError(err.Error())

		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.RequestURI, userCreated.ID))
	responses.JSON(w, http.StatusCreated, userCreated)
}

// OptionsUsers – Используется для подготовки соединения
func (server *Server) OptionsUsers(w http.ResponseWriter, r *http.Request) {
	// Проверка авторизации
	if !CheckPermission(server.DB, GetUserIDByToken(w, r), "USER-OPTIONS") {
		return
	}

	responses.JSON(w, http.StatusOK, []byte("Запрос OPTIONS обработан"))
}

// GetUsers - all users
func (server *Server) GetUsers(w http.ResponseWriter, r *http.Request) {
	// Проверка авторизации
	if !CheckPermission(server.DB, GetUserIDByToken(w, r), "USER-GET") {
		return
	}

	user := models.User{}

	users, err := user.FindAllUsers(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, users)
}

// GetUser - Получение информации о пользователе
func (server *Server) GetUser(w http.ResponseWriter, r *http.Request) {
	// Проверка авторизации
	if !CheckPermission(server.DB, GetUserIDByToken(w, r), "USER-GET") {
		return
	}

	vars := mux.Vars(r)
	uid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	user := models.User{}
	userGotten, err := user.FindUserByID(server.DB, uid)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	responses.JSON(w, http.StatusOK, userGotten)
}

// UpdateUser - Обновление информации о пользователе
func (server *Server) UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uid, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	user := models.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	// Проверка авторизации
	tokenID := GetUserIDByToken(w, r)
	if !CheckPermission(server.DB, tokenID, "USER-PUT") {
		return
	}
	if tokenID != uid {
		responses.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}

	user.Prepare()
	err = user.Validate("update")
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	updatedUser, err := user.UpdateAUser(server.DB, uid)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	responses.JSON(w, http.StatusOK, updatedUser)
}

// DeleteUser - Удаление пользователя
func (server *Server) DeleteUser(w http.ResponseWriter, r *http.Request) {
	// Проверка авторизации
	tokenID := GetUserIDByToken(w, r)
	if !CheckPermission(server.DB, tokenID, "USER-DELETE") {
		return
	}

	vars := mux.Vars(r)
	user := models.User{}

	uid, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	if tokenID != uid {
		responses.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}

	if tokenID != 0 && tokenID != uid {
		responses.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}
	_, err = user.DeleteAUser(server.DB, uid)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	w.Header().Set("Entity", fmt.Sprintf("%d", uid))
	responses.JSON(w, http.StatusNoContent, "")
}
