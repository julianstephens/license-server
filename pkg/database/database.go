package database

import (
	"fmt"

	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DB    *gorm.DB
	err   error
	DBErr error
)

type Database struct {
	*gorm.DB
}

func Setup() error {
	var db *gorm.DB
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s", viper.GetString("postgres_host"), viper.GetString("postgres_user"), viper.GetString("postgres_password"), viper.GetString("postgres_db"), "5432")

	db, err = gorm.Open(postgres.Open(dsn))
	if err != nil {
		DBErr = err
		return err
	}

	DB = db

	return nil
}

func GetDB() *gorm.DB {
	return DB
}

func GetDBErr() error {
	return DBErr
}
