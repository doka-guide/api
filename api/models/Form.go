package models

import (
	"errors"
	"html"
	"strings"
	"time"
	"os"

	"github.com/jinzhu/gorm"
)

// Произвольная форма
type Form struct {
	ID        uint64    `gorm:"primary_key;auto_increment" json:"id"`
	Type      string    `gorm:"size:255;not null;" json:"type"`
	Data      string    `gorm:"type:JSONB;not null;" json:"data"`
	Author    User      `json:"author"`
	AuthorID  uint32    `gorm:"not null" json:"author_id"`
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

// FindAllForms - Вывести все формы (максимальное количество задаётся параметром GET_LIMIT)
func (p *Form) FindAllForms(db *gorm.DB) (*[]Form, error) {
	var err error
	posts := []Form{}
	err = db.Debug().Model(&Form{}).Limit(os.Getenv("GET_LIMIT")).Find(&posts).Error
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

// FindFormByID - Вывести данные формы с ID
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

// UpdateAForm - Обновить форму
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

// DeleteAForm - Удалить форму
func (p *Form) DeleteAForm(db *gorm.DB, pid uint64, uid uint32) (int64, error) {
	db = db.Debug().Model(&Form{}).Where("id = ? and author_id = ?", pid, uid).Take(&Form{}).Delete(&Form{})
	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, errors.New("Form not found")
		}
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
