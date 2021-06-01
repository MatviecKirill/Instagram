package main

import (
	"errors"
	"fmt"
	"github.com/TheForgotten69/goinsta/v2"
	"os"
	"time"
)

type flwType string

const (
	Followers  flwType = "Followers"
	Followings flwType = "Followings"
)

var userFollowers, userFollowings map[string][]goinsta.User

func main() {
	//https://pkg.go.dev/github.com/TheForgotten69/goinsta/v2@v2.6.0

	config := initConfig()
	if insta, err := login(&config); err == nil {
		if followersList, err := getUserFlws(config.TARGETUSER, Followers, insta, config.REQUEST_DELAY_MIN, config.REQUEST_DELAY_MAX, 200); err == nil {
			userFollowers = make(map[string][]goinsta.User)
			userFollowers[config.TARGETUSER] = followersList
		} else {
			fmt.Println(err)
		}
		if followingsList, err := getUserFlws(config.TARGETUSER, Followings, insta, config.REQUEST_DELAY_MIN, config.REQUEST_DELAY_MAX, 200); err == nil {
			userFollowings = make(map[string][]goinsta.User)
			userFollowings[config.TARGETUSER] = followingsList
		} else {
			fmt.Println(err)
		}
	} else {
		fmt.Println("Login", err)
	}
}

func getListsDifference(usersList1, usersList2 []goinsta.User) (diffList []goinsta.User) {
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

func getUserFlws(targetUser string, fType flwType, insta *goinsta.Instagram, delayMin int, delayMax int, limit ...int) (flws []goinsta.User, err error) {
	if searchResult, err := insta.Search.User(targetUser); err == nil {
		flwUsers := make([]goinsta.User, 0)
		var flws *goinsta.Users

		switch fType {
		case Followers:
			flws = searchResult.Users[0].Followers()
		case Followings:
			flws = searchResult.Users[0].Following()
		}

		if flws != nil {
			for flws.Next() {
				flwUsers = append(flwUsers, flws.Users...)

				delay := getRandomNumber(delayMin-getRandomNumber(0, 200), delayMax+getRandomNumber(0, 500))
				time.Sleep(time.Duration(delay) * time.Millisecond)
				fmt.Printf("Delay: %v, %v users: %v \n", delay, fType, len(flwUsers))

				if len(limit) != 0 && len(flwUsers) >= limit[0] {
					return flwUsers, nil
				}
			}
			return flwUsers, nil
		} else {
			return nil, errors.New(string(fType) + " not found")
		}
	} else {
		return nil, err
	}
}

func login(config *Config) (insta *goinsta.Instagram, err error) {
	if workDir, err := os.Getwd(); err == nil {
		if insta, err := goinsta.Import(workDir + "\\accounts\\" + config.USERNAME + ".json"); err != nil {
			insta = goinsta.New(config.USERNAME, config.PASSWORD)

			if err := insta.Login(); err == nil {
				fmt.Println("Login successfully")
				if err := insta.Export(workDir + "\\accounts\\" + config.USERNAME + ".json"); err != nil {
					return nil, err
				} else {
					fmt.Println("Login data export successfully")
				}
			} else {
				return nil, err
			}

			return insta, nil
		} else {
			fmt.Println("Login data import successfully")
			return insta, nil
		}
	} else {
		return nil, err
	}
}
