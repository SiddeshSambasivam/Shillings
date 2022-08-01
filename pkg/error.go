package pkg

import (
	"fmt"
	"os"
)

func HandleErrorWithExt(err error) {
	if err != nil {
		fmt.Println("error (my message): ", err.Error())
		os.Exit(1)
	}
}
