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

// CreateSubscription – Создание записи о новой отправленной подписке
func (server *Server) CreateSubscription(w http.ResponseWriter, r *http.Request) {
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

	// Проверка авторизации
	uid := GetUserIDByToken(w, r)
	if uid == 0 {
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

	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.URL.Path, subscription.ID))
	responses.JSON(w, http.StatusCreated, subscription)
}

// OptionsSubscriptions – Для предварительной загрузки (prefetch)
func (server *Server) OptionsSubscriptions(w http.ResponseWriter, r *http.Request) {
	responses.JSON(w, http.StatusOK, []byte("Запрос OPTIONS обработан"))
}

// GetSubscriptions – Вывод всех форм
func (server *Server) GetSubscriptions(w http.ResponseWriter, r *http.Request) {
	if GetUserIDByToken(w, r) == 0 {
		return
	}

	form := models.Subscription{}
	forms, err := form.FindAllSubscriptions(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, forms)
}

// GetSubscription – Вывод подписки по ID
func (server *Server) GetSubscription(w http.ResponseWriter, r *http.Request) {
	if GetUserIDByToken(w, r) == 0 {
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
	responses.JSON(w, http.StatusOK, formReceived)
}

// UpdateSubscription – Обновление информации в подписке
func (server *Server) UpdateSubscription(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	// Валидация информации о подписке
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	// Проверка авторизации
	uid := GetUserIDByToken(w, r)
	if uid == 0 {
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
	responses.JSON(w, http.StatusOK, formUpdated)
}

// DeleteSubscription – Удаляет данные подписки из базы данных
func (server *Server) DeleteSubscription(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	// Валидация подписки
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	// Проверка авторизации
	uid := GetUserIDByToken(w, r)
	if uid == 0 {
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

// GetSubscriptionFormsWithHash – Вывод адресов электронной почты и настроек с указанием хэша
func (server *Server) GetSubscriptionFormsWithHash(w http.ResponseWriter, r *http.Request) {
	if GetUserIDByToken(w, r) == 0 {
		return
	}

	vars := mux.Vars(r)
	start := vars["start"]
	end := vars["end"]

	form := models.Form{}
	report := form.SubscriptionFormsWithHash(server.DB, start, end)
	responses.JSON(w, http.StatusOK, report)
}
