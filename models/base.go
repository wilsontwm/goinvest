package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"    // To install dependencies first
	_ "github.com/jinzhu/gorm/dialects/postgres" // To install dependencies first
	"github.com/joho/godotenv"
	"github.com/satori/go.uuid"
	_ "goinvest/news" // To install dependencies first
	"log"
	"os"
	"time"
)

var db *gorm.DB // database
var dbURI string
var dbDriver string
var jwtKey string

// Base contains common columns for all tables.
type Base struct {
	ID        uuid.UUID `gorm:"type:varchar(255);primary_key;"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

// BeforeCreate will set a UUID rather than numeric ID.
func (base *Base) BeforeCreate(scope *gorm.Scope) error {
	uuid := uuid.NewV4()
	return scope.SetColumn("ID", uuid)
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Print(err)
	}

	dbDriver = os.Getenv("DB_DRIVER")
	dbURI = os.Getenv("DB_URI_LOCAL")
	if os.Getenv("IS_PRODUCTION") == "True" {
		dbURI = os.Getenv("DB_URI_PRODUCTION")
	}
	jwtKey = os.Getenv("JWT_KEY")

	migrateDatabase()
}

// Datebase migration
func migrateDatabase() {
	db := GetDB()

	db.Debug().AutoMigrate(
		&User{}, &Account{},
	)

	// Migration scripts
	//db.Model(&Attendee{}).AddForeignKey("parent_id", "attendees(id)", "SET NULL", "RESTRICT")
}

// GetDB : To get the DB instance
func GetDB() *gorm.DB {
	// Making connection to the database
	db, err := gorm.Open(dbDriver, dbURI)
	if err != nil {
		log.Println(err)
	}

	retryCount := 30
	for {
		err := db.DB().Ping()
		if err != nil {
			if retryCount == 0 {
				log.Fatalf("Not able to establish connection to database")
			}

			log.Printf(fmt.Sprintf("Could not connect to database. Wait 2 seconds. %d retries left...", retryCount))
			retryCount--
			time.Sleep(2 * time.Second)
		} else {
			break
		}
	}

	return db
}
