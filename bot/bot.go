package bot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

type ConfigBot struct {
	Token string
	Debug bool
}

func NewBot(conf *ConfigBot) *tgbotapi.BotAPI {
	fmt.Printf("values: %t, %s \n", conf.Debug, conf.Token)

	bot, err := tgbotapi.NewBotAPI(conf.Token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = conf.Debug

	log.Printf("Authorized on account %s \n", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if len(update.Message.Photo) != 0 {
			go HandleImage(bot, &update)
			continue
		}

		if update.Message.Text != "" {
			go HandleMessage(bot, &update)
		}
	}

	return bot
}