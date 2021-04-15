package validation

import "strings"

func IsExist(input string, stringArray []string) bool {
	input = strings.ToLower(input)
	for _, str := range stringArray {
		if strings.ToLower(str) == input {
			return true
		}
	}
	return false
}
