package main

import (
	redisDB "InstagramStatistic/internal/Database"
	insta "InstagramStatistic/internal/Insta"
	telegram "InstagramStatistic/internal/Telegram"
	"fmt"
	"net/http"
	"os"
	"strings"
)

var config Config
var telegramMessageChannel, webServerChannel chan string
var telegramMessage string

func main() {
	go startWebServer()
	config = initConfig()
	telegramMessageChannel = make(chan string)

	if err := redisDB.Init(); err == nil {
		if err := insta.Init(config.INSTAGRAM_USERNAME, config.INSTAGRAM_PASSWORD, config.PROXY_URL, config.PROXY_LOGIN, config.PROXY_PASSWORD, config.REQUEST_DELAY_MIN, config.REQUEST_DELAY_MAX); err == nil {
			go telegram.Init(config.TELEGRAM_TOKEN, telegramMessageChannel)

			for {
				telegramMessage = <-telegramMessageChannel

				if strings.HasPrefix(telegramMessage, "/") {
					if strings.HasPrefix(telegramMessage, "/взаимные") {
						username := strings.Trim(strings.TrimPrefix(telegramMessage, "/взаимные"), " ")
						telegram.SendMessage("Собираю данные по пользователю: " + username + ". Ожидайте...")
						if message, err := insta.GetNonMutualFollowersMessage(username); err == nil {
							telegram.SendMessage(message)
							fmt.Print(message)
						} else {
							fmt.Println(err)
						}
					}

					if strings.HasPrefix(telegramMessage, "/отписались") {
						username := strings.Trim(strings.TrimPrefix(telegramMessage, "/отписались"), " ")
						telegram.SendMessage("Собираю данные по пользователю: " + username + ". Ожидайте...")
						if message, err := insta.GetUnsubscribedFollowersMessage(username); err == nil {
							telegram.SendMessage(message)
							fmt.Println(message)
						} else {
							fmt.Println(err)
						}
					}
				} else {
					tgMessage := "▫️ Анализ взаимных подписок. Команда:\n /взаимные имя пользователя\n"
					tgMessage += "▫️ Анализ отписавшихся пользователей. Команда:\n /отписались имя пользователя\n"
					telegram.SendMessage(tgMessage)
				}
			}
		} else {
			fmt.Println(err)
		}
	} else {
		fmt.Println(err)
	}
}

//Запуск веб сервера, для ответов на запросы.
//Необходим, чтобы телеграм бот получил ответ при запуске веб хука. И чтобы heroku не падал из-за отсутсвия обработчика запросов.
func startWebServer() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":"+os.Getenv("PORT"), nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%v", "InstagramStatistic")
}