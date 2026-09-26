// Package controllers - пакет для обработки данных запросов
package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"

	"github.com/doka-guide/api/api/models"
	"github.com/doka-guide/api/api/responses"
	"github.com/doka-guide/api/api/utils/formaterror"
	"github.com/doka-guide/api/api/utils/mail"
	"github.com/gorilla/mux"
)

type SubscriptionInfo struct {
	ID    uint64 `json:"id"`
	Email string `json:"email"`
	Data  string `json:"data"`
}

// CreateSubscription – Создание записи о новой отправленной подписке
func (server *Server) CreateSubscription(w http.ResponseWriter, r *http.Request) {
	// Проверка авторизации
	uid := GetUserIDByToken(w, r)
	if !CheckPermission(server.DB, uid, "SUBSCRIPTION-POST") {
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	subForm := models.Subscription{}
	err = json.Unmarshal(body, &subForm)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	subForm.Prepare()
	err = subForm.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	if uid != subForm.AuthorID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}
	subscription, err := subForm.SaveSubscription(server.DB)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	profileLinkForm := models.ProfileLink{}
	profileLinkForm.Prepare()
	profileLinkForm.AuthorID = subForm.AuthorID
	profileLinkForm.ProfileID = subForm.ID
	err = profileLinkForm.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	if uid != profileLinkForm.AuthorID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}
	_, err = profileLinkForm.SaveProfileLink(server.DB)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}

	hiImages := os.Getenv("MAIL_IMAGES_HI_HTML")
	imagesRegex := regexp.MustCompile(`\.\/images`)

	hiTxt, err := ioutil.ReadFile(os.Getenv("MAIL_BODY_HI_TEXT"))
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	hiHTML, err := ioutil.ReadFile(os.Getenv("MAIL_BODY_HI_HTML"))
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	varRegex := regexp.MustCompile(`{{ hash }}`)

	mail.SendMail(
		"Дорогой участник",
		subForm.Email,
		os.Getenv("MAIL_TITLE"),
		string(varRegex.ReplaceAllString(string(hiTxt), profileLinkForm.Hash)),
		string(imagesRegex.ReplaceAllString(
			string(varRegex.ReplaceAllString(string(hiHTML), profileLinkForm.Hash)),
			string(hiImages),
		)),
		false,
	)

	s := SubscriptionInfo{}
	s.ID = subscription.ID
	s.Email = subscription.Email
	s.Data = subscription.Data
	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.URL.Path, s.ID))
	responses.JSON(w, http.StatusCreated, s)
}

// OptionsSubscriptions – Для предварительной загрузки (prefetch)
func (server *Server) OptionsSubscriptions(w http.ResponseWriter, r *http.Request) {
	responses.JSON(w, http.StatusOK, []byte("Запрос OPTIONS обработан"))
}

// GetSubscriptions – Вывод всех форм
func (server *Server) GetSubscriptions(w http.ResponseWriter, r *http.Request) {
	// Проверка авторизации
	if !CheckPermission(server.DB, GetUserIDByToken(w, r), "SUBSCRIPTION-GET") {
		return
	}

	form := models.Subscription{}
	forms, err := form.FindAllSubscriptions(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	count := len(*forms)
	responseForms := []SubscriptionInfo{}
	for i := 0; i < count; i++ {
		formInfo := SubscriptionInfo{}
		formInfo.ID = (*forms)[i].ID
		formInfo.Email = (*forms)[i].Email
		formInfo.Data = (*forms)[i].Data
		responseForms = append(responseForms, formInfo)
	}
	responses.JSON(w, http.StatusOK, forms)
}

// GetSubscription – Вывод подписки по ID
func (server *Server) GetSubscription(w http.ResponseWriter, r *http.Request) {
	// Проверка авторизации
	if !CheckPermission(server.DB, GetUserIDByToken(w, r), "SUBSCRIPTION-GET") {
		return
	}

	vars := mux.Vars(r)
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	form := models.Subscription{}

	formReceived, err := form.FindSubscriptionByID(server.DB, pid)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	s := SubscriptionInfo{}
	s.ID = formReceived.ID
	s.Email = formReceived.Email
	s.Data = formReceived.Data
	responses.JSON(w, http.StatusOK, s)
}

// UpdateSubscription – Обновление информации в подписке
func (server *Server) UpdateSubscription(w http.ResponseWriter, r *http.Request) {
	// Проверка авторизации
	uid := GetUserIDByToken(w, r)
	if !CheckPermission(server.DB, uid, "SUBSCRIPTION-PUT") {
		return
	}

	vars := mux.Vars(r)

	// Валидация информации о подписке
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	// Проверка существования подписки
	form := models.Subscription{}
	err = server.DB.Debug().Model(models.Subscription{}).Where("id = ?", pid).Take(&form).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("Subscription not found"))
		return
	}

	// Если пользователь захочет обновить форму, которая отправлена не от него
	if uid != form.AuthorID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	// Чтение данных подписки
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	// Начало обработки данных подписки
	formUpdate := models.Subscription{}
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
	formUpdated, err := formUpdate.UpdateASubscription(server.DB)

	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}

	s := SubscriptionInfo{}
	s.ID = formUpdated.ID
	s.Email = formUpdated.Email
	s.Data = formUpdated.Data
	responses.JSON(w, http.StatusOK, s)
}

// DeleteSubscription – Удаляет данные подписки из базы данных
func (server *Server) DeleteSubscription(w http.ResponseWriter, r *http.Request) {
	// Проверка авторизации
	uid := GetUserIDByToken(w, r)
	if !CheckPermission(server.DB, uid, "SUBSCRIPTION-DELETE") {
		return
	}

	vars := mux.Vars(r)

	// Валидация подписки
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	// Проверка наличия подписки
	form := models.Subscription{}
	err = server.DB.Debug().Model(models.Subscription{}).Where("id = ?", pid).Take(&form).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("Unauthorized"))
		return
	}

	// Проверка принадлежности подписки пользователю
	if uid != form.AuthorID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	_, err = form.DeleteASubscription(server.DB, pid, uid)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	w.Header().Set("Entity", fmt.Sprintf("%d", pid))

	responses.JSON(w, http.StatusNoContent, "")
}

// GetSubscriptionFormsWithHash – Вывод адресов электронной почты и настроек с указанием хеша
func (server *Server) GetSubscriptionFormsWithHash(w http.ResponseWriter, r *http.Request) {
	// Проверка авторизации
	if !CheckPermission(server.DB, GetUserIDByToken(w, r), "SUBSCRIPTION-GET") {
		return
	}

	vars := mux.Vars(r)
	start := vars["start"]
	end := vars["end"]

	form := models.Form{}
	report := form.SubscriptionFormsWithHash(server.DB, start, end)
	responses.JSON(w, http.StatusOK, report)
}
