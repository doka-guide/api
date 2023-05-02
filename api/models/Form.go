// Package models - пакет для описания моделей, которые используются для хранения данных
package models

import (
	"errors"
	"html"
	"os"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

// Form - произвольная форма
type Form struct {
	ID        uint64    `gorm:"primary_key;auto_increment" json:"id"`
	Type      string    `gorm:"size:255;not null;" json:"type"`
	Data      string    `gorm:"type:JSONB;not null;" json:"data"`
	Author    User      `json:"author"`
	AuthorID  uint64    `gorm:"not null" json:"author_id"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// Prepare - Подготовка формы
func (p *Form) Prepare() {
	p.ID = 0
	p.Type = html.EscapeString(strings.TrimSpace(p.Type))
	p.Data = strings.Replace(string([]byte(html.EscapeString(strings.TrimSpace(p.Data)))), "&#34;", "\"", -1)
	p.Author = User{}
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
}

// Validate - Валидация формы
func (p *Form) Validate() error {
	if p.Type == "" {
		return errors.New("Необходимо указать тип формы")
	}
	if p.Data == "" {
		return errors.New("Нужна дата и время отправки формы")
	}
	if p.AuthorID < 1 {
		return errors.New("Необходимо указать ID пользователя")
	}
	return nil
}

// SaveForm - Сохранение формы
func (p *Form) SaveForm(db *gorm.DB) (*Form, error) {
	var err error
	err = db.Debug().Model(&Form{}).Create(&p).Error
	if err != nil {
		return &Form{}, err
	}
	if p.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", p.AuthorID).Take(&p.Author).Error
		if err != nil {
			return &Form{}, err
		}
	}
	return p, nil
}

// FindAllForms - Вывод все формы (максимальное количество задаётся параметром GET_LIMIT)
func (p *Form) FindAllForms(db *gorm.DB) (*[]Form, error) {
	var err error
	posts := []Form{}
	err = db.Debug().Model(&Form{}).Order("id DESC").Limit(os.Getenv("GET_LIMIT")).Find(&posts).Error
	if err != nil {
		return &[]Form{}, err
	}
	if len(posts) > 0 {
		for i := range posts {
			err := db.Debug().Model(&User{}).Where("id = ?", posts[i].AuthorID).Take(&posts[i].Author).Error
			if err != nil {
				return &[]Form{}, err
			}
		}
	}
	return &posts, nil
}

// FindFormByID - Вывод данных формы с ID
func (p *Form) FindFormByID(db *gorm.DB, pid uint64) (*Form, error) {
	var err error
	err = db.Debug().Model(&Form{}).Where("id = ?", pid).Take(&p).Error
	if err != nil {
		return &Form{}, err
	}
	if p.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", p.AuthorID).Take(&p.Author).Error
		if err != nil {
			return &Form{}, err
		}
	}
	return p, nil
}

// UpdateAForm - Обновление формы
func (p *Form) UpdateAForm(db *gorm.DB) (*Form, error) {
	var err error
	err = db.Debug().Model(&Form{}).Where("id = ?", p.ID).Updates(Form{Type: p.Type, Data: p.Data, UpdatedAt: time.Now()}).Error
	if err != nil {
		return &Form{}, err
	}
	if p.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", p.AuthorID).Take(&p.Author).Error
		if err != nil {
			return &Form{}, err
		}
	}
	return p, nil
}

// DeleteAForm - Удаление формы
func (p *Form) DeleteAForm(db *gorm.DB, pid uint64, uid uint64) (int64, error) {
	db = db.Debug().Model(&Form{}).Where("id = ? and author_id = ?", pid, uid).Take(&Form{}).Delete(&Form{})
	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, errors.New("Form not found")
		}
		return 0, db.Error
	}
	return db.RowsAffected, nil
}

type FormsGroupedByDataResult struct {
	Data  string `gorm:"type:JSONB;not null;" json:"data"`
	Count int    `gorm:"not null" json:"id"`
}

// FeedbackFormsGroupedByData - Вывод агрегированных данных по лайкам / замечаниям для материалов
func (p *Form) FeedbackFormsGroupedByData(db *gorm.DB, start string, end string) *[]FormsGroupedByDataResult {
	posts := []FormsGroupedByDataResult{}
	db.Raw("SELECT (data,count(data)) FROM forms WHERE type = 'feedback' AND created_at >= ? AND created_at <= ? GROUP BY data", start, end).Scan(&posts)
	return &posts
}

type QuestionFormsResult struct {
	Data string `gorm:"type:JSONB;not null;" json:"data"`
}

// QuestionForms - Вывод агрегированных данных по лайкам / замечаниям для материалов
func (p *Form) QuestionForms(db *gorm.DB, start string, end string) *[]QuestionFormsResult {
	posts := []QuestionFormsResult{}
	db.Raw("SELECT data FROM forms WHERE type = 'question' AND created_at >= ? AND created_at <= ?", start, end).Scan(&posts)
	return &posts
}
