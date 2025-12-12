package database

import (
	"log"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	"carefinder/models"
)

var DB *gorm.DB

func Connect() {
	var err error
	DB, err = gorm.Open(sqlite.Open("carefinder.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database")
	}

	log.Println("Database connection established")

	// Auto-migrate the schema
	err = DB.AutoMigrate(
		&models.User{},
		&models.CaregiverProfile{},
		&models.License{},
		&models.Availability{},
		&models.Contract{},
		&models.Review{},
	)
	if err != nil {
		log.Fatal("Failed to migrate database")
	}

	log.Println("Database migrated")
}
