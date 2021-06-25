package main

import (
	redisDB "InstagramStatistic/internal/Database"
	insta "InstagramStatistic/internal/Insta"
	telegram "InstagramStatistic/internal/Telegram"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var config Config
var telegramMessageChannel chan tgbotapi.Message
var webServerChannel chan string

func main() {
	go startWebServer()
	config = initConfig()
	telegramMessageChannel = make(chan tgbotapi.Message)

	if err := redisDB.Init(); err == nil {
		if err := insta.Init(config.INSTAGRAM_USERNAME, config.INSTAGRAM_PASSWORD, config.PROXY_URL, config.PROXY_LOGIN, config.PROXY_PASSWORD, config.REQUEST_DELAY_MIN, config.REQUEST_DELAY_MAX); err == nil {
			go telegram.Init(config.TELEGRAM_TOKEN, telegramMessageChannel)

			for telegramMessage := range telegramMessageChannel {
				if username := redisDB.Get(strconv.Itoa(telegramMessage.From.ID) + "_username"); username != "" || strings.HasPrefix(telegramMessage.Text, "/") {
					if username != "" {
						fmt.Println("Привязанный аккаунт: " + username + " TelegramID: " + strconv.Itoa(telegramMessage.From.ID))
					}

					if strings.HasPrefix(telegramMessage.Text, "/") {
						if telegram.ExecuteCommand(&username, telegramMessage) {
							continue
						}
						if username == "" {
							telegram.SendMessage("Введи имя аккаунта:")
						} else {
							telegram.SendMessage("Указанной команды не существует: " + strings.Fields(telegramMessage.Text)[0])
						}
					} else {
						telegram.SendMessage("Для получения справки введи команду: /help")
					}
				} else {
					if err := insta.GetUserInfo(telegramMessage.Text); err == nil {
						redisDB.Set(strconv.Itoa(telegramMessage.From.ID)+"_username", telegramMessage.Text)
						telegram.SendMessage("Привязано новое имя аккаунта: " + telegramMessage.Text)
					} else {
						telegram.SendMessage("Пользователь " + telegramMessage.Text + " не найден.\nВведи имя аккаунта:")
					}
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
