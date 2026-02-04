package database

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/lafathalfath/go-chatserver/helpers"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type ConnectionObj struct {
	DB     *gorm.DB
	RDB    *redis.Client
	RDBCtx context.Context
}

var DBConnection = ConnectionObj{RDBCtx: context.Background()}

func Connect() {
	connectDB()
	connectRDB()
}

func connectDB() {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		helpers.Env("DB_HOST"),
		helpers.Env("DB_USER"),
		helpers.Env("DB_PASSWORD"),
		helpers.Env("DB_NAME"),
		helpers.Env("DB_PORT"),
		helpers.Env("SSL_MODE"),
		helpers.Env("TIMEZONE"),
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	autoMigrate(db)

	DBConnection.DB = db
	log.Println("Database Connected")
}

func connectRDB() {
	var address string
	var password string
	address = helpers.Env("REDIS_HOST") + ":" + helpers.Env("REDIS_PORT")
	password = helpers.Env("REDIS_PASSWORD")
	db, err := strconv.Atoi(helpers.Env("REDIS_DB"))
	if err != nil {
		panic(err)
	}
	DBConnection.RDB = redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       db,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := DBConnection.RDB.Ping(ctx).Err(); err != nil {
		panic(err)
	}
	log.Println("Redis Connected")
}
