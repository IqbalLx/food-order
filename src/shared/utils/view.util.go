package utils

import (
	"strconv"
	"strings"
)

func FormatPrice(price int) string {
	currencyStr := strconv.Itoa(price)
	length := len(currencyStr)

	var formatted string

	if length <= 3 {
		formatted = currencyStr
	} else {
		formatted = currencyStr[:length-3] + "." + currencyStr[length-3:]
	}

	return formatted
}

func EncodeQuerystring(input string) string {
	// Split the input string into words
	words := strings.Fields(input)

	// Join the words with '%20'
	result := strings.Join(words, "%20")

	return result
}