// Package seed - пакет для установки структуры БД и записей по умолчанию
package seed

import (
	"log"
	"os"

	"github.com/doka-guide/api/api/models"
	"github.com/jinzhu/gorm"
)

// Load - loading of the DB
func Load(db *gorm.DB) {
	var users = []models.User{
		{
			Nickname: os.Getenv("USER_NAME"),
			Email:    os.Getenv("USER_MAIL"),
			Password: os.Getenv("USER_PASS"),
		},
	}

	if os.Getenv("MODE") == "DEBUG" {
		err := db.Debug().DropTableIfExists(&models.Form{}, &models.ProfileLink{}, &models.SubscriptionReport{}, &models.Subscription{}, &models.User{}).Error
		if err != nil {
			log.Fatalf("cannot drop table: %v", err)
		}
		err = db.Debug().AutoMigrate(&models.User{}, &models.Subscription{}, &models.ProfileLink{}, &models.SubscriptionReport{}, &models.Form{}).Error
		if err != nil {
			log.Fatalf("cannot migrate table: %v", err)
		}
		err = db.Debug().Model(&models.Form{}).AddForeignKey("author_id", "users(id)", "cascade", "cascade").Error
		if err != nil {
			log.Fatalf("attaching foreign key error (forms -> users): %v", err)
		}
		err = db.Debug().Model(&models.Subscription{}).AddForeignKey("author_id", "users(id)", "cascade", "cascade").Error
		if err != nil {
			log.Fatalf("attaching foreign key error (subscriptions -> users): %v", err)
		}
		err = db.Debug().Model(&models.ProfileLink{}).AddForeignKey("author_id", "users(id)", "cascade", "cascade").Error
		if err != nil {
			log.Fatalf("attaching foreign key error (profileLinks -> users): %v", err)
		}
		err = db.Debug().Model(&models.ProfileLink{}).AddForeignKey("profile_id", "subscriptions(id)", "cascade", "cascade").Error
		if err != nil {
			log.Fatalf("attaching foreign key error (profileLinks -> subscriptions): %v", err)
		}
		err = db.Debug().Model(&models.SubscriptionReport{}).AddForeignKey("profile_id", "subscriptions(id)", "cascade", "cascade").Error
		if err != nil {
			log.Fatalf("attaching foreign key error (subscriptionreport -> users): %v", err)
		}
		err = db.Debug().Model(&models.SubscriptionReport{}).AddForeignKey("profile_id", "subscriptions(id)", "cascade", "cascade").Error
		if err != nil {
			log.Fatalf("attaching foreign key error (subscriptionreport -> subscriptions): %v", err)
		}
		for i := range users {
			err = db.Debug().Model(&models.User{}).Create(&users[i]).Error
			if err != nil {
				log.Fatalf("cannot seed users table: %v", err)
			}
		}
	}
}
