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

// Permission — произвольные разрешения
type Permission struct {
	ID        uint64    `gorm:"primary_key;auto_increment" json:"id"`
	Name      string    `gorm:"size:255;not null;unique" json:"nickname"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// Prepare - Подготовка информации о произвольных разрешениях
func (u *Permission) Prepare() {
	u.ID = 0
	u.Name = html.EscapeString(strings.TrimSpace(u.Name))
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
}

// SavePermission - Сохранение информации о произвольных разрешениях
func (u *Permission) SavePermission(db *gorm.DB) (*Permission, error) {
	var err = db.Debug().Create(&u).Error
	if err != nil {
		return &Permission{}, err
	}
	return u, nil
}

// FindAllPermissions - Вывод всех произвольных разрешений (максимальное количество задаётся параметром GET_LIMIT)
func (u *Permission) FindAllPermissions(db *gorm.DB) (*[]Permission, error) {
	var err error
	users := []Permission{}
	err = db.Debug().Model(&Permission{}).Limit(os.Getenv("GET_LIMIT")).Find(&users).Error
	if err != nil {
		return &[]Permission{}, err
	}
	return &users, err
}

// FindPermissionByID - Вывод информации о произвольных разрешениях с ID
func (u *Permission) FindPermissionByID(db *gorm.DB, uid uint64) (*Permission, error) {
	var err = db.Debug().Model(User{}).Where("id = ?", uid).Take(&u).Error
	if err != nil {
		return &Permission{}, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return &Permission{}, errors.New("Permission Not Found")
	}
	return u, err
}

// UpdateAPermission - Обновление информации о произвольных разрешениях
func (u *Permission) UpdateAPermission(db *gorm.DB, uid uint64) (*Permission, error) {
	db = db.Debug().Model(&Permission{}).Where("id = ?", uid).Take(&Permission{}).UpdateColumns(
		map[string]interface{}{
			"name":      u.Name,
			"update_at": time.Now(),
		},
	)
	if db.Error != nil {
		return &Permission{}, db.Error
	}
	// Вывод обновленной информации о произвольных разрешениях
	err := db.Debug().Model(&Permission{}).Where("id = ?", uid).Take(&u).Error
	if err != nil {
		return &Permission{}, err
	}
	return u, nil
}

// DeleteAPermission - Удаление произвольных разрешений
func (u *Permission) DeleteAPermission(db *gorm.DB, uid uint64) (int64, error) {
	db = db.Debug().Model(&Permission{}).Where("id = ?", uid).Take(&Permission{}).Delete(&Permission{})
	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
