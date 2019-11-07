package utils

import (
	"log"
	"time"
)

func ExecTime(funcName string, start time.Time) {
	tc := time.Since(start)
	log.Printf("func %s speed time: %v\n", funcName, tc)
}
