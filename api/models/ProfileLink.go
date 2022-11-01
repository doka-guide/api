// Package models - пакет для описания моделей, которые используются для хранения данных
package models

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"html"
	"os"
	"strings"
	"time"

	"github.com/doka-guide/api/api/utils/randomize"
	"github.com/jinzhu/gorm"
)

// ProfileLink - ссылка на профиль подписчика
type ProfileLink struct {
	ID        uint64       `gorm:"primary_key;auto_increment" json:"id"`
	Hash      string       `gorm:"size:255;not null;unique;" json:"hash"`
	Author    User         `json:"author"`
	Profile   Subscription `json:"profile"`
	AuthorID  uint64       `gorm:"not null" json:"author_id"`
	ProfileID uint64       `gorm:"not null" json:"profile_id"`
	CreatedAt time.Time    `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time    `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// Prepare - Подготовка ссылок на профили подписчиков
func (p *ProfileLink) Prepare() {
	p.ID = 0
	if p.Hash == "" {
		hash := sha256.Sum256([]byte(randomize.GetRandomString(1024)))
		p.Hash = fmt.Sprintf("%x", hash[:])
	} else {
		hash := sha256.Sum256([]byte(html.EscapeString(strings.TrimSpace(p.Hash))))
		p.Hash = fmt.Sprintf("%x", hash[:])
	}
	p.Profile = Subscription{}
	p.Author = User{}
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
}

// Validate - Валидация ссылок на профили подписчиков
func (p *ProfileLink) Validate() error {
	if p.Hash == "" {
		return errors.New("Необходимо сгенерировать хэш для ссылки")
	}
	if p.AuthorID < 1 {
		return errors.New("Необходимо указать ID пользователя")
	}
	return nil
}

// SaveProfileLink - Сохранение ссылок на профили подписчиков
func (p *ProfileLink) SaveProfileLink(db *gorm.DB) (*ProfileLink, error) {
	var err error
	err = db.Debug().Model(&ProfileLink{}).Create(&p).Error
	if err != nil {
		return &ProfileLink{}, err
	}
	if p.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", p.AuthorID).Take(&p.Author).Error
		if err != nil {
			return &ProfileLink{}, err
		}
	}
	return p, nil
}

// FindAllProfileLinks - Вывод всех ссылок  на профили подписчиков (максимальное количество задаётся параметром GET_LIMIT)
func (p *ProfileLink) FindAllProfileLinks(db *gorm.DB) (*[]ProfileLink, error) {
	var err error
	posts := []ProfileLink{}
	err = db.Debug().Model(&ProfileLink{}).Order("id DESC").Limit(os.Getenv("GET_LIMIT")).Find(&posts).Error
	if err != nil {
		return &[]ProfileLink{}, err
	}
	if len(posts) > 0 {
		for i := range posts {
			err := db.Debug().Model(&User{}).Where("id = ?", posts[i].AuthorID).Take(&posts[i].Author).Error
			if err != nil {
				return &[]ProfileLink{}, err
			}
		}
	}
	return &posts, nil
}

// FindProfileLinkByHash - Вывод данных ссылки на профиль подписчика с Hash
func (p *ProfileLink) FindProfileLinkByHash(db *gorm.DB, hash string) (*ProfileLink, error) {
	var err error
	err = db.Debug().Model(&ProfileLink{}).Where("hash = ?", hash).Take(&p).Error
	if err != nil {
		return &ProfileLink{}, err
	}
	if p.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", p.AuthorID).Take(&p.Author).Error
		if err != nil {
			return &ProfileLink{}, err
		}
		err = db.Debug().Model(&Subscription{}).Where("id = ?", p.ProfileID).Take(&p.Profile).Error
		if err != nil {
			return &ProfileLink{}, err
		}
	}
	return p, nil
}

// DeleteAProfileLink - Удаление ссылок на профили подписчиков
func (p *ProfileLink) DeleteAProfileLink(db *gorm.DB, pid uint64, uid uint64) (int64, error) {
	db = db.Debug().Model(&ProfileLink{}).Where("id = ? and author_id = ?", pid, uid).Take(&ProfileLink{}).Delete(&ProfileLink{})
	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, errors.New("ProfileLink not found")
		}
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
