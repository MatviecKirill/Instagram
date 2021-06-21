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

			for message := range telegramMessageChannel {
				if username := redisDB.Get(strconv.Itoa(message.From.ID) + "_username"); username != "" || strings.HasPrefix(message.Text, "/") {
					if username != "" {
						fmt.Println("Привязанный аккаунт: " + username + " TelegramID: " + strconv.Itoa(message.From.ID))
					}

					if strings.HasPrefix(message.Text, "/") {
						if strings.HasPrefix(message.Text, "/справка") {
							tgMessage := ""
							if username != "" {
								tgMessage += "✅️ Привязанный аккаунт: " + username + "\n"
							}
							tgMessage += "▫️ Анализ взаимных подписок. Команда:\n /взаимные имя пользователя\n"
							tgMessage += "▫️ Анализ отписавшихся пользователей. Команда:\n /отписались имя пользователя\n"
							tgMessage += "▫️ Привязать новый аккаунт. Команда:\n /аккаунт имя пользователя\n"
							telegram.SendMessage(tgMessage)
							break
						}

						if strings.HasPrefix(message.Text, "/аккаунт") {
							username = strings.Trim(strings.TrimPrefix(message.Text, "/аккаунт"), " ")
							if err := insta.GetUserInfo(username); err == nil {
								redisDB.Set(strconv.Itoa(message.From.ID)+"_username", username)
								telegram.SendMessage("Привязано новое имя аккаунта: " + username)
							} else {
								telegram.SendMessage("Пользователь " + message.Text + " не найден.")
							}
							break
						}

						if strings.HasPrefix(message.Text, "/невзаимные") {
							if username == "" {
								username = strings.Trim(strings.TrimPrefix(message.Text, "/невзаимные"), " ")
							}
							telegram.SendMessage("Собираю данные по пользователю: " + username + ". Ожидайте...")
							if message, err := insta.GetNonMutualFollowersMessage(username); err == nil {
								telegram.SendMessage(message)
								fmt.Print(message)
							} else {
								fmt.Println(err)
							}
							break
						}

						if strings.HasPrefix(message.Text, "/отписались") {
							if username == "" {
								username = strings.Trim(strings.TrimPrefix(message.Text, "/отписались"), " ")
							}
							telegram.SendMessage("Собираю данные по пользователю: " + username + ". Ожидайте...")
							if message, err := insta.GetUnsubscribedFollowersMessage(username); err == nil {
								telegram.SendMessage(message)
								fmt.Println(message)
							} else {
								fmt.Println(err)
							}
							break
						}

						telegram.SendMessage("Указанной команды не существует: " + strings.Fields(message.Text)[0])
					} else {
						telegram.SendMessage("Для получения справки введи команду: /справка")
					}
				} else {
					if err := insta.GetUserInfo(message.Text); err == nil {
						redisDB.Set(strconv.Itoa(message.From.ID)+"_username", message.Text)
						telegram.SendMessage("Привязано новое имя аккаунта: " + message.Text)
					} else {
						telegram.SendMessage("Пользователь " + message.Text + " не найден.\nВведи имя аккаунта:")
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
