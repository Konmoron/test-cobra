package utils

import "fmt"

func CheckErr(errMess error) {
	if errMess != nil {
		fmt.Println(errMess)
	}
}
