package model

import (
	"go-zero-api/service/internal/types"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func init() {

	//open a db connection
	db, err := gorm.Open(sqlite.Open("/data/go-zero.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	//Migrate the schema
	db.AutoMigrate(&types.RegisterRequest{})
}
