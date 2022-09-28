// Package models - пакет для описания моделей, которые используются для хранения данных
package models

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"html"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

// Link - произвольная форма
type Link struct {
	ID        uint64    `gorm:"primary_key;auto_increment" json:"id"`
	Hash      string    `gorm:"size:255;not null;unique;" json:"hash"`
	Data      string    `gorm:"type:JSONB;not null;" json:"data"`
	Author    User      `json:"author"`
	AuthorID  uint32    `gorm:"not null" json:"author_id"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const stringLength int = 1024

// getRandomString — генерация строки определённой длинны из случайных символов из набора
func getStringWithCharset(length int, charset string) string {
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

// Prepare - Подготовка ссылок
func (p *Link) Prepare() {
	p.ID = 0
	if p.Hash == "" {
		hash := sha256.Sum256([]byte(getStringWithCharset(stringLength, charset)))
		p.Hash = fmt.Sprintf("%x", hash[:])
	} else {
		hash := sha256.Sum256([]byte(html.EscapeString(strings.TrimSpace(p.Hash))))
		p.Hash = fmt.Sprintf("%x", hash[:])
	}
	p.Data = strings.Replace(string([]byte(html.EscapeString(strings.TrimSpace(p.Data)))), "&#34;", "\"", -1)
	p.Author = User{}
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
}

// Validate - Валидация ссылок
func (p *Link) Validate() error {
	if p.Hash == "" {
		return errors.New("Необходимо сгенерировать хэш для ссылки")
	}
	if p.Data == "" {
		return errors.New("Нужна дата и время отправки ссылки")
	}
	if p.AuthorID < 1 {
		return errors.New("Необходимо указать ID пользователя")
	}
	return nil
}

// SaveLink - Сохранение ссылок
func (p *Link) SaveLink(db *gorm.DB) (*Link, error) {
	var err error
	err = db.Debug().Model(&Link{}).Create(&p).Error
	if err != nil {
		return &Link{}, err
	}
	if p.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", p.AuthorID).Take(&p.Author).Error
		if err != nil {
			return &Link{}, err
		}
	}
	return p, nil
}

// FindAllLinks - Вывод всех ссылок (максимальное количество задаётся параметром GET_LIMIT)
func (p *Link) FindAllLinks(db *gorm.DB) (*[]Link, error) {
	var err error
	posts := []Link{}
	err = db.Debug().Model(&Link{}).Order("id DESC").Limit(os.Getenv("GET_LIMIT")).Find(&posts).Error
	if err != nil {
		return &[]Link{}, err
	}
	if len(posts) > 0 {
		for i := range posts {
			err := db.Debug().Model(&User{}).Where("id = ?", posts[i].AuthorID).Take(&posts[i].Author).Error
			if err != nil {
				return &[]Link{}, err
			}
		}
	}
	return &posts, nil
}

// FindLinkByHash - Вывод данных ссылки с Hash
func (p *Link) FindLinkByHash(db *gorm.DB, hash string) (*Link, error) {
	var err error
	err = db.Debug().Model(&Link{}).Where("hash = ?", hash).Take(&p).Error
	if err != nil {
		return &Link{}, err
	}
	if p.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", p.AuthorID).Take(&p.Author).Error
		if err != nil {
			return &Link{}, err
		}
	}
	return p, nil
}

// DeleteALink - Удаление ссылок
func (p *Link) DeleteALink(db *gorm.DB, pid uint64, uid uint32) (int64, error) {
	db = db.Debug().Model(&Link{}).Where("id = ? and author_id = ?", pid, uid).Take(&Link{}).Delete(&Link{})
	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, errors.New("Link not found")
		}
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
