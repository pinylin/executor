package db

import (
	"database/sql"
	"fmt"
	"os"
	"sync"
	"time"
)

var (
	db   *sql.DB
	once sync.Once
	mu   sync.Mutex
)

func New() {
	if db == nil {
		once.Do(InitPostgres)
	} else {
		fmt.Println("db already exists")
	}
}

func GetDB() *sql.DB {
	if db == nil {
		InitPostgres()
	}
	return db
}

func InitPostgres() {
	mu.Lock()
	defer mu.Unlock()
	if db != nil {
		return
	}
	var err error

	//"user=pqgotest dbname=pqgotest sslmode=verify-full"

	db, err = sql.Open("postgres", os.Getenv("POSTGRES"))
	if err != nil {
		errMsg := fmt.Sprintf("数据库连接失败1: %v", err)
		fmt.Printf(errMsg)
		panic(err)
	}
	db.SetConnMaxLifetime(time.Second * 10)
	if err = db.Ping(); err != nil {
		errMsg := fmt.Sprintf("数据库连接失败2: %v", err)
		fmt.Printf(errMsg)
		panic(err)
	}
}
