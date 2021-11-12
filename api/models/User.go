package models

import (
	"errors"
	"html"
	"log"
	"strings"
	"time"
	"os"

	"github.com/badoux/checkmail"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

// Произвольный пользователь
type User struct {
	ID        uint32    `gorm:"primary_key;auto_increment" json:"id"`
	Nickname  string    `gorm:"size:255;not null;unique" json:"nickname"`
	Email     string    `gorm:"size:100;not null;unique" json:"email"`
	Password  string    `gorm:"size:100;not null;" json:"password"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// Функция хеширования
func Hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

// VerifyPassword - Проверка пароля пользователя
func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// BeforeSave - Подготовка к сохранению пользователя
func (u *User) BeforeSave() error {
	hashedPassword, err := Hash(u.Password)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// Prepare - Подготовка информации о пользователе
func (u *User) Prepare() {
	u.ID = 0
	u.Nickname = html.EscapeString(strings.TrimSpace(u.Nickname))
	u.Email = html.EscapeString(strings.TrimSpace(u.Email))
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
}

// Validate - Валидация учётных данных пользователя
func (u *User) Validate(action string) error {
	switch strings.ToLower(action) {
	case "update":
		if u.Nickname == "" {
			return errors.New("Требуется имя пользователя")
		}
		if u.Password == "" {
			return errors.New("Требуется пароль")
		}
		if u.Email == "" {
			return errors.New("Требуется электронная почта")
		}
		if err := checkmail.ValidateFormat(u.Email); err != nil {
			return errors.New("Такой почты быть не может")
		}

		return nil
	case "login":
		if u.Password == "" {
			return errors.New("Требуется пароль")
		}
		if u.Email == "" {
			return errors.New("Требуется электронная почта")
		}
		if err := checkmail.ValidateFormat(u.Email); err != nil {
			return errors.New("Такой почты быть не может")
		}
		return nil

	default:
		if u.Nickname == "" {
			return errors.New("Требуется имя пользователя")
		}
		if u.Password == "" {
			return errors.New("Требуется пароль")
		}
		if u.Email == "" {
			return errors.New("Требуется электронная почта")
		}
		if err := checkmail.ValidateFormat(u.Email); err != nil {
			return errors.New("Такой почты быть не может")
		}
		return nil
	}
}

// SaveUser - Сохранение информации о пользователе
func (u *User) SaveUser(db *gorm.DB) (*User, error) {

	var err error
	err = db.Debug().Create(&u).Error
	if err != nil {
		return &User{}, err
	}
	return u, nil
}

// FindAllUsers - Вывод всех пользователей (максимальное количество задаётся параметром GET_LIMIT)
func (u *User) FindAllUsers(db *gorm.DB) (*[]User, error) {
	var err error
	users := []User{}
	err = db.Debug().Model(&User{}).Limit(os.Getenv("GET_LIMIT")).Find(&users).Error
	if err != nil {
		return &[]User{}, err
	}
	return &users, err
}

// FindUserByID - Вывод информации о пользователе с ID
func (u *User) FindUserByID(db *gorm.DB, uid uint32) (*User, error) {
	var err error
	err = db.Debug().Model(User{}).Where("id = ?", uid).Take(&u).Error
	if err != nil {
		return &User{}, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return &User{}, errors.New("User Not Found")
	}
	return u, err
}

// UpdateAUser - Обновление информации о пользователе
func (u *User) UpdateAUser(db *gorm.DB, uid uint32) (*User, error) {

	// Хеширование пароля
	err := u.BeforeSave()
	if err != nil {
		log.Fatal(err)
	}
	db = db.Debug().Model(&User{}).Where("id = ?", uid).Take(&User{}).UpdateColumns(
		map[string]interface{}{
			"password":  u.Password,
			"nickname":  u.Nickname,
			"email":     u.Email,
			"update_at": time.Now(),
		},
	)
	if db.Error != nil {
		return &User{}, db.Error
	}
	// Вывод обновленной информации о пользователе
	err = db.Debug().Model(&User{}).Where("id = ?", uid).Take(&u).Error
	if err != nil {
		return &User{}, err
	}
	return u, nil
}

// DeleteAUser - Удаление пользователя
func (u *User) DeleteAUser(db *gorm.DB, uid uint32) (int64, error) {

	db = db.Debug().Model(&User{}).Where("id = ?", uid).Take(&User{}).Delete(&User{})

	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
