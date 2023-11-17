package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var DB *sql.DB

func Start() error {
	if ld := godotenv.Load(); ld != nil{
		log.Fatal(ld)
	}
	var err error
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		os.Getenv("HOST"), os.Getenv("PORT"), os.Getenv("USER"), os.Getenv("PASSWORD"), os.Getenv("DBNAME"))
	DB, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		return err
	}
	createUserTable := `
        CREATE TABLE IF NOT EXISTS users (
            id SERIAL PRIMARY KEY,
            username VARCHAR(50) UNIQUE NOT NULL,
            email VARCHAR(100) UNIQUE NOT NULL,
			password VARCHAR(255) NOT NULL,
			pfpurl VARCHAR(100),
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        );
    `
	createFileTable := `
        CREATE TABLE IF NOT EXISTS files (
            id SERIAL PRIMARY KEY,
            user_id INTEGER REFERENCES users(id) NOT NULL,
            filename VARCHAR(255) NOT NULL,
            file_path VARCHAR(255) NOT NULL,
            share_link VARCHAR(255) UNIQUE NOT NULL,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        );
    `
	_, err = DB.Exec(createUserTable)
	if err != nil {
		return err
	}

	_, err = DB.Exec(createFileTable)
	if err != nil {
		return err
	}
	return nil
}
