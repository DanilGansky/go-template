package db

import (
	"log"
	"sync"

	"gorm.io/driver/postgres"

	"gorm.io/gorm"
)

var (
	once sync.Once
	db   *gorm.DB
)

func Get(dsn string) *gorm.DB {
	once.Do(func() {
		conn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Fatalf("failed to connect to db: %s", err.Error())
		}

		db = conn
	})

	return db
}
