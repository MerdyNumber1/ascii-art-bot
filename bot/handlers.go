package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

func HandleImage(bot *tgbotapi.BotAPI, update *tgbotapi.Update)  {
	log.Printf(bot.GetFileDirectURL(update.Message.Photo[0].FileID))
}

func HandleMessage(bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
	msg.ReplyToMessageID = update.Message.MessageID

	if _, err := bot.Send(msg); err != nil {
		panic(err)
	}
}
