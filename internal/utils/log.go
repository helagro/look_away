package utils

import (
	"fmt"
	"time"
)

func Log(message string) {
	timeStamp := time.Now().Format("2006-01-02 15:04:05")

	fmt.Println(timeStamp, message)
}
