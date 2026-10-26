package utils

import (
	"fmt"
	"regexp"
	"slices"
	"strings"
)

// these funcs are exclusively used in the commands
// so there is no point making them take strings

func HasFlag(splitText []string, flag string) bool {
	return slices.Contains(splitText, "-"+flag)
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

func GetWrapped(splitText []string, flag string) string {
	text := strings.Join(splitText, " ")
	pattern := regexp.MustCompile(fmt.Sprintf("%s=\"([^\"]*)\"", flag))

	match := pattern.FindStringSubmatch(text)
	if len(match) < 2 {
		return ""
	}

	return match[1]
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
