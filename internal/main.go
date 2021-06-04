package main

import (
	stat "InstagramStatistic/internal/Stat"
	"fmt"
)

func main() {
	if err := stat.Init(); err == nil {
		if message, err := stat.GetNonMutualFollowersMessage(""); err == nil {
			fmt.Print(message)
		} else {
			fmt.Println(err)
		}
	} else {
		fmt.Println(err)
	}
}
