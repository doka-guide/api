// Package models - пакет для описания моделей, которые используются для хранения данных
package models

import (
	"errors"
	"os"
	"time"

	"github.com/jinzhu/gorm"
)

// GroupPermission - пара группа-разрешение
type GroupPermission struct {
	ID        uint64     `gorm:"primary_key;auto_increment" json:"id"`
	Group     UserGroup  `json:"group"`
	Perms     Permission `json:"permission"`
	GroupID   uint64     `gorm:"not null" json:"group_id"`
	PermsID   uint64     `gorm:"not null" json:"perms_id"`
	CreatedAt time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// Prepare - Подготовка на пар группа-разрешение
func (p *GroupPermission) Prepare() {
	p.ID = 0
	p.Group = UserGroup{}
	p.Perms = Permission{}
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
}

// Validate - Валидация на пару группа-разрешение
func (p *GroupPermission) Validate() error {
	if p.GroupID < 1 {
		return errors.New("Необходимо указать ID пары группа-разрешение")
	}
	return nil
}

// SaveGroupPermission - Сохранение на пару группа-разрешение
func (p *GroupPermission) SaveGroupPermission(db *gorm.DB) (*GroupPermission, error) {
	var err error
	err = db.Debug().Model(&GroupPermission{}).Create(&p).Error
	if err != nil {
		return &GroupPermission{}, err
	}
	if p.ID != 0 {
		err = db.Debug().Model(&UserGroup{}).Where("id = ?", p.GroupID).Take(&p.Group).Error
		if err != nil {
			return &GroupPermission{}, err
		}
	}
	return p, nil
}

// FindAllGroupPermission - Вывод всех пар группа-разрешение (максимальное количество задаётся параметром GET_LIMIT)
func (p *GroupPermission) FindAllGroupPermission(db *gorm.DB) (*[]GroupPermission, error) {
	var err error
	posts := []GroupPermission{}
	err = db.Debug().Model(&GroupPermission{}).Order("id DESC").Limit(os.Getenv("GET_LIMIT")).Find(&posts).Error
	if err != nil {
		return &[]GroupPermission{}, err
	}
	if len(posts) > 0 {
		for i := range posts {
			err := db.Debug().Model(&UserGroup{}).Where("id = ?", posts[i].GroupID).Take(&posts[i].Group).Error
			if err != nil {
				return &[]GroupPermission{}, err
			}
		}
	}
	return &posts, nil
}

// FindGroupPermissionByID - Вывод данных пары группа-разрешение с ID
func (p *GroupPermission) FindGroupPermissionByID(db *gorm.DB, id string) (*GroupPermission, error) {
	var err error
	err = db.Debug().Model(&GroupPermission{}).Where("id = ?", id).Take(&p).Error
	if err != nil {
		return &GroupPermission{}, err
	}
	if p.ID != 0 {
		err = db.Debug().Model(&UserGroup{}).Where("id = ?", p.GroupID).Take(&p.Group).Error
		if err != nil {
			return &GroupPermission{}, err
		}
		err = db.Debug().Model(&Permission{}).Where("id = ?", p.PermsID).Take(&p.Perms).Error
		if err != nil {
			return &GroupPermission{}, err
		}
	}
	return p, nil
}

// DeleteAGroupPermission - Удаление на пар группа-разрешение
func (p *GroupPermission) DeleteAGroupPermission(db *gorm.DB, uid uint64) (int64, error) {
	db = db.Debug().Model(&GroupPermission{}).Where("id = ?", uid).Take(&GroupPermission{}).Delete(&GroupPermission{})
	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, errors.New("GroupPermission not found")
		}
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
