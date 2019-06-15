package util

import (
	"regexp"
	"strings"
)


func DomainSafeString(input string) string  {
	reg, _ := regexp.Compile("[^a-zA-Z0-9-]")

	//trim whitespace on sides
	input = strings.TrimSpace(input)
	//first, replace underscores with dashes
	input = strings.Replace(input,"_", "-", -1)
	//next, replace spaces with dashes
	input = strings.Replace(input," ", "-", -1)
	input = reg.ReplaceAllString(input, "")

	return strings.ToLower(input)
}
