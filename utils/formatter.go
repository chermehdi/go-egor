package utils

import (
	"fmt"
	"strings"
)

func GetHeader(headerText, headerBorder string, headerRepeat int) string {
	return fmt.Sprintf("%s\n%s", headerText, strings.Repeat(headerBorder, headerRepeat))
}
