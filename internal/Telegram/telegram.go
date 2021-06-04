package Telegram

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
)

var bot *tgbotapi.BotAPI
var update tgbotapi.Update

func Init(token string, messageChannel chan string) {
	if b, err := tgbotapi.NewBotAPI(token); err == nil {
		bot = b
	}else {
		fmt.Println(err)
	}
	fmt.Println("Telegram bot init successfully")

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	if updateChannel, err := bot.GetUpdatesChan(u); err == nil {
		for update = range updateChannel {
			if update.Message == nil {
				continue
			}

			log.Printf("Message from [%s]: %s", update.Message.From.FirstName, update.Message.Text)
			messageChannel <- update.Message.Text
		}
	} else {
		fmt.Println(err)
	}
}

func SendMessage(message string) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, message)

	bot.Send(msg)
}
