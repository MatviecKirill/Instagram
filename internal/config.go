package main

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	TOKEN             string
	USERNAME          string
	PASSWORD          string
	REQUEST_DELAY_MIN int
	REQUEST_DELAY_MAX int
}

func initConfig() (config Config) {
	var USERNAME, PASSWORD, TOKEN string
	var REQUEST_DELAY_MIN, REQUEST_DELAY_MAX int

	if TOKEN = os.Getenv("TELEGRAM_TOKEN"); TOKEN == "" {
		fmt.Println("Telegram token not found")
		return
	}
	if USERNAME = os.Getenv("INSTAGRAM_USERNAME"); USERNAME == "" {
		fmt.Print("Enter username: ")
		fmt.Scan(&USERNAME)
	}
	if PASSWORD = os.Getenv("INSTAGRAM_PASSWORD"); PASSWORD == "" {
		fmt.Print("Enter password: ")
		fmt.Scan(&PASSWORD)
	}

	if REQUEST_DELAY_MIN, _ = strconv.Atoi(os.Getenv("REQUEST_DELAY_MIN")); REQUEST_DELAY_MIN == 0 {
		REQUEST_DELAY_MIN = 800
	}
	if REQUEST_DELAY_MAX, _ = strconv.Atoi(os.Getenv("REQUEST_DELAY_MAX")); REQUEST_DELAY_MAX == 0 {
		REQUEST_DELAY_MAX = 3500
	}

	config = Config{TOKEN, USERNAME, PASSWORD, REQUEST_DELAY_MIN, REQUEST_DELAY_MAX}
	fmt.Println("Config initialized")
	return config
}
