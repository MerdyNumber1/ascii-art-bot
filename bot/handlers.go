package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

func HandleImage(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	url, err := bot.GetFileDirectURL(update.Message.Photo[len(update.Message.Photo) - 1].FileID)
	if err != nil {
		HandleError(bot, update)
		return
	}

	img, _, width, height, err := DownloadImage(url)
	if err != nil {
		HandleError(bot, update)
		panic(err)
		return
	}

	text := ConvertImageToAscii(img)
	imgBytes, err := GenerateImageFromText(text, "#ece7ea", "#350b23", 16, width * 2, height * 2)
	if err != nil {
		HandleError(bot, update)
		panic(err)
		return
	}

	msg := tgbotapi.NewPhoto(update.Message.Chat.ID, tgbotapi.FileBytes{Bytes: imgBytes, Name: "1"})
	if _, err := bot.Send(msg); err != nil {
		log.Println(err)
	}
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
