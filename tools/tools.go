package tools

import (
	"regexp"
)

func ContrainPercentSign(context string) bool {
	reg := regexp.MustCompile(`%[^@]|%$`)
	return reg.MatchString(context)
}