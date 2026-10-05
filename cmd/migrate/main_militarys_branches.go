package main

import (
	"log"

	"army/internal/app/ds"
	"army/internal/app/dsn"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	_ = godotenv.Load()

	db, err := gorm.Open(
		postgres.Open(dsn.FromEnv()),
		&gorm.Config{},
	)
	if err != nil {
		log.Fatal("failed to connect database")
	}

	err = db.AutoMigrate(
		&ds.User{},
		&ds.MilitaryBranch{},
		&ds.MilitaryBranchLike{},
	)
	if err != nil {
		log.Fatal("cant migrate db")
	}

	log.Println("database migration completed")
}
