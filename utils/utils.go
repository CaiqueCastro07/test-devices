package utils

import (
	"math"
	math2 "math/rand/v2"
	"regexp"
	"strconv"
	"strings"
)

var onlyNumbersRegex = regexp.MustCompile(`[^0-9]`)

func OnlyNumbers(str string) string {
	if len(str) < 1 {
		return ""
	}
	return onlyNumbersRegex.ReplaceAllString(str, "")
}

var lowerCaseNoSpaceRegex = regexp.MustCompile(`/\s/g`)

func LowerCaseNoSpace(str string) string {
	if len(str) < 1 {
		return ""
	}
	return strings.ToLower(lowerCaseNoSpaceRegex.ReplaceAllString(str, ""))
}

func FloatToStr(num float64) string {
	return strconv.FormatFloat(num, 'f', 0, 64)
}

func RandomNumber(min int, max int) int {

	if max < 1 || min < 0 {
		return 0
	}

	return math2.IntN(max-min) + min
}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func ToFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}

func RemoveNonAlphanumeric(s string) string {
	var result strings.Builder
	for i := 0; i < len(s); i++ {
		b := s[i]
		if ('a' <= b && b <= 'z') ||
			('A' <= b && b <= 'Z') ||
			('0' <= b && b <= '9') ||
			b == ' ' {
			result.WriteByte(b)
		}
	}

	return strings.Trim(result.String(), " ")
}
