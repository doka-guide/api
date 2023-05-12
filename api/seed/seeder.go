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
	// Пользователи по умолчанию
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

	// Группы пользователей по умолчанию
	var groups = []models.UserGroup{
		{
			Name:  os.Getenv("USER_GROUP_NAME"),
			Email: os.Getenv("USER_GROUP_MAIL"),
		},
		{
			Name:  os.Getenv("ADMIN_GROUP_NAME"),
			Email: os.Getenv("ADMIN_GROUP_MAIL"),
		},
	}

	var groupedUsers = []models.GroupedUser{
		{
			GroupID: 1,
			UserID:  1,
		},
		{
			GroupID: 2,
			UserID:  2,
		},
	}

	// Права на запросы для групп пользователей
	var permissions = []models.Permission{
		{Name: os.Getenv("PERMISSION_REQUEST_OPTIONS") + os.Getenv("PERMISSION_ENTITY_USER")},
		{Name: os.Getenv("PERMISSION_REQUEST_GET") + os.Getenv("PERMISSION_ENTITY_USER")},
		{Name: os.Getenv("PERMISSION_REQUEST_POST") + os.Getenv("PERMISSION_ENTITY_USER")},
		{Name: os.Getenv("PERMISSION_REQUEST_PUT") + os.Getenv("PERMISSION_ENTITY_USER")},
		{Name: os.Getenv("PERMISSION_REQUEST_DELETE") + os.Getenv("PERMISSION_ENTITY_USER")},

		{Name: os.Getenv("PERMISSION_REQUEST_OPTIONS") + os.Getenv("PERMISSION_ENTITY_FORM")},
		{Name: os.Getenv("PERMISSION_REQUEST_GET") + os.Getenv("PERMISSION_ENTITY_FORM")},
		{Name: os.Getenv("PERMISSION_REQUEST_POST") + os.Getenv("PERMISSION_ENTITY_FORM")},
		{Name: os.Getenv("PERMISSION_REQUEST_PUT") + os.Getenv("PERMISSION_ENTITY_FORM")},
		{Name: os.Getenv("PERMISSION_REQUEST_DELETE") + os.Getenv("PERMISSION_ENTITY_FORM")},

		{Name: os.Getenv("PERMISSION_REQUEST_OPTIONS") + os.Getenv("PERMISSION_ENTITY_PROFILE_LINK")},
		{Name: os.Getenv("PERMISSION_REQUEST_GET") + os.Getenv("PERMISSION_ENTITY_PROFILE_LINK")},
		{Name: os.Getenv("PERMISSION_REQUEST_POST") + os.Getenv("PERMISSION_ENTITY_PROFILE_LINK")},
		{Name: os.Getenv("PERMISSION_REQUEST_PUT") + os.Getenv("PERMISSION_ENTITY_PROFILE_LINK")},
		{Name: os.Getenv("PERMISSION_REQUEST_DELETE") + os.Getenv("PERMISSION_ENTITY_PROFILE_LINK")},

		{Name: os.Getenv("PERMISSION_REQUEST_OPTIONS") + os.Getenv("PERMISSION_ENTITY_SUBSCRIPTION")},
		{Name: os.Getenv("PERMISSION_REQUEST_GET") + os.Getenv("PERMISSION_ENTITY_SUBSCRIPTION")},
		{Name: os.Getenv("PERMISSION_REQUEST_POST") + os.Getenv("PERMISSION_ENTITY_SUBSCRIPTION")},
		{Name: os.Getenv("PERMISSION_REQUEST_PUT") + os.Getenv("PERMISSION_ENTITY_SUBSCRIPTION")},
		{Name: os.Getenv("PERMISSION_REQUEST_DELETE") + os.Getenv("PERMISSION_ENTITY_SUBSCRIPTION")},

		{Name: os.Getenv("PERMISSION_REQUEST_OPTIONS") + os.Getenv("PERMISSION_ENTITY_SUBSCRIPTION_REPORT")},
		{Name: os.Getenv("PERMISSION_REQUEST_GET") + os.Getenv("PERMISSION_ENTITY_SUBSCRIPTION_REPORT")},
		{Name: os.Getenv("PERMISSION_REQUEST_POST") + os.Getenv("PERMISSION_ENTITY_SUBSCRIPTION_REPORT")},
		{Name: os.Getenv("PERMISSION_REQUEST_PUT") + os.Getenv("PERMISSION_ENTITY_SUBSCRIPTION_REPORT")},
		{Name: os.Getenv("PERMISSION_REQUEST_DELETE") + os.Getenv("PERMISSION_ENTITY_SUBSCRIPTION_REPORT")},
	}

	var groupPermissions = []models.GroupPermission{
		{
			GroupID: 1,
			PermsID: 1,
		},
		{
			GroupID: 1,
			PermsID: 3,
		},
		{
			GroupID: 2,
			PermsID: 1,
		},
		{
			GroupID: 2,
			PermsID: 2,
		},
		{
			GroupID: 2,
			PermsID: 3,
		},
		{
			GroupID: 2,
			PermsID: 4,
		},
		{
			GroupID: 2,
			PermsID: 5,
		},
	}

	// Создание записей по умолчанию в режиме отладки
	if os.Getenv("MODE") == "DEBUG" {
		// Удаление таблиц из базы данных
		err := db.Debug().DropTableIfExists(&models.Form{}, &models.ProfileLink{}, &models.SubscriptionReport{}, &models.Subscription{}, &models.GroupedUser{}, &models.User{}, &models.GroupPermission{}, &models.UserGroup{}, &models.Permission{}).Error
		if err != nil {
			log.Fatalf("Не удаётся удалить таблицу: %v", err)
		}

		// Автоматическая миграция  схемы базы данных
		err = db.Debug().AutoMigrate(&models.User{}, &models.UserGroup{}, &models.GroupedUser{}, &models.Permission{}, &models.GroupPermission{}, &models.Subscription{}, &models.ProfileLink{}, &models.SubscriptionReport{}, &models.Form{}).Error
		if err != nil {
			log.Fatalf("Не удаётся произвести миграцию: %v", err)
		}

		// Установка внешних ключей (связей)
		err = db.Debug().Model(&models.GroupedUser{}).AddForeignKey("user_id", "users(id)", "cascade", "cascade").Error
		if err != nil {
			log.Fatalf("Установка внешнего ключа завершилась неудачей (grouped_users -> users): %v", err)
		}
		err = db.Debug().Model(&models.GroupedUser{}).AddForeignKey("group_id", "user_groups(id)", "cascade", "cascade").Error
		if err != nil {
			log.Fatalf("Установка внешнего ключа завершилась неудачей (grouped_users -> user_group): %v", err)
		}
		err = db.Debug().Model(&models.GroupPermission{}).AddForeignKey("perms_id", "permissions(id)", "cascade", "cascade").Error
		if err != nil {
			log.Fatalf("Установка внешнего ключа завершилась неудачей (group_permissions -> permissions): %v", err)
		}
		err = db.Debug().Model(&models.GroupPermission{}).AddForeignKey("group_id", "user_groups(id)", "cascade", "cascade").Error
		if err != nil {
			log.Fatalf("Установка внешнего ключа завершилась неудачей (group_permissions -> user_group): %v", err)
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
			log.Fatalf("Установка внешнего ключа завершилась неудачей (subscription report -> users): %v", err)
		}
		err = db.Debug().Model(&models.SubscriptionReport{}).AddForeignKey("profile_id", "subscriptions(id)", "cascade", "cascade").Error
		if err != nil {
			log.Fatalf("Установка внешнего ключа завершилась неудачей (subscription report -> subscriptions): %v", err)
		}

		// Запись записей по умолчанию
		for i := range users {
			err = db.Debug().Model(&models.User{}).Create(&users[i]).Error
			if err != nil {
				log.Fatalf("Не удаётся добавить пользователей: %v", err)
			}
		}
		for i := range groups {
			err = db.Debug().Model(&models.UserGroup{}).Create(&groups[i]).Error
			if err != nil {
				log.Fatalf("Не удаётся добавить группы пользователей: %v", err)
			}
		}
		for i := range groupedUsers {
			err = db.Debug().Model(&models.GroupedUser{}).Create(&groupedUsers[i]).Error
			if err != nil {
				log.Fatalf("Не удаётся добавить соответствие пользователя и группы пользователей: %v", err)
			}
		}
		for i := range permissions {
			err = db.Debug().Model(&models.Permission{}).Create(&permissions[i]).Error
			if err != nil {
				log.Fatalf("Не удаётся добавить права пользователей: %v", err)
			}
		}
		for i := range groupPermissions {
			err = db.Debug().Model(&models.GroupPermission{}).Create(&groupPermissions[i]).Error
			if err != nil {
				log.Fatalf("Не удаётся добавить соответствие группы пользователей и прав пользователей: %v", err)
			}
		}
	}
}
