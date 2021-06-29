package Telegram

import (
	redisDB "InstagramStatistic/internal/Database"
	insta "InstagramStatistic/internal/Insta"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"strconv"
	"strings"
)

func ExecuteCommand(username *string, telegramMessage tgbotapi.Message) bool {
	if strings.HasPrefix(telegramMessage.Text, "/help") {
		tgMessage := ""
		if *username != "" {
			tgMessage += "✅️ Привязанный аккаунт: " + *username + "\n\n"
		}
		tgMessage += "▫️ Анализ взаимных подписок. Команда:\n /nonmutual имя пользователя\n"
		tgMessage += "▫️ Анализ отписавшихся пользователей. Команда:\n /unsubscribe имя пользователя\n"
		tgMessage += "▫️ Привязать новый аккаунт. Команда:\n /account имя пользователя\n"
		tgMessage += "▫️ Отвзяать аккаунт. Команда:\n /accountunbind\n"
		SendMessage(tgMessage)
		return true
	}

	if strings.HasPrefix(telegramMessage.Text, "/accountunbind") {
		if *username != "" {
			redisDB.Del(strconv.Itoa(telegramMessage.From.ID)+"_username")
			SendMessage("Аккаунт отвязан.")
		} else {
			SendMessage("Привязанный аккаунт не найден.")
		}
		return true
	}

	if strings.HasPrefix(telegramMessage.Text, "/account") {
		*username = strings.Trim(strings.TrimPrefix(telegramMessage.Text, "/account"), " ")
		if err := insta.GetUserInfo(*username); err == nil {
			redisDB.Set(strconv.Itoa(telegramMessage.From.ID)+"_username", *username)
			SendMessage("Привязано новое имя аккаунта: " + *username)
		} else {
			SendMessage("Пользователь " + *username + " не найден.")
		}
		return true
	}

	if strings.HasPrefix(telegramMessage.Text, "/nonmutual") {
		usernameFromCommand := strings.Trim(strings.TrimPrefix(telegramMessage.Text, "/nonmutual"), " ")
		if usernameFromCommand != "" {
			*username = usernameFromCommand
		}
		if *username == "" {
			return false
		}
		SendMessage("Собираю данные по пользователю: " + *username + ". Ожидайте...")
		if message, err := insta.GetNonMutualFollowersMessage(*username); err == nil {
			SendMessage(message)
			fmt.Print(message)
		} else {
			fmt.Println(err)
		}
		return true
	}

	if strings.HasPrefix(telegramMessage.Text, "/unsubscribe") {
		usernameFromCommand := strings.Trim(strings.TrimPrefix(telegramMessage.Text, "/unsubscribe"), " ")
		if usernameFromCommand != "" {
			*username = usernameFromCommand
		}
		if *username == "" {
			return false
		}
		SendMessage("Собираю данные по пользователю: " + *username + ". Ожидайте...")
		if message, err := insta.GetUnsubscribedFollowersMessage(*username); err == nil {
			SendMessage(message)
			fmt.Println(message)
		} else {
			fmt.Println(err)
		}
		return true
	}
	return false
}
