package models

import (
	"fmt"
)

// FormFeedback – Форма для отзыва пользователя о статье
type FormFeedback struct {
	Answer string `json:"answer"`
	Article string `json:"article_id"`
}

// ToString - Генерация форматированного текста из данных формы
func (p *FormFeedback) ToString() string {
	s := fmt.Sprintf("Отзыв пользователя:\n%s\n%s\n", p.Article, p.Answer)
	return s
}
