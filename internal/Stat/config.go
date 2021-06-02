package stat

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
)

type Config struct {
	USERNAME          string
	PASSWORD          string
	TARGETUSER        string
	REQUEST_DELAY_MIN int
	REQUEST_DELAY_MAX int
}

func initConfig() (config Config) {
	var USERNAME, PASSWORD, TARGETUSER string
	var REQUEST_DELAY_MIN, REQUEST_DELAY_MAX int

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
	if REQUEST_DELAY_MIN, _ = strconv.Atoi(os.Getenv("REQUEST_DELAY_MIN")); REQUEST_DELAY_MIN == 0 {
		REQUEST_DELAY_MIN = getRandomNumber(800, 1100)
	}
	if REQUEST_DELAY_MAX, _ = strconv.Atoi(os.Getenv("REQUEST_DELAY_MAX")); REQUEST_DELAY_MAX == 0 {
		REQUEST_DELAY_MAX = getRandomNumber(2500, 3500)
	}

	config = Config{USERNAME, PASSWORD, TARGETUSER, REQUEST_DELAY_MIN, REQUEST_DELAY_MAX}
	fmt.Println("Config initialized")
	return config
}

func getRandomNumber(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min) + min
}
