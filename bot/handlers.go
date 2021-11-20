package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

func HandleImage(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	url, err := bot.GetFileDirectURL(update.Message.Photo[len(update.Message.Photo) - 1].FileID)
	if err != nil {
		HandleError(bot, update)
	}

	img, _, err := DownloadImage(url)
	if err != nil {
		HandleError(bot, update)
	}

	ConvertImageToAscii(img)
}

func HandleMessage(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)

	if _, err := bot.Send(msg); err != nil {
		panic(err)
	}
}

func HandleError(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Something went wrong")
	if _, err := bot.Send(msg); err != nil {
		log.Println(err)
	}
}
