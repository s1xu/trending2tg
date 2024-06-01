package main

import (
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"
	"trending2telbot/config"
	"trending2telbot/helper"
	"trending2telbot/model"
)

func main() {
	startTime := time.Now()

	log.Println("Starting the application...")
	if len(config.Languages) == 0 {
		log.Fatalf("No languages configured. Exiting.")
	}

	db, err := helper.InitializeDatabase("trends.db")
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}
	defer db.Close()

	var wg sync.WaitGroup
	messageCounts := make(map[string]int)
	resultsChan := make(chan map[string][]model.Message)
	messagesChan := make(chan model.Message)

	wg.Add(1)
	go func() {
		defer wg.Done()
		fetchDataAndProcess(config.Languages, resultsChan)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		processResultsAndFormat(resultsChan, messagesChan, db)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		sendMessages(messagesChan, messageCounts)
	}()

	wg.Wait()

	elapsed := time.Since(startTime)
	sendSummary(messageCounts, elapsed)

	log.Println("Application finished.")
}

func fetchDataAndProcess(languages []string, resultsChan chan map[string][]model.Message) {
	defer close(resultsChan)

	for _, lang := range languages {
		results, err := helper.GetLanguageData(lang)
		if err != nil {
			log.Printf("Error fetching data for language %s: %v", lang, err)
			continue
		}
		resultsChan <- map[string][]model.Message{lang: results}
	}
}

func processResultsAndFormat(resultsChan chan map[string][]model.Message, messagesChan chan model.Message, db *sql.DB) {
	defer close(messagesChan)

	for combinedResults := range resultsChan {
		for lang, messages := range combinedResults {
			processLanguageResults(lang, messages, db, messagesChan)
		}
	}
}

func processLanguageResults(lang string, messages []model.Message, db *sql.DB, messagesChan chan model.Message) {
	for _, result := range messages {
		processSingleResult(result, lang, db, messagesChan)
	}
}

func processSingleResult(result model.Message, lang string, db *sql.DB, messagesChan chan model.Message) {
	trend := model.NewTrends(result.Title, result.Language, result.URL)
	if inserted, err := helper.InsertIfNotExists(db, trend); err != nil {
		log.Printf("Error inserting trend: %v", err)
		return
	} else if !inserted {
		return
	}

	log.Printf("Inserted trend: %s", result.Title)
	message := helper.FormatDate2TgMessage(result)
	messagesChan <- model.Message{Content: message, Language: lang}
}

func sendMessages(messagesChan chan model.Message, messageCounts map[string]int) {
	limit := time.Tick(3 * time.Second)
	for message := range messagesChan {
		<-limit // limits to 20 per minute
		helper.SendMessage2Telegram(message.Content)
		messageCounts[message.Language]++
	}
}

func sendSummary(messageCounts map[string]int, elapsed time.Duration) {
	currentDate := time.Now().Format("2006-01-02")
	summaryMessage := fmt.Sprintf("Today's List Updated: %s 🎉\n", currentDate)

	for lang, count := range messageCounts {
		summaryMessage += fmt.Sprintf("> #%s: %d News\n", lang, count)
	}

	elapsedInSeconds := int(elapsed.Seconds())
	summaryMessage += fmt.Sprintf("Time Elapsed: %ds 🚀", elapsedInSeconds)
	helper.SendMessage2Telegram(summaryMessage)
}
