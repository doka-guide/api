package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"../auth"
	"../models"
	"../responses"
	"../utils/formaterror"
	"../utils/mail"
	"github.com/gorilla/mux"
)

var users = []models.User{
	models.User{
		Nickname: "Igor Korovchenko",
		Email:    "igsekor@gmail.com",
		Password: "MdyVHvqJwU74SrL4cX23QpPOeltWTCjXLuoztDf9",
	},
}

// CreateForm – Создание записи о новой отправленной форме
func (server *Server) CreateForm(w http.ResponseWriter, r *http.Request) {

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
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	if uid != form.AuthorID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
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

	subject := "Новая форма была заполнена на сайте"
	emailBody := ""
	switch form.Type {
	case "feedback":
		submittedForm := models.FormFeedback{}
		err = json.Unmarshal([]byte(form.Data), &submittedForm)
		if err != nil {
			responses.ERROR(w, http.StatusUnprocessableEntity, err)
			return
		}
		emailBody = "Тип: Отзыв о статье\n" + submittedForm.ToString()
	}

	for i := range users {
		err = mail.SendMail(users[i].Nickname, users[i].Email, subject, emailBody)
		if err != nil {
			responses.ERROR(w, http.StatusUnprocessableEntity, err)
			return
		}
	}
}

// OptionsForms – Для предзагрузки (prefetch)
func (server *Server) OptionsForms(w http.ResponseWriter, r *http.Request) {
	responses.JSON(w, http.StatusOK, []byte("Request with options has been processed"))
}

// GetForms – Вывод всех форм
func (server *Server) GetForms(w http.ResponseWriter, r *http.Request) {
	_, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
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

	_, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
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

	vars := mux.Vars(r)

	// Check if the form id is valid
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
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

	formUpdate.ID = form.ID //this is important to tell the model the form id to update, the other update field are set above

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

	vars := mux.Vars(r)

	// Валидация формы
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
