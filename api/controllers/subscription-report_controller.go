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

// CreateSubscriptionReport – Создание отчёта о загрузке ссылки
func (server *Server) CreateSubscriptionReport(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	link := models.SubscriptionReport{}
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
	linkCreated, err := link.SaveSubscriptionReport(server.DB)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}

	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.URL.Path, linkCreated.ID))
	responses.JSON(w, http.StatusCreated, linkCreated)
}

// OptionsSubscriptionReports – Для предварительной загрузки (prefetch)
func (server *Server) OptionsSubscriptionReports(w http.ResponseWriter, r *http.Request) {
	responses.JSON(w, http.StatusOK, []byte("Request with options has been processed"))
}

// GetSubscriptionReports – Вывод всех отчёта о загрузке ссылок
func (server *Server) GetSubscriptionReports(w http.ResponseWriter, r *http.Request) {
	_, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	link := models.SubscriptionReport{}
	links, err := link.FindAllSubscriptionReports(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, links)
}

// GetSubscriptionReport – Вывод отчёта о загрузке ссылки по Hash
func (server *Server) GetSubscriptionReport(w http.ResponseWriter, r *http.Request) {

	_, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	vars := mux.Vars(r)
	path := vars["id"]
	link := models.SubscriptionReport{}

	linkReceived, err := link.FindSubscriptionReportByPath(server.DB, path)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, linkReceived)
}

// DeleteSubscriptionReport – Удаляет данные о отчёта о загрузке ссылке из базы данных
func (server *Server) DeleteSubscriptionReport(w http.ResponseWriter, r *http.Request) {

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
	link := models.SubscriptionReport{}
	err = server.DB.Debug().Model(models.SubscriptionReport{}).Where("id = ?", pid).Take(&link).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("Unauthorized"))
		return
	}

	// Проверка принадлежности подписки пользователю
	if uid != link.AuthorID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	_, err = link.DeleteASubscriptionReport(server.DB, pid, uid)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	w.Header().Set("Entity", fmt.Sprintf("%d", pid))
	responses.JSON(w, http.StatusNoContent, "")
}
