package main

import (
	"fmt"
	"github.com/TheForgotten69/goinsta/v2"
	"os"
	"time"
)

var config Config
var insta *goinsta.Instagram
var targetUsers map[string]*goinsta.User
var usersFollowers, usersFollowings map[string][]goinsta.User

func main() {
	defer db.Close()
	initDB()

	config = initConfig()
	//return
	if ins, err := login(); err == nil {
		insta = ins
		usersFollowers = make(map[string][]goinsta.User)
		usersFollowings = make(map[string][]goinsta.User)

		getUserInfo(config.TARGETUSER)
	} else {
		fmt.Println("Login", err)
	}
}

func getUserInfo(targetUserName string) {
	if targetUserInfo, err := insta.Profiles.ByName(targetUserName); err == nil {
		targetUsers[targetUserName] = targetUserInfo
	} else {
		fmt.Println(err)
	}
}

func getUserFollowers(targetUserName string){
	if followersList, err := getUserFlws(targetUsers[targetUserName].Followers(), 200); err == nil {
		usersFollowers[targetUserName] = followersList
	} else {
		fmt.Println(err)
	}
}

func getUserFollowings(targetUserName string)  {
	if followingsList, err := getUserFlws(targetUsers[targetUserName].Following(), 200); err == nil {
		usersFollowings[targetUserName] = followingsList
	} else {
		fmt.Println(err)
	}
}

func getUserFlws(users *goinsta.Users, limitFlwsCount ...int) (flwUsers []goinsta.User, err error) {
	flwUsers = make([]goinsta.User, 0)

	for users.Next() {
		flwUsers = append(flwUsers, users.Users...)

		delay := getRandomNumber(config.REQUEST_DELAY_MIN-getRandomNumber(0, 200), config.REQUEST_DELAY_MAX+getRandomNumber(0, 500))
		time.Sleep(time.Duration(delay) * time.Millisecond)
		fmt.Printf("Delay: %v; users: %v \n", delay, len(flwUsers))

		if len(limitFlwsCount) != 0 && len(flwUsers) >= limitFlwsCount[0] {
			return flwUsers, nil
		}
	}
	return flwUsers, nil
}

func login() (insta *goinsta.Instagram, err error) {
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
