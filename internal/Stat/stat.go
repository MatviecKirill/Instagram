package stat

import (
	"errors"
	"fmt"
	"github.com/TheForgotten69/goinsta/v2"
	"os"
	"time"
)

var config Config
var insta *goinsta.Instagram
var targetUsers map[string]*goinsta.User
var usersFollowers, usersFollowings map[string][]goinsta.User

func Init() error {
	/*defer db.Close()
	initDB()*/

	config = initConfig()
	if ins, err := login(); err == nil {
		insta = ins
		usersFollowers = make(map[string][]goinsta.User)
		usersFollowings = make(map[string][]goinsta.User)
		targetUsers = make(map[string]*goinsta.User)
	} else {
		return errors.New("Login " + err.Error())
	}
	return nil
}

func getNonMutualFollowers(targetUserName string) ([]goinsta.User, error) {
	if err := getUserInfo(targetUserName); err == nil {
		if err := getUserFollowers(targetUserName); err == nil {
			if err := getUserFollowings(targetUserName); err == nil {
				return getListsDifference(usersFollowings[targetUserName], usersFollowers[targetUserName]), nil
			} else {
				return nil, err
			}
		} else {
			return nil, err
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

func getUserFollowers(targetUserName string) error {
	if targetUsers[targetUserName] != nil {
		if followersList, err := getUserFlws(targetUsers[targetUserName].Followers(), targetUsers[targetUserName].FollowerCount, 0); err == nil {
			usersFollowers[targetUserName] = followersList
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
	flwUsers = make([]goinsta.User, 0, flwCount)

	fmt.Println("Start loading users")
	for users.Next() {
		flwUsers = append(flwUsers, users.Users...)

		delay := getRandomNumber(config.REQUEST_DELAY_MIN-getRandomNumber(0, 200), config.REQUEST_DELAY_MAX+getRandomNumber(0, 500))
		time.Sleep(time.Duration(delay) * time.Millisecond)
		fmt.Printf("Delay: %v; Loaded users count: %v \n", delay, len(flwUsers))

		if len(limit) != 0 && len(flwUsers) >= limit[0] {
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
