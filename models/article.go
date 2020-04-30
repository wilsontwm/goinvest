package models

import (
	"github.com/jinzhu/gorm"
	"goinvest/news"
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

type ArticleListFilter struct {
	NewSources        []int
	Limit             int
	PreviousID        int
	PreviousPublished time.Time
}

var NewsSources map[int]string

func init() {
	NewsSources = map[int]string{
		ChinaPress:      "China Press (中国报)",
		NanYang:         "Nanyang Siang Pau (南洋商报)",
		SinChew:         "Sin Chew Daily (星洲日报)",
		TheStar:         "The Star",
		TheEdge:         "The Edge",
		Investing:       "Investing.com",
		NewStraitsTimes: "New Straits Times",
		MalayMail:       "Malay Mail",
		BusinessInsider: "Business Insider",
	}
}

// Get list of articles
func GetArticleList(filter ArticleListFilter) ([]Article, error) {
	var articles []Article

	db := GetDB()
	defer db.Close()

	// Get the list of articles
	db.Table("articles").
		Scopes(filterByNewSources(filter.NewSources), filterByPrev(filter.PreviousID, filter.PreviousPublished)).
		Order("published_at desc, id desc").
		Limit(filter.Limit).
		Find(&articles)

	return articles, nil
}

// Filter by the previous ID
func filterByPrev(prevID int, prevPublished time.Time) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if prevID > 0 && !prevPublished.IsZero() {
			return db.Where("published_at < ? OR (published_at = ? AND id < ?)", prevPublished, prevPublished, prevID)
		}

		return db
	}
}

// Filter the article by the news sources
func filterByNewSources(sourceID []int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if len(sourceID) > 0 {
			return db.Where("source IN (?)", sourceID)
		}

		return db
	}
}

// Crawl the articles based on the news sources
func CrawlArticles(sourceIDs []int) bool {
	var filterSources []int
	for _, v := range sourceIDs {
		if _, found := NewsSources[v]; found {
			filterSources = append(filterSources, v)
		}
	}

	// Create the jobs and worker pools
	var numOfJobs = len(filterSources)
	var numOfWorkers = 4
	var outputs []bool
	if numOfJobs > 0 {
		sourceIDJobs := make(chan int, numOfJobs) // Accept the news source ID
		results := make(chan bool, numOfJobs)     // Use bool to return if it's run successfully

		for w := 1; w <= numOfWorkers; w++ {
			go crawlArticlesWorker(sourceIDJobs, results)
		}

		for _, s := range filterSources {
			sourceIDJobs <- s
		}

		close(sourceIDJobs)

		for a := 1; a <= numOfJobs; a++ {
			output := <-results
			if output {
				outputs = append(outputs, output)
			}
		}

	}

	return len(outputs) == numOfJobs
}

// Worker to crawl the article based on the source ID
func crawlArticlesWorker(sourceIDs <-chan int, results chan<- bool) {
	for sourceID := range sourceIDs {
		switch sourceID {
		case ChinaPress:
			results <- news.CrawlChinaPress()
		case NanYang:
			results <- news.CrawlNanYang()
		case SinChew:
			results <- news.CrawlSinChew()
		case TheStar:
			results <- news.CrawlTheStar()
		case TheEdge:
			results <- news.CrawlTheEdge()
		case Investing:
			results <- news.CrawlInvesting()
		case NewStraitsTimes:
			results <- news.CrawlNewStraitsTimes()
		case MalayMail:
			results <- news.CrawlMalayMail()
		case BusinessInsider:
			results <- news.CrawlBusinessInsider()
		default:
			results <- false
		}
	}
}
