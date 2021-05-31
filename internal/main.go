package main

import (
	"errors"
	"fmt"
	"github.com/TheForgotten69/goinsta/v2"
	"os"
)

type Config struct {
	USERNAME   string
	PASSWORD   string
	TARGETUSER string
}

type flwType string

const (
	Followers  flwType = "Followers"
	Followings flwType = "Followings"
)

var userFollowers, userFollowings map[string][]goinsta.User

func main() {
	//https://pkg.go.dev/github.com/TheForgotten69/goinsta/v2@v2.6.0

	config := initConfig()
	if insta, err := login(config); err == nil {
		if followersList, err := getUserFlws(config.TARGETUSER, Followers, insta, 200); err == nil {
			userFollowers = make(map[string][]goinsta.User)
			userFollowers[config.TARGETUSER] = followersList
		} else {
			fmt.Println(err)
		}
		if followingsList, err := getUserFlws(config.TARGETUSER, Followings, insta, 200); err == nil {
			userFollowings = make(map[string][]goinsta.User)
			userFollowings[config.TARGETUSER] = followingsList
		} else {
			fmt.Println(err)
		}
	} else {
		fmt.Println(err)
	}
}

func getUserFlws(targetUser string, fType flwType, insta *goinsta.Instagram, limit ...int) (flws []goinsta.User, err error) {
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

func initConfig() (config Config) {
	var USERNAME string
	var PASSWORD string
	var TARGETUSER string

	if USERNAME = os.Getenv("INSTAGRAM_USERNAME"); USERNAME == "" {
		fmt.Print("Enter username: ")
		fmt.Scan(&USERNAME)
	}
	if PASSWORD = os.Getenv("INSTAGRAM_PASSWORD"); PASSWORD == "" {
		fmt.Print("Enter password: ")
		fmt.Scan(&PASSWORD)
	}
	if TARGETUSER = os.Getenv("INSTAGRAM_TARGETUSER"); TARGETUSER == "" {
		fmt.Print("Enter target username: ")
		fmt.Scan(&TARGETUSER)
	}

	config = Config{USERNAME, PASSWORD, TARGETUSER}
	fmt.Println("#Config initialized ")
	return config
}

func login(config Config) (insta *goinsta.Instagram, err error) {
	if workDir, err := os.Getwd(); err == nil {
		if insta, err := goinsta.Import(workDir + "\\config\\LoginSettings.json"); err != nil {
			insta = goinsta.New(config.USERNAME, config.PASSWORD)

			if err := insta.Login(); err == nil {
				fmt.Println("#Login successfully")
				if err := insta.Export(workDir + "\\config\\LoginSettings.json"); err != nil {
					return nil, err
				} else {
					fmt.Println("#Login data export successfully")
				}
			} else {
				return nil, err
			}

			return insta, nil
		} else {
			fmt.Println("#Login data import successfully")
			return insta, nil
		}
	} else {
		return nil, err
	}
}
