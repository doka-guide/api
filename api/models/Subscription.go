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

// Subscription - форма подписки
type Subscription struct {
	ID        uint64    `gorm:"primary_key;auto_increment" json:"id"`
	Email     string    `gorm:"size:255;not null;" json:"email"`
	Data      string    `gorm:"type:JSONB;not null;" json:"data"`
	Author    User      `json:"author"`
	AuthorID  uint64    `gorm:"not null" json:"author_id"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// Prepare - Подготовка подписки
func (p *Subscription) Prepare() {
	p.ID = 0
	p.Email = html.EscapeString(strings.TrimSpace(p.Email))
	p.Data = strings.Replace(string([]byte(html.EscapeString(strings.TrimSpace(p.Data)))), "&#34;", "\"", -1)
	p.Author = User{}
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
}

// Validate - Валидация подписки
func (p *Subscription) Validate() error {
	if p.Email == "" {
		return errors.New("Необходимо указать электронную почту для подписки")
	}
	if p.Data == "" {
		return errors.New("Нужна дата и время отправки подписки")
	}
	if p.AuthorID < 1 {
		return errors.New("Необходимо указать ID пользователя")
	}
	return nil
}

// SaveSubscription - Сохранение подписки
func (p *Subscription) SaveSubscription(db *gorm.DB) (*Subscription, error) {
	var err error
	err = db.Debug().Model(&Subscription{}).Create(&p).Error
	if err != nil {
		return &Subscription{}, err
	}
	if p.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", p.AuthorID).Take(&p.Author).Error
		if err != nil {
			return &Subscription{}, err
		}
	}
	return p, nil
}

// FindAllSubscriptions - Вывод все подписки (максимальное количество задаётся параметром GET_LIMIT)
func (p *Subscription) FindAllSubscriptions(db *gorm.DB) (*[]Subscription, error) {
	var err error
	posts := []Subscription{}
	err = db.Debug().Model(&Subscription{}).Order("id DESC").Limit(os.Getenv("GET_LIMIT")).Find(&posts).Error
	if err != nil {
		return &[]Subscription{}, err
	}
	if len(posts) > 0 {
		for i := range posts {
			err := db.Debug().Model(&User{}).Where("id = ?", posts[i].AuthorID).Take(&posts[i].Author).Error
			if err != nil {
				return &[]Subscription{}, err
			}
		}
	}
	return &posts, nil
}

// FindSubscriptionByID - Вывод данных подписки с ID
func (p *Subscription) FindSubscriptionByID(db *gorm.DB, pid uint64) (*Subscription, error) {
	var err error
	err = db.Debug().Model(&Subscription{}).Where("id = ?", pid).Take(&p).Error
	if err != nil {
		return &Subscription{}, err
	}
	if p.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", p.AuthorID).Take(&p.Author).Error
		if err != nil {
			return &Subscription{}, err
		}
	}
	return p, nil
}

// UpdateASubscription - Обновление подписки
func (p *Subscription) UpdateASubscription(db *gorm.DB) (*Subscription, error) {
	var err error
	err = db.Debug().Model(&Subscription{}).Where("id = ?", p.ID).Updates(Subscription{Email: p.Email, Data: p.Data, UpdatedAt: time.Now()}).Error
	if err != nil {
		return &Subscription{}, err
	}
	if p.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", p.AuthorID).Take(&p.Author).Error
		if err != nil {
			return &Subscription{}, err
		}
	}
	return p, nil
}

// DeleteASubscription - Удаление подписки
func (p *Subscription) DeleteASubscription(db *gorm.DB, pid uint64, uid uint64) (int64, error) {
	db = db.Debug().Model(&Subscription{}).Where("id = ? and author_id = ?", pid, uid).Take(&Subscription{}).Delete(&Subscription{})
	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, errors.New("Subscription not found")
		}
		return 0, db.Error
	}
	return db.RowsAffected, nil
}

type SubscriptionFormsWithHashResult struct {
	Email string `json:"email"`
	Hash  string `json:"hash"`
	Data  string `json:"data"`
}

// SubscriptionFormsWithHash - Вывод адресов электронной почты и настроек с указанием хэша
func (p *Form) SubscriptionFormsWithHash(db *gorm.DB, start string, end string) *[]SubscriptionFormsWithHashResult {
	posts := []SubscriptionFormsWithHashResult{}
	db.Raw("SELECT (email,hash,data) FROM subscriptions JOIN profile_links ON subscriptions.id=profile_links.profile_id WHERE subscriptions.created_at >= ? AND subscriptions.created_at <= ? ORDER BY subscriptions.created_at ASC", start, end).Scan(&posts)
	return &posts
}
