package stringx

import (
	"github.com/satori/go.uuid"
	"regexp"
)

var replaceCtrlChar = regexp.MustCompile(`[[:cntrl:]]`)

func IsStringUuid(input string) bool {
	_, err := uuid.FromString(input)
	if err != nil {
		return false
	}
	return true
}
func IsStrExistInArray(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func ReplaceCtrlChar(s string) string {
	r := replaceCtrlChar.ReplaceAllLiteralString(s, "")
	return r
}
