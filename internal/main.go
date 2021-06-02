package main

import (
	stat "InstagramStatistic/internal/Stat"
	"fmt"
)

func main() {
	if err:= stat.Init(); err == nil{

	} else {
		fmt.Println(err)
	}
}
