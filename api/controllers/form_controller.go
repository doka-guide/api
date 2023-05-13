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

// CreateForm – Создание записи о новой отправленной форме
func (server *Server) CreateForm(w http.ResponseWriter, r *http.Request) {
	// Проверка авторизации
	if CheckPermission(GetUserIDByToken(w, r), "FORM-POST") {
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	form := models.Form{}
	err = json.Unmarshal(body, &form)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	form.Prepare()
	err = form.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	formCreated, err := form.SaveForm(server.DB)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}

	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.URL.Path, formCreated.ID))
	responses.JSON(w, http.StatusCreated, formCreated)

	switch form.Type {
	case "feedback":
		submittedForm := models.FormFeedback{}
		err = json.Unmarshal([]byte(form.Data), &submittedForm)
		if err != nil {
			responses.ERROR(w, http.StatusUnprocessableEntity, err)
			return
		}
	}
}

// OptionsForms – Для предварительной загрузки (prefetch)
func (server *Server) OptionsForms(w http.ResponseWriter, r *http.Request) {
	responses.JSON(w, http.StatusOK, []byte("Запрос OPTIONS обработан"))
}

// GetForms – Вывод всех форм
func (server *Server) GetForms(w http.ResponseWriter, r *http.Request) {
	// Проверка авторизации
	if CheckPermission(GetUserIDByToken(w, r), "FORM-GET") {
		return
	}

	form := models.Form{}
	forms, err := form.FindAllForms(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, forms)
}

// GetForm – Вывод формы по ID
func (server *Server) GetForm(w http.ResponseWriter, r *http.Request) {
	// Проверка авторизации
	if CheckPermission(GetUserIDByToken(w, r), "FORM-GET") {
		return
	}

	vars := mux.Vars(r)
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	form := models.Form{}

	formReceived, err := form.FindFormByID(server.DB, pid)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, formReceived)
}

// UpdateForm – Обновление информации в форме
func (server *Server) UpdateForm(w http.ResponseWriter, r *http.Request) {
	// Проверка авторизации
	uid := GetUserIDByToken(w, r)
	if CheckPermission(uid, "FORM-PUT") {
		return
	}

	vars := mux.Vars(r)

	// Валидация полей формы
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	// Проверка существования формы
	form := models.Form{}
	err = server.DB.Debug().Model(models.Form{}).Where("id = ?", pid).Take(&form).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("Form not found"))
		return
	}

	// Если пользователь захочет обновить форму, которая отправлена не от него
	if uid != form.AuthorID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	// Чтение данных формы
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	// Начало обработки данных формы
	formUpdate := models.Form{}
	err = json.Unmarshal(body, &formUpdate)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	// Проверка авторизации
	if uid != formUpdate.AuthorID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	formUpdate.Prepare()
	err = formUpdate.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	formUpdate.ID = form.ID
	formUpdated, err := formUpdate.UpdateAForm(server.DB)

	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	responses.JSON(w, http.StatusOK, formUpdated)
}

// DeleteForm – Удаляет данные формы из базы данных
func (server *Server) DeleteForm(w http.ResponseWriter, r *http.Request) {
	// Проверка авторизации
	uid := GetUserIDByToken(w, r)
	if CheckPermission(uid, "FORM-DELETE") {
		return
	}

	vars := mux.Vars(r)

	// Валидация формы
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	// Проверка наличия формы
	form := models.Form{}
	err = server.DB.Debug().Model(models.Form{}).Where("id = ?", pid).Take(&form).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("Unauthorized"))
		return
	}

	// Проверка принадлежности формы пользователю
	if uid != form.AuthorID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	_, err = form.DeleteAForm(server.DB, pid, uid)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	w.Header().Set("Entity", fmt.Sprintf("%d", pid))
	responses.JSON(w, http.StatusNoContent, "")
}

// GetFeedbackForms – Вывод информации о заполненных формах обратной связи за период
func (server *Server) GetFeedbackForms(w http.ResponseWriter, r *http.Request) {
	// Проверка авторизации
	if CheckPermission(GetUserIDByToken(w, r), "FORM-GET") {
		return
	}

	vars := mux.Vars(r)
	start := vars["start"]
	end := vars["end"]

	form := models.Form{}
	report := form.FeedbackFormsGroupedByData(server.DB, start, end)
	responses.JSON(w, http.StatusOK, report)
}

// GetQuestionForms – Вывод информации о заполненных формах обратной связи за период
func (server *Server) GetQuestionForms(w http.ResponseWriter, r *http.Request) {
	// Проверка авторизации
	if CheckPermission(GetUserIDByToken(w, r), "FORM-GET") {
		return
	}

	vars := mux.Vars(r)
	start := vars["start"]
	end := vars["end"]

	form := models.Form{}
	report := form.QuestionForms(server.DB, start, end)
	responses.JSON(w, http.StatusOK, report)
}
