package initializers

import "github.com/variant-abdulaziz/models"

func SyncDB() {
	DB.AutoMigrate(&models.User{})
}
