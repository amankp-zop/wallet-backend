package database

import (
	"database/sql"
	"fmt"
	"time"
	_ "github.com/go-sql-driver/mysql"
)

func NewDatabase(dsn string) (*sql.DB,error){
	db,err := sql.Open("mysql",dsn);
	if err!=nil{
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5*time.Minute)

	err = db.Ping();
	if err!=nil{
		db.Close()
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	return db,nil
}