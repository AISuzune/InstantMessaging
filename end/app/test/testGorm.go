package test

import (
	"InstantMessaging/app/global"
	"InstantMessaging/app/internal/model"
	"fmt"
	"github.com/jinzhu/gorm"
	"log"
	"time"
)

func main() {
	config := g.Config.DataBase.Mysql

	db, err := gorm.Open("mysql", config.GetDsn())
	if err != nil {
		log.Printf("initialize mysql db failed, err: %v", err)
	}

	sqlDB := db.DB()
	sqlDB.SetConnMaxIdleTime(10 * time.Second)
	sqlDB.SetConnMaxLifetime(100 * time.Second)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	err = sqlDB.Ping()
	if err != nil {
		fmt.Printf("connect to mysql db failed, err: %v", err)
	}

	fmt.Printf("initialize mysql db successfully")
	g.MysqlDB = db

	db.AutoMigrate(&model.Contact{})
}
