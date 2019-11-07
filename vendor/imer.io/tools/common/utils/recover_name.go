package utils

import "fmt"

func RecoverName() {
	if r := recover(); r != nil {
		fmt.Println("recovered from", r)
	}
}
