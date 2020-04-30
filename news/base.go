package news

import (
	"fmt"
	"github.com/jinzhu/gorm"
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
	ID          int `gorm:"primary_key;"`
	Title       string
	Content     string
	Source      int
	URL         string
	Thumbnail   string
	PublishedAt time.Time
	CreatedAt   time.Time
}

var db *gorm.DB // database
var username, password, dbName, dbHost, dbPort string

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Print(err)
	}

	username = os.Getenv("DB_USER")
	password = os.Getenv("DB_PASSWORD")
	dbName = os.Getenv("DB_NAME")
	dbHost = os.Getenv("DB_HOST")
	dbPort = os.Getenv("DB_PORT")
}

func GetDB() *gorm.DB {
	dbUri := fmt.Sprintf("postgres://%v@%v:%v/%v?sslmode=disable&password=%v", username, dbHost, dbPort, dbName, password)

	// Making connection to the database
	db, err := gorm.Open("postgres", dbUri)
	if err != nil {
		log.Println(err)
	}

	return db
}

// Create new article
func CreateArticle(input *Article) (*Article, error) {
	article := input

	// db := GetDB()
	// defer db.Close()

	// // Create the user
	// db.Create(article)

	// if article.ID <= 0 {
	// 	return nil, fmt.Errorf("Article is not created.")
	// }

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
