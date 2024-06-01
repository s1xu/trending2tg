package helper

import (
	"fmt"
	"log"
	"strings"
	"time"
	"trending2telbot/config"
	"trending2telbot/model"

	"github.com/PuerkitoBio/goquery"
)

func GetLanguageData(lang string) ([]model.Message, error) {
	var results []model.Message
	var err error

	// Retry a maximum of three times
	maxRetries := 3
	// Retry latency
	retryInterval := time.Second * 2

	for attempt := 1; attempt <= maxRetries; attempt++ {
		results, err = scrapeLanguageData(lang)
		if err == nil {
			return results, nil
		}

		log.Printf("Attempt %d: Error getting language data for %s: %v", attempt, lang, err)
		time.Sleep(retryInterval)
	}

	return nil, fmt.Errorf("after %d attempts, last error: %w", maxRetries, err)
}

func scrapeLanguageData(lang string) ([]model.Message, error) {
	url := fmt.Sprintf("https://github.com/trending/%s", lang)
	resp, err := config.Client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	var results []model.Message
	doc.Find("article.Box-row").Each(func(i int, s *goquery.Selection) {
		title := CleanText(s.Find(".lh-condensed a").Text())
		description := CleanText(s.Find("p.col-9").Text())
		link, _ := s.Find(".lh-condensed a").Attr("href")
		link = "https://github.com" + link
		todayStars := CleanText(s.Find(".d-inline-block.float-sm-right").Text())
		todayStars = strings.Split(todayStars, " ")[0]
		if lang == "" {
			lang = "all"
		}
		results = append(results, model.Message{
			Title:       title,
			Language:    lang,
			Description: description,
			URL:         link,
			TodayStars:  todayStars,
		})
	})
	log.Printf("Scraped %d results for language: %s", len(results), lang)
	return results, nil
}
