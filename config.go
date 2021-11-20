package main

import (
	"ascii-art-bot/bot"
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	Bot bot.ConfigBot
}

func NewConfig(path string) *Config {
	viper.SetConfigFile("config.yaml")
	viper.AddConfigPath(path)

	if err := viper.ReadInConfig(); err != nil {
		log.Panicf("Fatal error reading config file: %w \n", err)
	}

	var config Config

	if err := viper.Unmarshal(&config); err != nil {
		log.Panicf("Fatal error unmarshalling config file: %w \n", err)
	}

	return &config
}