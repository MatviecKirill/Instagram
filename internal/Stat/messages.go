package stat

import "strconv"

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
