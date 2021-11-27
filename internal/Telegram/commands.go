package Telegram

import (
	redisDB "InstagramStatistic/internal/Database"
	insta "InstagramStatistic/internal/Insta"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"strconv"
	"strings"
)

func ExecuteCommand(username *string, chatId int64, telegramMessage tgbotapi.Message) bool {
	if strings.HasPrefix(telegramMessage.Text, "/help") {
		tgMessage := ""
		if *username != "" {
			tgMessage += "✅️ Привязанный аккаунт: " + *username + "\n\n"
		}
		tgMessage += "▫️ Вся статистика. Команда:\n /scan имя пользователя\n"
		tgMessage += "▫️ Привязать новый аккаунт. Команда:\n /account имя пользователя\n"
		tgMessage += "▫️ Отвзяать аккаунт. Команда:\n /accountunbind\n"
		SendMessage(tgMessage, chatId)
		return true
	}

	if strings.HasPrefix(telegramMessage.Text, "/accountunbind") {
		if *username != "" {
			redisDB.Del(strconv.Itoa(telegramMessage.From.ID)+"_username")
			SendMessage("Аккаунт отвязан.", chatId)
		} else {
			SendMessage("Привязанный аккаунт не найден.", chatId)
		}
		return true
	}

	if strings.HasPrefix(telegramMessage.Text, "/account") {
		*username = strings.Trim(strings.TrimPrefix(telegramMessage.Text, "/account"), " ")
		if err := insta.GetUserInfo(*username); err == nil {
			redisDB.Set(strconv.Itoa(telegramMessage.From.ID)+"_username", *username)
			SendMessage("Привязано новое имя аккаунта: " + *username, chatId)
		} else {
			SendMessage("Пользователь " + *username + " не найден.", chatId)
		}
		return true
	}

	if strings.HasPrefix(telegramMessage.Text, "/scan") {
		usernameFromCommand := strings.Trim(strings.TrimPrefix(telegramMessage.Text, "/scan"), " ")
		if usernameFromCommand != "" {
			*username = usernameFromCommand
		}
		if *username == "" {
			return false
		}
		SendMessage("Собираю данные по пользователю: " + *username + ". Ожидайте...", chatId)
		if message, err := insta.GetScanMessage(*username); err == nil {
			SendMessage(message, chatId)
			fmt.Println(message)
		} else {
			fmt.Println(err)
		}
		return true
	}
	return false
}
