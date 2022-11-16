// Package models - пакет для описания моделей, которые используются для хранения данных
package models

import (
	"errors"
	"os"
	"time"

	"github.com/jinzhu/gorm"
)

// SubscriptionReport - ссылка на ресурсы, которые запросил пользователь
type SubscriptionReport struct {
	ID        uint64       `gorm:"primary_key;auto_increment" json:"id"`
	Path      string       `gorm:"size:255;not null;" json:"path"`
	Author    User         `json:"author"`
	Profile   Subscription `json:"profile"`
	AuthorID  uint64       `gorm:"not null" json:"author_id"`
	ProfileID uint64       `gorm:"not null" json:"profile_id"`
	CreatedAt time.Time    `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time    `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// Prepare - Подготовка ссылок на ресурсы, которые запросил пользователь
func (p *SubscriptionReport) Prepare() {
	p.ID = 0
	p.Profile = Subscription{}
	p.Author = User{}
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
}

// Validate - Валидация ссылок на ресурсы, которые запросил пользователь
func (p *SubscriptionReport) Validate() error {
	if p.Path == "" {
		return errors.New("Необходимо указать ссылку к ресурсу")
	}
	if p.AuthorID < 1 {
		return errors.New("Необходимо указать ID пользователя")
	}
	return nil
}

// SaveSubscriptionReport - Сохранение ссылок на ресурсы, которые запросил пользователь
func (p *SubscriptionReport) SaveSubscriptionReport(db *gorm.DB) (*SubscriptionReport, error) {
	var err error
	err = db.Debug().Model(&SubscriptionReport{}).Create(&p).Error
	if err != nil {
		return &SubscriptionReport{}, err
	}
	if p.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", p.AuthorID).Take(&p.Author).Error
		if err != nil {
			return &SubscriptionReport{}, err
		}
	}
	return p, nil
}

// FindAllSubscriptionReports - Вывод всех ссылок на ресурсы, которые запросил пользователь (максимальное количество задаётся параметром GET_LIMIT)
func (p *SubscriptionReport) FindAllSubscriptionReports(db *gorm.DB) (*[]SubscriptionReport, error) {
	var err error
	posts := []SubscriptionReport{}
	err = db.Debug().Model(&SubscriptionReport{}).Order("id DESC").Limit(os.Getenv("GET_LIMIT")).Find(&posts).Error
	if err != nil {
		return &[]SubscriptionReport{}, err
	}
	if len(posts) > 0 {
		for i := range posts {
			err := db.Debug().Model(&User{}).Where("id = ?", posts[i].AuthorID).Take(&posts[i].Author).Error
			if err != nil {
				return &[]SubscriptionReport{}, err
			}
		}
	}
	return &posts, nil
}

// FindSubscriptionReportByPath - Вывод данных ссылки на ресурсы, которые запросили пользователи
func (p *SubscriptionReport) FindSubscriptionReportByPath(db *gorm.DB, path string) (*SubscriptionReport, error) {
	var err error
	err = db.Debug().Model(&SubscriptionReport{}).Where("path = ?", path).Take(&p).Error
	if err != nil {
		return &SubscriptionReport{}, err
	}
	if p.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", p.AuthorID).Take(&p.Author).Error
		if err != nil {
			return &SubscriptionReport{}, err
		}
		err = db.Debug().Model(&Subscription{}).Where("id = ?", p.ProfileID).Take(&p.Profile).Error
		if err != nil {
			return &SubscriptionReport{}, err
		}
	}
	return p, nil
}

// DeleteASubscriptionReport - Удаление ссылок на ресурсы, которые запросил пользователь
func (p *SubscriptionReport) DeleteASubscriptionReport(db *gorm.DB, pid uint64, uid uint64) (int64, error) {
	db = db.Debug().Model(&SubscriptionReport{}).Where("id = ? and author_id = ?", pid, uid).Take(&SubscriptionReport{}).Delete(&SubscriptionReport{})
	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, errors.New("SubscriptionReport not found")
		}
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
