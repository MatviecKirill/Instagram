package insta

import (
	redisDB "InstagramStatistic/internal/Database"
	"errors"
	"fmt"
	"github.com/TheForgotten69/goinsta/v2"
	"math/rand"
	"os"
	"strconv"
	"time"
)

var username, password, proxyURL, proxyLogin, proxyPassword string
var minDelay, maxDelay int
var insta *goinsta.Instagram
var targetUsers map[string]*goinsta.User
var usersFollowers, usersFollowings map[string][]goinsta.User

func Init(username_, password_, proxyURL_, proxyLogin_, proxyPassword_ string, minDelay_, maxDelay_ int) error {
	username = username_
	password = password_
	proxyURL = proxyURL_
	proxyLogin = proxyLogin_
	proxyPassword = proxyPassword_
	minDelay = minDelay_
	maxDelay = maxDelay_

	if ins, err := login(); err == nil {
		insta = ins
		if err := insta.SetProxy("http://"+proxyLogin+":"+proxyPassword+"@"+proxyURL, true); err == nil {
			fmt.Println("Login successfully")
			usersFollowers = make(map[string][]goinsta.User)
			usersFollowings = make(map[string][]goinsta.User)
			targetUsers = make(map[string]*goinsta.User)
		} else {
			return err
		}
	} else {
		return errors.New("Login " + err.Error())
	}
	return nil
}

func login() (insta *goinsta.Instagram, err error) {
	if path, err := getWorkDir(); err == nil {
		if insta, err := goinsta.Import(path + username + ".json"); err != nil {
			insta = goinsta.New(username, password)

			if err := insta.Login(); err == nil {
				if err := insta.Export(path + username + ".json"); err != nil {
					return nil, err
				} else {
					fmt.Println("Login data export successfully")
				}
				return insta, nil
			} else {
				return nil, err
			}
		} else {
			fmt.Println("Login data import successfully")
			return insta, nil
		}
	} else {
		return nil, err
	}
}

func getUserInfo(targetUserName string) error {
	if targetUserInfo, err := insta.Profiles.ByName(targetUserName); err == nil {
		targetUsers[targetUserName] = targetUserInfo
		return nil
	} else {
		return err
	}
}

func getNonMutualFollowers(targetUserName string) ([]goinsta.User, error) {
	if err := getUserInfo(targetUserName); err == nil {
		if len(usersFollowers[targetUserName]) != targetUsers[targetUserName].FollowerCount || len(usersFollowings[targetUserName]) != targetUsers[targetUserName].FollowingCount {
			if err := getUserFollowers(targetUserName); err == nil {
				if err := getUserFollowings(targetUserName); err == nil {
					return getListsDifferenceUsers(usersFollowings[targetUserName], usersFollowers[targetUserName]), nil
				} else {
					return nil, err
				}
			} else {
				return nil, err
			}
		} else {
			return getListsDifferenceUsers(usersFollowings[targetUserName], usersFollowers[targetUserName]), nil
		}
	} else {
		return nil, err
	}
}

func getUnsubscribedFollowers(targetUserName string) ([]string, error) {
	if err := getUserInfo(targetUserName); err == nil {
		if len(usersFollowers[targetUserName]) != targetUsers[targetUserName].FollowerCount || !redisDB.Exist(targetUserName+"_followers") {
			if err := getUserFollowers(targetUserName); err != nil {
				return nil, err
			}
		}
		users := getListsDifferenceStrings(redisDB.SMembers(targetUserName+"_followers"), getUsersNamesList(usersFollowers[targetUserName]))
		if len(users) != 0 {
			redisDB.SAdd(targetUserName+"_followers", getUsersNamesList(usersFollowers[targetUserName]))
		}
		return users, nil
	} else {
		return nil, err
	}
}

func getUserFollowers(targetUserName string) error {
	if targetUsers[targetUserName] != nil {
		if followersList, err := getUserFlws(targetUsers[targetUserName].Followers(), targetUsers[targetUserName].FollowerCount, 0); err == nil {
			usersFollowers[targetUserName] = followersList
			if !redisDB.Exist(targetUserName + "_followers") {
				redisDB.Set(targetUserName+"_followers_time", timeMoscow().Format("02.01.2006 15:04"))
				redisDB.SAdd(targetUserName+"_followers", getUsersNamesList(followersList))
			}
			return nil
		} else {
			return err
		}
	}
	return errors.New(targetUserName + " not found")
}

func getUserFollowings(targetUserName string) error {
	if targetUsers[targetUserName] != nil {
		if followingsList, err := getUserFlws(targetUsers[targetUserName].Following(), targetUsers[targetUserName].FollowingCount, 0); err == nil {
			usersFollowings[targetUserName] = followingsList
			return nil
		} else {
			return err
		}
	}
	return errors.New(targetUserName + " not found")
}

func getUserFlws(users *goinsta.Users, flwCount int, limit ...int) (flwUsers []goinsta.User, err error) {
	progressStages := 0
	flwUsers = make([]goinsta.User, 0, flwCount)

	fmt.Println("Start loading users. Count: " + strconv.Itoa(flwCount))
	for users.Next() {
		flwUsers = append(flwUsers, users.Users...)
		delay := getRandomNumber(minDelay-getRandomNumber(0, 200), maxDelay+getRandomNumber(0, 500))
		time.Sleep(time.Duration(delay) * time.Millisecond)

		if progressStages != len(flwUsers)/1000 {
			progressStages = len(flwUsers) / 1000
			fmt.Printf("Loaded users count: %v/%v \n", len(flwUsers), strconv.Itoa(flwCount))
		}

		if len(limit) != 0 && limit[0] != 0 && len(flwUsers) >= limit[0] {
			return flwUsers, nil
		}
	}
	fmt.Println("Loading finished")
	return flwUsers, nil
}

func getListsDifferenceUsers(usersList1, usersList2 []goinsta.User) (diffList []goinsta.User) {
	usersMap := make(map[int64]struct{}, len(usersList2))
	for _, user := range usersList2 {
		usersMap[user.ID] = struct{}{}
	}
	for _, user := range usersList1 {
		if _, found := usersMap[user.ID]; !found {
			diffList = append(diffList, user)
		}
	}
	return diffList
}

func getListsDifferenceStrings(usersList1, usersList2 []string) (diffList []string) {
	usersMap := make(map[string]struct{}, len(usersList2))
	for _, user := range usersList2 {
		usersMap[user] = struct{}{}
	}
	for _, user := range usersList1 {
		if _, found := usersMap[user]; !found {
			diffList = append(diffList, user)
		}
	}
	return diffList
}

func getUsersNamesList(users []goinsta.User) []string {
	usersNames := make([]string, 0, len(users))
	for _, user := range users {
		usersNames = append(usersNames, user.Username)
	}
	return usersNames
}

func getWorkDir() (path string, err error) {
	if workDir, err := os.Getwd(); err == nil {
		path = workDir + "\\accounts"
		if _, err := os.Stat(path); os.IsNotExist(err) {
			if err := os.Mkdir(path, os.ModeDir); err == nil {
				path = path + "\\"
				return path, nil
			} else {
				path = "./"
				return path, nil
			}
		} else {
			path = path + "\\"
			return path, nil
		}
	} else {
		return "", err
	}
}

func getRandomNumber(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min) + min
}

func timeMoscow() time.Time {
	loc, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		panic(err)
	}
	return time.Now().In(loc)
}
