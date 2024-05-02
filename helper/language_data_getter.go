package helper

import (
	"fmt"
	"log"
	"trending2telbot/config"
	"trending2telbot/model"

	"github.com/PuerkitoBio/goquery"
)

func GetLanguageData(lang string) []model.Message {
	results, err := scrapeLanguageData(lang)
	if err != nil {
		log.Printf("Error getting language data for %s: %v", lang, err)
		return nil
	}
	return results
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

		results = append(results, model.Message{
			Title:       title,
			Language:    lang,
			Description: description,
			URL:         link,
		})
	})
	return results, nil
}
