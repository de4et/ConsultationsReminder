package main

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TGBot struct {
	b *tgbotapi.BotAPI
}

func CreateBot(key string) *TGBot {
	bot, err := tgbotapi.NewBotAPI(key)
	if err != nil {
		panic(err)
	}
	return &TGBot{b: bot}
}

func (b *TGBot) SendNotification(text string) {
	msg := tgbotapi.NewMessage(cfg.ChannelChatId, text)
	if _, err := b.b.Send(msg); err != nil {
		log.Fatal(err)
	}
}
