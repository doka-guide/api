// Package formaterror - пакет для проверки формата данных
package formaterror

import (
	"errors"
	"strings"
)

func FormatError(err string) error {
	if strings.Contains(err, "nickname") {
		return errors.New("Пользователь с таким именем уже существует")
	}
	if strings.Contains(err, "email") {
		return errors.New("Пользователь с такой электронной почтой уже существует")
	}
	if strings.Contains(err, "hashedPassword") {
		return errors.New("Некорректный пароль")
	}
	return errors.New("Некорректный ввод: " + err)
}
