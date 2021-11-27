package insta

import (
	redisDB "InstagramStatistic/internal/Database"
	"strconv"
)

const instaURL = "https://www.instagram.com/"

func GetScanMessage(targetUserName string) (message string, err error) {
	if nonMutualUsers, unsubscribedUsers, subscribedUsers, err := getAllStatistics(targetUserName); err == nil {
		message = "Статистика для пользователя " + targetUserName + ".\n"
		message = message + "С даты: " + redisDB.Get(targetUserName+"_followers_time") + "\n"
		message = message + "Подписались: " + strconv.Itoa(len(subscribedUsers)) + ". Отписались: " + strconv.Itoa(len(unsubscribedUsers)) + ".\n"

		if subscribedUsers != nil {
			message = message + "\nСписок подписавшихся:\n"
			for i, user := range subscribedUsers {
				message = message + strconv.Itoa(i+1) + ". " + instaURL + user + "\n"
			}
		} else {
			message = message + "Нет подписавшихся.\n"
		}

		if unsubscribedUsers != nil {
			message = message + "\nСписок отписавшихся:\n"
			for i, user := range unsubscribedUsers {
				message = message + strconv.Itoa(i+1) + ". " + instaURL + user + "\n"
			}
		} else {
			message = message + "Нет отписавшихся.\n"
		}

		message = message + "\nНевзаимные подписки:\n"
		for i, user := range nonMutualUsers {
			message = message + strconv.Itoa(i+1) + ". " + user.FullName + " " + instaURL + user.Username + "\n"
		}

	} else {
		return "", err
	}
	redisDB.Set(targetUserName+"_followers_time", timeMoscow().Format("02.01.2006 15:04"))
	return message, nil
}

func GetNonMutualFollowersMessage(targetUserName string) (message string, err error) {
	if users, err := getNonMutualFollowers(targetUserName); err == nil {
		message = "Невзаимные подписки для " + targetUserName + ":\n"
		for i, user := range users {
			message = message + strconv.Itoa(i+1) + ". " + user.FullName + " " + instaURL + user.Username + "\n"
		}
		return message, nil
	} else {
		return "", err
	}
	return "", err
}

func GetUnsubscribedFollowersMessage(targetUserName string) (message string, err error) {
	if users, err := getUnsubscribedFollowers(targetUserName); err == nil {
		if users != nil {
			message = "Все отписавшиеся пользователи с даты: " + redisDB.Get(targetUserName+"_followers_time") + "\n"
			for i, user := range users {
				message = message + strconv.Itoa(i+1) + ". " + instaURL + user + "\n"
			}
		} else {
			message = "Не найдено отписавшихся пользователей для " + targetUserName + " с даты: " + redisDB.Get(targetUserName+"_followers_time")
		}
	} else {
		return "", err
	}
	redisDB.Set(targetUserName+"_followers_time", timeMoscow().Format("02.01.2006 15:04"))
	return message, nil
}
