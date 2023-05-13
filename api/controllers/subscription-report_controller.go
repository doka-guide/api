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

// CreateSubscriptionReport – Создание отчёта о загрузке ссылки
func (server *Server) CreateSubscriptionReport(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	report := models.SubscriptionReport{}
	err = json.Unmarshal(body, &report)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	report.Prepare()
	err = report.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	reportCreated, err := report.SaveSubscriptionReport(server.DB)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}

	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.URL.Path, reportCreated.ID))
	responses.JSON(w, http.StatusCreated, reportCreated)
}

// OptionsSubscriptionReports – Для предварительной загрузки (prefetch)
func (server *Server) OptionsSubscriptionReports(w http.ResponseWriter, r *http.Request) {
	responses.JSON(w, http.StatusOK, []byte("Запрос OPTIONS обработан"))
}

// GetSubscriptionReports – Вывод всех отчёта о загрузке ссылок
func (server *Server) GetSubscriptionReports(w http.ResponseWriter, r *http.Request) {
	// Проверка авторизации
	if CheckPermission(server.DB, GetUserIDByToken(w, r), "SUBSCRIPTION-REPORT-GET") {
		return
	}

	report := models.SubscriptionReport{}
	reports, err := report.FindAllSubscriptionReports(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, reports)
}

// GetSubscriptionReport – Вывод отчёта о загрузке ссылки по Hash
func (server *Server) GetSubscriptionReport(w http.ResponseWriter, r *http.Request) {
	// Проверка авторизации
	if CheckPermission(server.DB, GetUserIDByToken(w, r), "SUBSCRIPTION-REPORT-GET") {
		return
	}

	vars := mux.Vars(r)
	path := vars["id"]
	report := models.SubscriptionReport{}

	reportReceived, err := report.FindSubscriptionReportByPath(server.DB, path)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, reportReceived)
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
	uid := GetUserIDByToken(w, r)
	if CheckPermission(server.DB, uid, "SUBSCRIPTION-REPORT-DELETE") {
		return
	}

	// Проверка наличия подписки
	report := models.SubscriptionReport{}
	err = server.DB.Debug().Model(models.SubscriptionReport{}).Where("id = ?", pid).Take(&report).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("Unauthorized"))
		return
	}

	// Проверка принадлежности подписки пользователю
	if uid != report.Profile.AuthorID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	_, err = report.DeleteASubscriptionReport(server.DB, pid, uid)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	w.Header().Set("Entity", fmt.Sprintf("%d", pid))
	responses.JSON(w, http.StatusNoContent, "")
}
