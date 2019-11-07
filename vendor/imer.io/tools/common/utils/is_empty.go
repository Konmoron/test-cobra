package utils

func IsEmptyStringArray(arr []string) bool {
	for _, s := range arr {
		if !IsEmptyString(s) {
			return false
		}
	}
	return true
}

func IsEmptyString(s string) bool {
	if s != "" && s != " " {
		return false
	} else {
		return true
	}
}
