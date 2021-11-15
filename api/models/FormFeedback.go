package models

import (
	"fmt"
)

// FormFeedback – Форма для отзыва пользователя о статье
type FormFeedback struct {
	Answer string `json:"answer"`
}

// ToString - Генерация форматированного текста из данных формы
func (p *FormFeedback) ToString() string {
	s := fmt.Sprintf("Отзыв пользователя:\n\n%s\n", p.Answer)
	return s
}
