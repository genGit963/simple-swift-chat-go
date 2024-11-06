package errorutils

import "log"

func AnyErrorCaptureLog(errorLogType string, orderOfError int8, err any) {
	if err != nil {
		log.Println("\n\t-->", errorLogType, " Error ", orderOfError, " : ", err)
	} else {
		return
	}
}
