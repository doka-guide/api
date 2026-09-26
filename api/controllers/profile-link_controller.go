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

type ProfileLinkInfo struct {
	ID        uint64           `json:"id"`
	Hash      string           `json:"hash"`
	ProfileID uint64           `json:"profile_id"`
	Profile   SubscriptionInfo `json:"profile"`
}

// CreateProfileLink – Создание ссылки
func (server *Server) CreateProfileLink(w http.ResponseWriter, r *http.Request) {
	// Проверка авторизации
	if !CheckPermission(server.DB, GetUserIDByToken(w, r), "PROFILE-LINK-POST") {
		return
	}

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

	linkCreated, err := link.SaveProfileLink(server.DB)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}

	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.URL.Path, linkCreated.ID))

	s := ProfileLinkInfo{}
	s.ID = linkCreated.ID
	s.Hash = linkCreated.Hash
	s.ProfileID = linkCreated.ProfileID
	s.Profile = SubscriptionInfo{}
	s.Profile.ID = linkCreated.Profile.ID
	s.Profile.Email = linkCreated.Profile.Email
	s.Profile.Data = linkCreated.Profile.Data
	responses.JSON(w, http.StatusCreated, s)
}

// OptionsProfileLinks – Для предварительной загрузки (prefetch)
func (server *Server) OptionsProfileLinks(w http.ResponseWriter, r *http.Request) {
	responses.JSON(w, http.StatusOK, []byte("Запрос OPTIONS обработан"))
}

// GetProfileLinks – Вывод всех ссылок
func (server *Server) GetProfileLinks(w http.ResponseWriter, r *http.Request) {
	// Проверка авторизации
	if !CheckPermission(server.DB, GetUserIDByToken(w, r), "PROFILE-LINK-GET") {
		return
	}

	link := models.ProfileLink{}
	links, err := link.FindAllProfileLinks(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	count := len(*links)
	responseForms := []ProfileLinkInfo{}
	for i := 0; i < count; i++ {
		formInfo := ProfileLinkInfo{}
		formInfo.ID = (*links)[i].ID
		formInfo.Hash = (*links)[i].Hash
		formInfo.ProfileID = (*links)[i].ProfileID
		formInfo.Profile = SubscriptionInfo{}
		formInfo.Profile.ID = (*links)[i].Profile.ID
		formInfo.Profile.Email = (*links)[i].Profile.Email
		formInfo.Profile.Data = (*links)[i].Profile.Data
		responseForms = append(responseForms, formInfo)
	}
	responses.JSON(w, http.StatusOK, links)
}

// GetProfileLink – Вывод ссылки по Hash
func (server *Server) GetProfileLink(w http.ResponseWriter, r *http.Request) {
	// Проверка авторизации
	if !CheckPermission(server.DB, GetUserIDByToken(w, r), "PROFILE-LINK-GET") {
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

	s := ProfileLinkInfo{}
	s.ID = linkReceived.ID
	s.Hash = linkReceived.Hash
	s.ProfileID = linkReceived.ProfileID
	s.Profile = SubscriptionInfo{}
	s.Profile.ID = linkReceived.Profile.ID
	s.Profile.Email = linkReceived.Profile.Email
	s.Profile.Data = linkReceived.Profile.Data
	responses.JSON(w, http.StatusOK, s)
}

// DeleteProfileLink – Удаляет данные о ссылке из базы данных
func (server *Server) DeleteProfileLink(w http.ResponseWriter, r *http.Request) {
	// Проверка авторизации
	uid := GetUserIDByToken(w, r)
	if !CheckPermission(server.DB, uid, "PROFILE-LINK-DELETE") {
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
