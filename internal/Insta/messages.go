package insta

import (
	redisDB "InstagramStatistic/internal/Database"
	"strconv"
)

const instaURL = "https://www.instagram.com/"

func GetNonMutualFollowersMessage(targetUserName string) (message string, err error) {
	if users, err := getNonMutualFollowers(targetUserName); err == nil {
		message = "Невзаимные подписки:\n"
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
			message = "Не найдено отписавшихся пользователей с даты: " + redisDB.Get(targetUserName+"_followers_time")
		}
	} else {
		return "", err
	}
	redisDB.Set(targetUserName+"_followers_time", timeMoscow().Format("02.01.2006 15:04"))
	return message, nil
}
