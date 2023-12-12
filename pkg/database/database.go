package database

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
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

	newDBLogger := logger.New(
		log.New(getWriter(), "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level (Silent, Error, Warn, Info)
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,       // Disable color
		},
	)

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		// Logger: logger.Default.LogMode(logger.Info),
		Logger: newDBLogger,
	})
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

func getWriter() io.Writer {
	file, err := os.OpenFile("ls.db.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return os.Stdout
	} else {
		return file
	}
}
