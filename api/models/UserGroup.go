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

// UserGroup - произвольная группа пользователей
type UserGroup struct {
	ID        uint64    `gorm:"primary_key;auto_increment" json:"id"`
	Name      string    `gorm:"size:255;not null;unique" json:"nickname"`
	Email     string    `gorm:"size:100;not null;unique" json:"email"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// Prepare - Подготовка информации о группе пользователей
func (u *UserGroup) Prepare() {
	u.ID = 0
	u.Name = html.EscapeString(strings.TrimSpace(u.Name))
	u.Email = html.EscapeString(strings.TrimSpace(u.Email))
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
}

// SaveUserGroup - Сохранение информации о группе пользователей
func (u *UserGroup) SaveUserGroup(db *gorm.DB) (*UserGroup, error) {
	var err = db.Debug().Create(&u).Error
	if err != nil {
		return &UserGroup{}, err
	}
	return u, nil
}

// FindAllUserGroups - Вывод всех групп пользователей (максимальное количество задаётся параметром GET_LIMIT)
func (u *UserGroup) FindAllUserGroups(db *gorm.DB) (*[]UserGroup, error) {
	var err error
	users := []UserGroup{}
	err = db.Debug().Model(&UserGroup{}).Limit(os.Getenv("GET_LIMIT")).Find(&users).Error
	if err != nil {
		return &[]UserGroup{}, err
	}
	return &users, err
}

// FindUserGroupByID - Вывод информации о группе пользователей с ID
func (u *UserGroup) FindUserGroupByID(db *gorm.DB, uid uint64) (*UserGroup, error) {
	var err = db.Debug().Model(User{}).Where("id = ?", uid).Take(&u).Error
	if err != nil {
		return &UserGroup{}, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return &UserGroup{}, errors.New("UserGroup Not Found")
	}
	return u, err
}

// UpdateAUserGroup - Обновление информации о группе пользователей
func (u *UserGroup) UpdateAUserGroup(db *gorm.DB, uid uint64) (*UserGroup, error) {
	db = db.Debug().Model(&UserGroup{}).Where("id = ?", uid).Take(&UserGroup{}).UpdateColumns(
		map[string]interface{}{
			"name":      u.Name,
			"email":     u.Email,
			"update_at": time.Now(),
		},
	)
	if db.Error != nil {
		return &UserGroup{}, db.Error
	}
	// Вывод обновленной информации о группе пользователей
	err := db.Debug().Model(&UserGroup{}).Where("id = ?", uid).Take(&u).Error
	if err != nil {
		return &UserGroup{}, err
	}
	return u, nil
}

// DeleteAUserGroup - Удаление группы пользователей
func (u *UserGroup) DeleteAUserGroup(db *gorm.DB, uid uint64) (int64, error) {
	db = db.Debug().Model(&UserGroup{}).Where("id = ?", uid).Take(&UserGroup{}).Delete(&UserGroup{})
	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
