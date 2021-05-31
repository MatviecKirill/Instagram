package main

import (
	"fmt"
	"github.com/TheForgotten69/goinsta/v2"
	"os"
)

type Config struct {
	USERNAME   string
	PASSWORD   string
	TARGETUSER string
}

func main() {
	//https://pkg.go.dev/github.com/TheForgotten69/goinsta/v2@v2.6.0

	config := initConfig()
	if insta, err := login(config); err == nil {
		if followers, err := getUserFollowers(config.TARGETUSER, insta, 200); err == nil {
			for _, follower := range followers {
				fmt.Println(follower)
			}
		} else {
			fmt.Println(err)
		}
	} else {
		fmt.Println(err)
	}
}

func login(config Config) (insta *goinsta.Instagram, err error) {
	if workDir, err := os.Getwd(); err == nil {
		if insta, err := goinsta.Import(workDir + "/config/LoginSettings.json"); err != nil {
			insta = goinsta.New(config.USERNAME, config.PASSWORD)

			if err := insta.Login(); err != nil {
				return nil, err
			}

			insta.Export(workDir + "/config/LoginSettings.json")
			return insta, nil
		} else {
			return insta, nil
		}
	} else {
		return nil, err
	}
}

func getUserFollowers(userName string, insta *goinsta.Instagram, limit ...int) (followers []goinsta.User, err error) {
	if searchResult, err := insta.Search.User(userName); err == nil {
		followerUsers := make([]goinsta.User, 0)
		followers := searchResult.Users[0].Followers()

		for followers.Next() {
			followerUsers = append(followerUsers, followers.Users...)

			if len(limit) != 0 && len(followerUsers) >= limit[0] {
				return followerUsers, nil
			}
		}
		return followerUsers, nil
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
	fmt.Print("CONFIG INITIALIZED ")
	return config
}
