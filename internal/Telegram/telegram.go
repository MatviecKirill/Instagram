package Telegram

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"os"
	"strconv"
	"time"
)

var bot *tgbotapi.BotAPI
var update tgbotapi.Update

func Init(token string, messageChannel chan tgbotapi.Message) {
	if b, err := tgbotapi.NewBotAPI(token); err == nil {
		bot = b
	} else {
		return
		fmt.Println(err)
	}

	//Телеграм бот запускается с использованием webhook'a
	//Когда кто-то напишет боту, он обратится к моему серверу
	//heroku labs:enable runtime-dyno-metadata -a <app name> - включить метаданные для приложения. HEROKU_APP_NAME
	setWebhook()
	updates := bot.ListenForWebhook("/" + bot.Token)
	for update = range updates {
		if update.Message == nil {
			continue
		}
		fmt.Printf("Message from [%s]: %s\n", update.Message.From.FirstName, update.Message.Text)
		messageChannel <- *update.Message
	}
}

func setWebhook() {
	for attemptCount := 0; attemptCount < 10; attemptCount++ {
		if _, err := bot.SetWebhook(tgbotapi.NewWebhook("https://" + os.Getenv("HEROKU_APP_NAME") + ".herokuapp.com/" + bot.Token)); err == nil {
			if info, err := bot.GetWebhookInfo(); err == nil {
				if info.LastErrorDate == 0 {
					return
				} else {
					fmt.Printf("Telegram callback failed: %s\n", info.LastErrorMessage)
					fmt.Println("Attempt count: " + strconv.Itoa(attemptCount+1))
				}
			} else {
				fmt.Println(err)
			}
		} else {
			fmt.Println(err)
		}
		time.Sleep(time.Second * 1)
		bot.SetWebhook(tgbotapi.NewWebhook("http://" + os.Getenv("HEROKU_APP_NAME") + ".herokuapp.com/" + bot.Token))
		time.Sleep(time.Second * 1)
	}
}

func SendMessage(message string, chatId int64) {
	msg := tgbotapi.NewMessage(chatId, message)

	bot.Send(msg)
}
