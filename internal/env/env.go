package env

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

var Global Env

// Database ...
type Database struct {
	Driver   string
	HostName string
	User     string
	Password string
	Name     string
	Schema   string

	URL string
}

// Env ...
type Env struct {
	DB Database
}

func init() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}
	Global.DB = Database{
		Driver:   os.Getenv("DB_DRIVER"),
		HostName: os.Getenv("DB_HOST_NAME"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Name:     os.Getenv("DB_NAME"),
		Schema:   os.Getenv("DB_SCHEMA"),
	}
	Global.DB.URL = fmt.Sprintf(
		"postgres://%s:%s@%s:5432/%s",
		Global.DB.User, Global.DB.Password, Global.DB.HostName, Global.DB.Name)
}
