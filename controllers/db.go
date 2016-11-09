package controllers

import (
	"log"
	"math/rand"
	"time"

	"github.com/deviceMP/api-server/models"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var db *gorm.DB

func init() {
	rand.Seed(time.Now().UnixNano())

	var err error
	db, err = gorm.Open("mysql", "root:1@/iotgatewaydb?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		log.Fatalf("Error occurs when connecting to the database %s", err)
	}

	db.AutoMigrate(
		&models.Org{},
		&models.Project{},
		&models.Device{},
		&models.App{},
		&models.ProjectEnv{},
		&models.DeviceEnv{},
	)
}