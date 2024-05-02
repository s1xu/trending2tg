package config

import (
	"net/http"
	"os"
	"strings"
	"time"
)

var (
	BotToken     string
	ChatID       string
	TranslateURL string
	Languages    []string
)

var Client = &http.Client{
	Timeout: time.Second * 10,
}

func init() {
	BotToken = os.Getenv("BOT_TOKEN")
	ChatID = os.Getenv("CHAT_ID")
	TranslateURL = os.Getenv("TRANSLATE_URL")
	language := os.Getenv("LANGUAGES")
	Languages = strings.Split(language, ",")
}
