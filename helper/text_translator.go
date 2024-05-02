package helper

import (
	"bytes"
	"encoding/json"
	"log"
	"trending2telbot/config"
	"trending2telbot/model"
)

func TranslateText(url, text, sourceLang, targetLang string) string {
	translateURL := url
	payload := model.Translate{
		Text:       text,
		SourceLang: sourceLang,
		TargetLang: targetLang,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		log.Println("Error marshalling payload:", err)
		return ""
	}

	resp, err := config.Client.Post(translateURL, "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		log.Println("Error making request:", err)
		return ""
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Println("Error decoding response:", err)
		return ""
	}

	translatedText, ok := result["data"].(string)
	if !ok {
		log.Println("Invalid response format")
		return ""
	}

	return translatedText
}
