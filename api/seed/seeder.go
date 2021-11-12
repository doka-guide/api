package seed

import (
	"log"
	"os"

	"../models"
	"github.com/jinzhu/gorm"
)

var users = []models.User{
	models.User{
		Nickname: os.Getenv("USER_NAME"),
		Email:    os.Getenv("USER_MAIL"),
		Password: os.Getenv("USER_PASS"),
	},
}

// Load - loading of the DB
func Load(db *gorm.DB) {

	if os.Getenv("MODE") == "DEBUG" {
		err := db.Debug().DropTableIfExists(&models.Form{}, &models.User{}).Error
		if err != nil {
			log.Fatalf("cannot drop table: %v", err)
		}
		err = db.Debug().AutoMigrate(&models.User{}, &models.Form{}).Error
		if err != nil {
			log.Fatalf("cannot migrate table: %v", err)
		}
		err = db.Debug().Model(&models.Form{}).AddForeignKey("author_id", "users(id)", "cascade", "cascade").Error
		if err != nil {
			log.Fatalf("attaching foreign key error: %v", err)
		}
		for i := range users {
			err = db.Debug().Model(&models.User{}).Create(&users[i]).Error
			if err != nil {
				log.Fatalf("cannot seed users table: %v", err)
			}
		}
	}
}
