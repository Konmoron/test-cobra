package utils

import (
	"regexp"
	"strings"
)

// 0.0.0.0 ~ 255.255.255.255
// https://socketloop.com/tutorials/golang-validate-ip-address
func ValidIP4(ipAddress string) bool {
	ipAddress = strings.Trim(ipAddress, " ")

	// [0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5]): 匹配0~255
	re, _ := regexp.Compile(`^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$`)
	if re.MatchString(ipAddress) {
		return true
	}
	return false
}
