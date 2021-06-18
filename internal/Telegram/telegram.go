package Telegram

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"os"
)

var bot *tgbotapi.BotAPI
var update tgbotapi.Update

func Init(token string, messageChannel chan string) {
	if b, err := tgbotapi.NewBotAPI(token); err == nil {
		bot = b
	} else {
		return
		fmt.Println(err)
	}

	//Телеграм бот запускается с использованием webhook'a
	//Когда кто-то напишет боту, он обратится к моему серверу
	if _, err := bot.SetWebhook(tgbotapi.NewWebhook("https://" + os.Getenv("HEROKU_APP_NAME") + ".herokuapp.com/" + bot.Token)); err == nil {
		if info, err := bot.GetWebhookInfo(); err == nil {
			if info.LastErrorDate == 0 {
				fmt.Println("Telegram bot init successfully")

				updates := bot.ListenForWebhook("/" + bot.Token)
				for update = range updates {
					if update.Message == nil {
						continue
					}

					fmt.Printf("Message from [%s]: %s\n", update.Message.From.FirstName, update.Message.Text)
					messageChannel <- update.Message.Text
				}

			} else {
				fmt.Printf("Telegram callback failed: %s", info.LastErrorMessage)
			}
		} else {
			fmt.Println(err)
		}
	} else {
		fmt.Println(err)
	}
}

func SendMessage(message string) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, message)

	bot.Send(msg)
}
