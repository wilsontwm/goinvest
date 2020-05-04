package news

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
	"log"
	"os"
	"time"
)

const (
	ChinaPress      = 1
	NanYang         = 2
	SinChew         = 3
	TheStar         = 4
	TheEdge         = 5
	Investing       = 6
	NewStraitsTimes = 7
	MalayMail       = 8
	BusinessInsider = 9
)

type Article struct {
	ID          int    `gorm:"primary_key;"`
	Title       string `sql:"type:longtext"`
	Content     string `sql:"type:longtext"`
	Source      int
	URL         string
	Thumbnail   string
	PublishedAt time.Time
	CreatedAt   time.Time
}

var db *gorm.DB // database
var dbURI string
var dbDriver string

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

	db.Debug().AutoMigrate(&Article{})

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

// Create new article
func CreateArticle(input *Article) (*Article, error) {
	article := input

	db := GetDB()
	defer db.Close()

	// Create the user
	db.Create(article)

	if article.ID <= 0 {
		return nil, fmt.Errorf("Article is not created.")
	}

	return article, nil
}

// Get article links by new source
func GetArticlesBySource(sourceID int) []string {
	var links []string

	db := GetDB()
	defer db.Close()

	db.Table("articles").Where("source = ?", sourceID).Pluck("url", &links)

	return links
}
