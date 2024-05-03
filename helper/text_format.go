package helper

import (
	"fmt"
	"strings"
	"time"
	"trending2telbot/config"
	"trending2telbot/model"
)

func CleanText(text string) string {
	return strings.Join(strings.Fields(text), " ")
}

func FormatDate2TgMessage(msg model.Message) string {
	translated := TranslateText(config.TranslateURL, msg.Description, "en", "zh")
	now := time.Now().Format("20060102")
	return fmt.Sprintf("📌*%s*\n%s```\n%s\n```\n#D%s  #%s  🔗[RepoURL](%s)  ⭐️%s", msg.Title, msg.Description, translated, now, msg.Language, msg.URL, msg.TodayStars)
}
