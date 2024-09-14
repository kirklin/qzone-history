package utils

import (
	"regexp"
	"strconv"
	"strings"
)

// extractStringBetween 提取两个字符串之间的内容
func extractStringBetween(sourceString, startString, endString string) string {
	startIndex := strings.Index(sourceString, startString) + len(startString)
	endIndex := strings.Index(sourceString, endString)
	if startIndex < 0 || endIndex < 0 || startIndex >= endIndex {
		return ""
	}
	return sourceString[startIndex : len(sourceString)-37]
}

// replaceMultipleSpaces 去除多余的空格
func replaceMultipleSpaces(s string) string {
	re := regexp.MustCompile(`\s+`)
	return re.ReplaceAllString(s, " ")
}

// ProcessOldHTML 替换十六进制编码并处理HTML
func ProcessOldHTML(message string) string {
	re := regexp.MustCompile(`\\x[0-9a-fA-F]{2}`)
	newText := re.ReplaceAllStringFunc(message, func(hex string) string {
		byteValue, err := strconv.ParseUint(hex[2:], 16, 8)
		if err != nil {
			return hex
		}
		return string(rune(byteValue))
	})

	startString := "html:'"
	endString := "',opuin"
	newText = extractStringBetween(newText, startString, endString)
	newText = replaceMultipleSpaces(newText)
	newText = strings.ReplaceAll(newText, "\\", "")
	return newText
}
