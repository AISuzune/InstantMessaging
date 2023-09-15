package boot

import (
	"InstantMessaging/app/global"
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
	"time"
)

func MysqlDBSetup() {
	config := g.Config.DataBase.Mysql

	db, err := gorm.Open("mysql", config.GetDsn())
	if err != nil {
		log.Printf("initialize mysql db failed, err: %v\n", err)
	}
	db.LogMode(true)
	sqlDB := db.DB()
	sqlDB.SetConnMaxIdleTime(10 * time.Second)
	sqlDB.SetConnMaxLifetime(100 * time.Second)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	err = sqlDB.Ping()
	if err != nil {
		log.Fatalf("connect to mysql db failed, err: %v\n", err)
	}

	fmt.Printf("initialize mysql db successfully\n")
	g.MysqlDB = db

}

func RedisSetup() {
	config := g.Config.DataBase.Redis

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.Addr, config.Port),
		Username: "",
		Password: config.Password,
		DB:       config.Db,
		PoolSize: 10000,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("connect to redis instance failed, err: %v\n", err)
	}

	g.Rdb = rdb

	fmt.Printf("initialize redis client successfully\n")
}
