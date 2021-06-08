package main

import (
	stat "InstagramStatistic/internal/Stat"
	telegram "InstagramStatistic/internal/Telegram"
	"fmt"
	"strings"
)

var config Config
var messageChannel chan string
var telegramMessage string

func main() {
	config = initConfig()
	messageChannel = make(chan string)

	if err := stat.Init(config.INSTAGRAM_USERNAME, config.INSTAGRAM_PASSWORD, config.PROXY_URL, config.PROXY_LOGIN, config.PROXY_PASSWORD, config.REQUEST_DELAY_MIN, config.REQUEST_DELAY_MAX); err == nil {
		go func() {
			telegram.Init(config.TELEGRAM_TOKEN, messageChannel)
		}()

		for {
			telegramMessage = <-messageChannel

			if strings.HasPrefix(telegramMessage, "/get") {
				username := strings.Trim(strings.TrimPrefix(telegramMessage, "/get"), " ")
				telegram.SendMessage("Собираю данные по пользователю: " + username + ". Ожидайте...")
				if message, err := stat.GetNonMutualFollowersMessage(username); err == nil {
					telegram.SendMessage(message)
					fmt.Print(message)
				} else {
					fmt.Println(err)
				}
			} else {
				telegram.SendMessage("Чтобы начать анализ введи команду:\n /get имя пользователя")
			}
		}
	} else {
		fmt.Println(err)
	}
}
