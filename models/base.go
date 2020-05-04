package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
	"github.com/satori/go.uuid"
	_ "goinvest/news"
	"log"
	"os"
	"time"
)

var db *gorm.DB // database
var dbURI string
var dbDriver string

// Base contains common columns for all tables.
type Base struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;"`
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

	migrateDatabase()
}

// Datebase migration
func migrateDatabase() {
	db := GetDB()

	db.Debug().AutoMigrate()

	// Migration scripts
	//db.Model(&Attendee{}).AddForeignKey("parent_id", "attendees(id)", "SET NULL", "RESTRICT")
}

func GetDB() *gorm.DB {
	// Making connection to the database
	db, err := gorm.Open(dbDriver, dbURI)
	if err != nil {
		log.Println(err)
	}

	return db
}
