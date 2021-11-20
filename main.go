package main

import (
	"ascii-art-bot/bot"
)

func main() {
	conf := NewConfig("")
	bot.NewBot(&conf.Bot)
}
