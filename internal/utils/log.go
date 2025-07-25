package utils

import (
	"fmt"
	"os"
	"time"
)

func Log(message string, error bool) {
	timeStamp := time.Now().Format("2006-01-02 15:04:05")
	var formattedMessage = timeStamp + " - " + message

	if error {
		fmt.Fprintln(os.Stderr, formattedMessage)
	} else {
		fmt.Println(formattedMessage)
	}
}
