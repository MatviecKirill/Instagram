package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	TELEGRAM_ACCOUNTS_WHITELIST []string
	TELEGRAM_TOKEN              string
	INSTAGRAM_USERNAME          string
	INSTAGRAM_PASSWORD          string
	PROXY_URL                   string
	PROXY_LOGIN                 string
	PROXY_PASSWORD              string
	REQUEST_DELAY_MIN           int
	REQUEST_DELAY_MAX           int
}

func initConfig() (config Config) {
	var INSTAGRAM_USERNAME, INSTAGRAM_PASSWORD, PROXY_URL, PROXY_LOGIN, PROXY_PASSWORD, TELEGRAM_TOKEN string
	var REQUEST_DELAY_MIN, REQUEST_DELAY_MAX int
	var TELEGRAM_ACCOUNTS_WHITELIST []string

	TELEGRAM_ACCOUNTS_WHITELIST = strings.Split(os.Getenv("TELEGRAM_ACCOUNTS_WHITELIST"), ",")

	if TELEGRAM_TOKEN = os.Getenv("TELEGRAM_TOKEN"); TELEGRAM_TOKEN == "" {
		fmt.Println("Telegram token not found")
		return
	}
	if PROXY_URL = os.Getenv("PROXY_URL"); PROXY_URL == "" {
		fmt.Println("Proxy URL not found")
		return
	}
	if PROXY_LOGIN = os.Getenv("PROXY_LOGIN"); PROXY_LOGIN == "" {
		fmt.Println("Proxy login not found")
		return
	}
	if PROXY_PASSWORD = os.Getenv("PROXY_PASSWORD"); PROXY_PASSWORD == "" {
		fmt.Println("Proxy password not found")
		return
	}
	if INSTAGRAM_USERNAME = os.Getenv("INSTAGRAM_USERNAME"); INSTAGRAM_USERNAME == "" {
		fmt.Print("Enter instagram username: ")
		fmt.Scan(&INSTAGRAM_USERNAME)
	}
	if INSTAGRAM_PASSWORD = os.Getenv("INSTAGRAM_PASSWORD"); INSTAGRAM_PASSWORD == "" {
		fmt.Print("Enter instagram password: ")
		fmt.Scan(&INSTAGRAM_PASSWORD)
	}

	if REQUEST_DELAY_MIN, _ = strconv.Atoi(os.Getenv("REQUEST_DELAY_MIN")); REQUEST_DELAY_MIN == 0 {
		REQUEST_DELAY_MIN = 800
	}
	if REQUEST_DELAY_MAX, _ = strconv.Atoi(os.Getenv("REQUEST_DELAY_MAX")); REQUEST_DELAY_MAX == 0 {
		REQUEST_DELAY_MAX = 3500
	}

	config = Config{ TELEGRAM_ACCOUNTS_WHITELIST,
		TELEGRAM_TOKEN,
		INSTAGRAM_USERNAME,
		INSTAGRAM_PASSWORD,
		PROXY_URL,
		PROXY_LOGIN,
		PROXY_PASSWORD,
		REQUEST_DELAY_MIN,
		REQUEST_DELAY_MAX}
	fmt.Println("Config initialized")
	return config
}
