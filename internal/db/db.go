package db

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"sync"
)

var Db *sql.DB
var dbOnce sync.Once
var dbErr error

// Подключение к БД
func ConnectionDb() (*sql.DB, error) {
	dbOnce.Do(func() {
		Db, dbErr = sql.Open("mysql", "root:@/pet_projects") //newuser:password@/pet_projects
		if dbErr != nil {
			log.Printf("Ошибка подключения к базе данных: %v", dbErr)
			return
		}
	})

	return Db, dbErr
}
