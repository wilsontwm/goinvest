package models

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"    // To install dependencies first
	_ "github.com/jinzhu/gorm/dialects/postgres" // To install dependencies first
	"github.com/joho/godotenv"
	"github.com/satori/go.uuid"
	_ "goinvest/news" // To install dependencies first
	"log"
	"os"
	"reflect"
	"time"
)

var db *gorm.DB // database
var dbURI string
var dbDriver string
var jwtKey string

// use a single instance of Validate, it caches struct info
var validate *validator.Validate

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

	validate = validator.New()
}

// Datebase migration
func migrateDatabase() {
	db := GetDB()

	db.Debug().AutoMigrate(
		&User{}, &Account{}, &FundFlow{},
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

// Get the validation message
func getValidationMessage(err error) error {

	for _, errz := range err.(validator.ValidationErrors) {
		// Build the custom errors here
		switch tag := errz.ActualTag(); tag {
		case "required":
			return fmt.Errorf(errz.StructField() + " is required.")
		case "email":
			return fmt.Errorf(errz.StructField() + " is an invalid email address.")
		case "min":
			if errz.Type().Kind() == reflect.String {
				return fmt.Errorf(errz.StructField() + " must be more than or equal to " + errz.Param() + " character(s).")
			}
			return fmt.Errorf(errz.StructField() + " must be larger than " + errz.Param() + ".")

		case "max":
			if errz.Type().Kind() == reflect.String {
				return fmt.Errorf(errz.StructField() + " must be lesser than or equal to " + errz.Param() + " character(s).")
			}
			return fmt.Errorf(errz.StructField() + " must be smaller than " + errz.Param() + ".")

		default:
			return fmt.Errorf(errz.StructField() + " is invalid.")
		}
	}

	return nil
}
