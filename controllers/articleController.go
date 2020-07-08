package controllers

import (
	"encoding/json"
	"goinvest/models"
	"goinvest/utils"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// ArticleCrawlInput : The input for the articles
type ArticleCrawlInput struct {
	Sources []int
}

// ArticleList (GET) : List of articles
var ArticleList = func(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]interface{})

	newsSources := []int{}
	if sourcesParam := utils.GetParam(r, "sources"); sourcesParam != "" {
		sources := strings.Split(sourcesParam, ",")
		for _, source := range sources {
			if sourceID, err := strconv.Atoi(source); err == nil && sourceID > 0 {
				// Check if the source ID is within the list of source
				if _, found := models.NewsSources[sourceID]; found {
					newsSources = append(newsSources, sourceID)
				}
			}
		}
	}

	limit := 10
	if limitParam := utils.GetParam(r, "limit"); limitParam != "" {
		limit, _ = strconv.Atoi(limitParam)
	}

	previousID := 0
	if previousIDParam := utils.GetParam(r, "prev"); previousIDParam != "" {
		previousID, _ = strconv.Atoi(previousIDParam)
	}

	previousPublished := time.Time{}
	if previousPublishedParam := utils.GetParam(r, "prevPublished"); previousPublishedParam != "" {
		if previousPublishedUnix, err := strconv.Atoi(previousPublishedParam); err == nil && previousPublishedUnix > 0 {
			previousPublished = time.Unix(int64(previousPublishedUnix), 0)
		}

	}

	// Build the filter for the article list
	filter := models.ArticleListFilter{
		NewSources:        newsSources,
		Limit:             limit,
		PreviousID:        previousID,
		PreviousPublished: previousPublished,
	}

	// Get the articles list
	articles, err := models.GetArticleList(filter)

	if err != nil {
		utils.Fail(w, http.StatusBadRequest, resp, err.Error())
		return
	}

	utils.Success(w, http.StatusOK, resp, articles, "")
}

// ArticleCrawl (POST) : Crawl articles from various sources
var ArticleCrawl = func(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]interface{})

	input := ArticleCrawlInput{}
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		utils.Fail(w, http.StatusInternalServerError, resp, err.Error())
		return
	}

	if ok := models.CrawlArticles(input.Sources); !ok {
		utils.Fail(w, http.StatusInternalServerError, resp, "Not all the news sources have been crawled successfully.")
		return
	}

	utils.Success(w, http.StatusOK, resp, "", "You have crawled all the articles successfully.")
}
