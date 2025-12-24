package notifications

import (
	"os/exec"

	"github.com/helagro/look_away/internal/utils"
)

func getNotificationMessage(textCommand string) (string, error) {
	if textCommand != "" {
		cmd := exec.Command(textCommand)

		output, err := cmd.Output()
		if err == nil {
			utils.Log("Command output is: "+string(output), false)
			return string(output), nil
		} else {
			var errorMessage = "Error executing command: " + err.Error()
			utils.Log(errorMessage, true)

			return "Error executing command: " + err.Error(), err
		}
	} else {
		utils.Log("No text command provided, using default message", false)
		return "Time to rest your eyes! Look at least 6m away for at least 20 seconds!", nil
	}
}
