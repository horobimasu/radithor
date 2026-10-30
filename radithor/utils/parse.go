package utils

import (
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

func HasFlag(splitText []string, flag string) bool {
	return slices.Contains(splitText, "-"+flag)
}

func HasFlagStr(text string, flag string) bool {
	return HasFlag(strings.Fields(text), flag)
}

func UnwrapQuotes(splitText []string) string {
	text := strings.Join(splitText, " ")

	first := strings.Index(text, "\"")
	last := strings.LastIndex(text, "\"")

	if first == -1 || last == -1 || first == last {
		return ""
	}

	return text[first+1 : last]
}

func UnwrapQuotesStr(text string) string {
	return UnwrapQuotes(strings.Fields(text))
}

func GetWrapped(splitText []string, flag string) (string, string) { // the match, the remainder
	text := strings.Join(splitText, " ")
	pattern := regexp.MustCompile(fmt.Sprintf("%s=\"([^\"]*)\"", flag))

	matches := pattern.FindStringSubmatchIndex(text)
	if matches == nil {
		return "", text
	}

	match := text[matches[2]:matches[3]]
	remainder := text[:matches[0]] + text[matches[1]:]

	return match, strings.Join(strings.Fields(remainder), " ")
}

func GetWrappedStr(text string, flag string) (string, string) {
	return GetWrapped(strings.Fields(text), flag)
}

// if the number is too big, the 2nd argument will return the number as a string
// if all is well, this will return the parsed number, the remainder and if a number was returned
// last arg is if a number was returned at all
func GetWrappedNumber(splitText []string, flag string) (int64, string, bool) {
	text := strings.Join(splitText, " ")
	pattern := regexp.MustCompile(fmt.Sprintf("\\b%s=(\\d+)\\b", flag))

	matches := pattern.FindStringSubmatchIndex(text)
	if matches == nil {
		return 0, text, false
	}

	match := text[matches[2]:matches[3]]

	number, err := strconv.ParseInt(match, 10, 64)
	if err != nil {
		return 0, match, true
	}

	remainder := text[:matches[0]] + text[matches[1]:]

	return number, strings.Join(strings.Fields(remainder), " "), true
}

func GetWrappedNumberStr(text string, flag string) (int64, string, bool) {
	return GetWrappedNumber(strings.Fields(text), flag)
}

// word1 word2 "word 3" \"word 4\" "word\\5"
// -> [word1, word2, word 3, "word, 4", "word\5"]
func ParseArgs(splitText []string) []string {
	var args []string
	var current []rune

	escape := false
	quoted := false
	started := false

	text := strings.Join(splitText, " ")
	for _, char := range []rune(text) {
		if escape {
			current = append(current, char)
			escape = false
			started = true
			continue
		}

		switch char {
		case '\\':
			escape = true
			started = true

		case '"':
			quoted = !quoted
			started = true

		case ' ', '\n', '\t':
			if quoted {
				current = append(current, char)
			} else {
				args = append(args, string(current))
				current = nil
				started = false
			}

		default:
			current = append(current, char)
			started = true
		}
	}

	if escape {
		current = append(current, '\\')
	}

	if started {
		args = append(args, string(current))
	}

	return args
}
