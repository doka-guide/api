// Package models - пакет для описания моделей, которые используются для хранения данных
package models

import (
	"errors"
	"os"
	"time"

	"github.com/jinzhu/gorm"
)

// GroupedUser - пара группа-пользователей
type GroupedUser struct {
	ID        uint64    `gorm:"primary_key;auto_increment" json:"id"`
	Group     UserGroup `json:"group"`
	User      User      `json:"user"`
	GroupID   uint64    `gorm:"not null" json:"group_id"`
	UserID    uint64    `gorm:"not null" json:"user_id"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// Prepare - Подготовка пары группа-пользователей
func (p *GroupedUser) Prepare() {
	p.ID = 0
	p.Group = UserGroup{}
	p.User = User{}
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
}

// Validate - Валидация пары группа-пользователей
func (p *GroupedUser) Validate() error {
	if p.GroupID < 1 {
		return errors.New("Необходимо указать ID пары группа-пользователей")
	}
	return nil
}

// SaveGroupedUser - Сохранение пары группа-пользователей
func (p *GroupedUser) SaveGroupedUser(db *gorm.DB) (*GroupedUser, error) {
	var err error
	err = db.Debug().Model(&GroupedUser{}).Create(&p).Error
	if err != nil {
		return &GroupedUser{}, err
	}
	if p.ID != 0 {
		err = db.Debug().Model(&UserGroup{}).Where("id = ?", p.GroupID).Take(&p.Group).Error
		if err != nil {
			return &GroupedUser{}, err
		}
	}
	return p, nil
}

// FindAllGroupedUser - Вывод всех пар группа-пользователей (максимальное количество задаётся параметром GET_LIMIT)
func (p *GroupedUser) FindAllGroupedUser(db *gorm.DB) (*[]GroupedUser, error) {
	var err error
	posts := []GroupedUser{}
	err = db.Debug().Model(&GroupedUser{}).Order("id DESC").Limit(os.Getenv("GET_LIMIT")).Find(&posts).Error
	if err != nil {
		return &[]GroupedUser{}, err
	}
	if len(posts) > 0 {
		for i := range posts {
			err := db.Debug().Model(&UserGroup{}).Where("id = ?", posts[i].GroupID).Take(&posts[i].Group).Error
			if err != nil {
				return &[]GroupedUser{}, err
			}
		}
	}
	return &posts, nil
}

// FindAllGroupedUserWithUserID - Вывод всех пар группа-пользователей с определённым User ID
func (p *GroupedUser) FindAllGroupedUserWithUserID(db *gorm.DB, id string) (*[]GroupedUser, error) {
	var err error
	posts := []GroupedUser{}
	err = db.Debug().Model(&GroupedUser{}).Where("user_id = ?", id).Order("id DESC").Find(&posts).Error
	if err != nil {
		return &[]GroupedUser{}, err
	}
	if len(posts) > 0 {
		for i := range posts {
			err := db.Debug().Model(&UserGroup{}).Where("id = ?", posts[i].GroupID).Take(&posts[i].Group).Error
			if err != nil {
				return &[]GroupedUser{}, err
			}
			err = db.Debug().Model(&User{}).Where("id = ?", posts[i].UserID).Take(&posts[i].User).Error
			if err != nil {
				return &[]GroupedUser{}, err
			}
		}
	}
	return &posts, nil
}

// FindAllGroupedUserWithUserGroupID - Вывод всех пар группа-пользователей с определённым UserGroup ID
func (p *GroupedUser) FindAllGroupedUserWithUserGroupID(db *gorm.DB, id string) (*[]GroupedUser, error) {
	var err error
	posts := []GroupedUser{}
	err = db.Debug().Model(&GroupedUser{}).Where("group_id = ?", id).Order("id DESC").Find(&posts).Error
	if err != nil {
		return &[]GroupedUser{}, err
	}
	if len(posts) > 0 {
		for i := range posts {
			err := db.Debug().Model(&UserGroup{}).Where("id = ?", posts[i].GroupID).Take(&posts[i].Group).Error
			if err != nil {
				return &[]GroupedUser{}, err
			}
			err = db.Debug().Model(&User{}).Where("id = ?", posts[i].UserID).Take(&posts[i].User).Error
			if err != nil {
				return &[]GroupedUser{}, err
			}
		}
	}
	return &posts, nil
}

// FindGroupedUserByID - Вывод данных пары группа-пользователь с ID
func (p *GroupedUser) FindGroupedUserByID(db *gorm.DB, id string) (*GroupedUser, error) {
	var err error
	err = db.Debug().Model(&GroupedUser{}).Where("id = ?", id).Take(&p).Error
	if err != nil {
		return &GroupedUser{}, err
	}
	if p.ID != 0 {
		err = db.Debug().Model(&UserGroup{}).Where("id = ?", p.GroupID).Take(&p.Group).Error
		if err != nil {
			return &GroupedUser{}, err
		}
		err = db.Debug().Model(&User{}).Where("id = ?", p.UserID).Take(&p.User).Error
		if err != nil {
			return &GroupedUser{}, err
		}
	}
	return p, nil
}

// DeleteAGroupedUser - Удаление пары группа-пользователей
func (p *GroupedUser) DeleteAGroupedUser(db *gorm.DB, uid uint64) (int64, error) {
	db = db.Debug().Model(&GroupedUser{}).Where("id = ?", uid).Take(&GroupedUser{}).Delete(&GroupedUser{})
	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, errors.New("GroupedUser not found")
		}
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
