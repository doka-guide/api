package controllers

import (
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"../responses"
	"github.com/joho/godotenv"
)

// UploadFile – Загрузка файла из формы
func (server *Server) UploadFile(w http.ResponseWriter, r *http.Request) {

	err := godotenv.Load()
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
	}
	maxUploadSize, err := strconv.ParseInt(os.Getenv("UPLOAD_MAX_SIZE"), 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
	}
	uploadPath := os.Getenv("UPLOAD_FOLDER")

	// Проверка размера файла
	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)
	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	// Проверка и валидация параметров передачи файла из формы
	file, _, err := r.FormFile("file-ready-to-upload")
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	defer file.Close()
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	// Проверка типа файла (используется первые 512 байт)
	detectedFileType := http.DetectContentType(fileBytes)
	switch detectedFileType {
	case "application/gzip":
	case "application/msword":
	case "application/pdf":
	case "application/vnd.openxmlformats-officedocument.wordprocessingml.document":
	case "application/x-7z-compressed":
	case "application/x-rar-compressed":
	case "application/x-tar":
	case "application/zip":
	case "image/jpeg":
	case "image/tiff":
		break
	default:
		responses.JSON(w, http.StatusBadRequest, "INVALID_FILE_TYPE")
		return
	}
	b := make([]byte, 12)
	rand.Read(b)
	t := time.Now()
	fileName := fmt.Sprintf("%x–%d-%02d-%02dT%02d–%02d–%02d", b, t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
	fileEndings, err := mime.ExtensionsByType(detectedFileType)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	newPath := filepath.Join(uploadPath, fileName+fileEndings[0])
	successMessage := fmt.Sprintf("FileType: %s, File: %s\n", detectedFileType, newPath)

	// Запись файла на диск
	newFile, err := os.Create(newPath)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	defer newFile.Close() // Закрывает файл после записи
	if _, err := newFile.Write(fileBytes); err != nil || newFile.Close() != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, successMessage)
}
