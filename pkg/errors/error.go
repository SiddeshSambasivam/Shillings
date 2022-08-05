package error

import (
	"fmt"
	"os"
)

func HandleErrorWithExt(err error) {
	if err != nil {
		fmt.Println("Error: ", err.Error())
		os.Exit(1)
	}
}
