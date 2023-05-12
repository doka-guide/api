// Package seed - пакет для установки структуры БД и записей по умолчанию
package seed

import (
	"log"
	"os"

	"github.com/doka-guide/api/api/models"
	"github.com/jinzhu/gorm"
)

// Load - загрузка базы данных
func Load(db *gorm.DB) {
	var users = []models.User{
		{
			Nickname: os.Getenv("USER_NAME"),
			Email:    os.Getenv("USER_MAIL"),
			Password: os.Getenv("USER_PASS"),
		},
		{
			Nickname: os.Getenv("ADMIN_NAME"),
			Email:    os.Getenv("ADMIN_MAIL"),
			Password: os.Getenv("ADMIN_PASS"),
		},
	}

	// Создание записей по умолчанию в режиме отладки
	if os.Getenv("MODE") == "DEBUG" {
		err := db.Debug().DropTableIfExists(&models.Form{}, &models.ProfileLink{}, &models.SubscriptionReport{}, &models.Subscription{}, &models.User{}).Error
		if err != nil {
			log.Fatalf("Не удаётся удалить таблицу: %v", err)
		}
		err = db.Debug().AutoMigrate(&models.User{}, &models.Subscription{}, &models.ProfileLink{}, &models.SubscriptionReport{}, &models.Form{}).Error
		if err != nil {
			log.Fatalf("Не удаётся произвести миграцию: %v", err)
		}
		err = db.Debug().Model(&models.Form{}).AddForeignKey("author_id", "users(id)", "cascade", "cascade").Error
		if err != nil {
			log.Fatalf("Установка внешнего ключа завершилась неудачей (forms -> users): %v", err)
		}
		err = db.Debug().Model(&models.Subscription{}).AddForeignKey("author_id", "users(id)", "cascade", "cascade").Error
		if err != nil {
			log.Fatalf("Установка внешнего ключа завершилась неудачей (subscriptions -> users): %v", err)
		}
		err = db.Debug().Model(&models.ProfileLink{}).AddForeignKey("author_id", "users(id)", "cascade", "cascade").Error
		if err != nil {
			log.Fatalf("Установка внешнего ключа завершилась неудачей (profileLinks -> users): %v", err)
		}
		err = db.Debug().Model(&models.ProfileLink{}).AddForeignKey("profile_id", "subscriptions(id)", "cascade", "cascade").Error
		if err != nil {
			log.Fatalf("Установка внешнего ключа завершилась неудачей (profileLinks -> subscriptions): %v", err)
		}
		err = db.Debug().Model(&models.SubscriptionReport{}).AddForeignKey("profile_id", "subscriptions(id)", "cascade", "cascade").Error
		if err != nil {
			log.Fatalf("Установка внешнего ключа завершилась неудачей (subscriptionreport -> users): %v", err)
		}
		err = db.Debug().Model(&models.SubscriptionReport{}).AddForeignKey("profile_id", "subscriptions(id)", "cascade", "cascade").Error
		if err != nil {
			log.Fatalf("Установка внешнего ключа завершилась неудачей (subscriptionreport -> subscriptions): %v", err)
		}
		for i := range users {
			err = db.Debug().Model(&models.User{}).Create(&users[i]).Error
			if err != nil {
				log.Fatalf("Не удаётся добавить пользователей: %v", err)
			}
		}
	}
}
