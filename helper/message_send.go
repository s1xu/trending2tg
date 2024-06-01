package helper

import (
	"fmt"
	"log"
	"net/url"
	"trending2telbot/config"
)

func SendMessage2Telegram(message string) {
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", config.BotToken)
	queryParams := url.Values{
		"chat_id":    {config.ChatID},
		"text":       {message},
		"parse_mode": {"markdown"},
	}.Encode()

	fullURL := fmt.Sprintf("%s?%s", apiURL, queryParams)
	// log.Println(fullURL)
	resp, err := config.Client.Get(fullURL)
	if err != nil {
		log.Println("Error sending message to Telegram:", err)
		return
	}
	defer resp.Body.Close()
}
