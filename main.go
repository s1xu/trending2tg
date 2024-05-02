package main

import (
	"time"
	"trending2telbot/config"
	"trending2telbot/helper"
)

func main() {
	for _, lang := range config.Languages {
		processLanguage(lang)
	}
}

func processLanguage(lang string) {
	results := helper.GetLanguageData(lang)
	for _, result := range results {
		message := helper.FormatDate2TgMessage(result)
		helper.SendMessage2Telegram(message)
		time.Sleep(3 * time.Second) // limits to 20 per minute
	}
}
