// Package controllers - пакет для обработки данных запросов
package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/doka-guide/api/api/auth"
	"github.com/doka-guide/api/api/models"
	"github.com/doka-guide/api/api/responses"
	"github.com/doka-guide/api/api/utils/formaterror"
	"github.com/gorilla/mux"
)

// CreateProfileLink – Создание ссылки
func (server *Server) CreateProfileLink(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	link := models.ProfileLink{}
	err = json.Unmarshal(body, &link)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	link.Prepare()
	err = link.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	if uid != link.AuthorID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}
	linkCreated, err := link.SaveProfileLink(server.DB)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}

	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.URL.Path, linkCreated.ID))
	responses.JSON(w, http.StatusCreated, linkCreated)
}

// OptionsProfileLinks – Для предварительной загрузки (prefetch)
func (server *Server) OptionsProfileLinks(w http.ResponseWriter, r *http.Request) {
	responses.JSON(w, http.StatusOK, []byte("Request with options has been processed"))
}

// GetProfileLinks – Вывод всех ссылок
func (server *Server) GetProfileLinks(w http.ResponseWriter, r *http.Request) {
	_, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	link := models.ProfileLink{}
	links, err := link.FindAllProfileLinks(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, links)
}

// GetProfileLink – Вывод ссылки по Hash
func (server *Server) GetProfileLink(w http.ResponseWriter, r *http.Request) {
	_, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	vars := mux.Vars(r)
	hash := vars["id"]
	link := models.ProfileLink{}

	linkReceived, err := link.FindProfileLinkByHash(server.DB, hash)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, linkReceived)
}

// DeleteProfileLink – Удаляет данные о ссылке из базы данных
func (server *Server) DeleteProfileLink(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	// Валидация подписки
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	// Проверка авторизации
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	// Проверка наличия подписки
	link := models.ProfileLink{}
	err = server.DB.Debug().Model(models.ProfileLink{}).Where("id = ?", pid).Take(&link).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("Unauthorized"))
		return
	}

	// Проверка принадлежности подписки пользователю
	if uid != link.AuthorID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	_, err = link.DeleteAProfileLink(server.DB, pid, uid)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	w.Header().Set("Entity", fmt.Sprintf("%d", pid))
	responses.JSON(w, http.StatusNoContent, "")
}
