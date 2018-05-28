package htmlparse

import (
	"errors"
	"strings"
	"unicode"
)

var (
	NotTagError         = errors.New("it's not a tag")
	NotCssSelectorError = errors.New("css selector syntax error")
	TagsWithoutClose    = map[string]bool{
		"br":    true,
		"img":   true,
		"hr":    true,
		"input": true,
		"link":  true,
		"meta":  true,
	}
)

func ReadWord(s []byte) []byte {
	for i := 0; i < len(s); i++ {
		if unicode.IsLetter(rune(s[i])) || unicode.IsDigit(rune(s[i])) {
			continue
		}
		return s[:i]
	}
	return []byte{}
}

//read the bytes terminate with or followed by a '<' or '>'
func ReadSegment(s []byte, offset int64) (int64, []byte, error) {
	if offset < 0 || offset >= int64(len(s)-1) {
		return 0, []byte{}, errors.New("index out of range")
	}
	var inDoubleQuote bool = false
	var inSingleQuote bool = false
	var length int64 = int64(len(s))
	for i := offset; i < length; i++ {
		if s[i] == '"' {
			inDoubleQuote = !inDoubleQuote
		}
		if s[i] == '\'' {
			inSingleQuote = !inSingleQuote
		}
		if inSingleQuote || inDoubleQuote {
			continue
		}
		if i > offset {
			if s[i] == '>' {
				return i - offset + 1, s[offset : i+1], nil
			} else if s[i] == '<' {
				return i - offset, s[offset:i], nil
			} else {
				continue
			}
		}
	}
	return length, []byte{}, nil
}

func ReadTagname(s []byte) (string, error) {
	if len(s) < 3 {
		return "", NotTagError
	}
	if s[0] != '<' {
		return "", NotTagError
	}
	for i := 0; i < len(s); i++ {
		if s[1] != '/' && (s[i] == ' ' || s[i] == '>' || i == len(s)-1) {
			return string(s[1:i]), nil
		} else if s[1] == '/' && (s[i] == ' ' || s[i] == '>' || i == len(s)-1) {
			return string(s[2:i]), nil
		} else {
			continue
		}
	}
	return "", NotTagError
}

func IsTag(text []byte) bool {
	return len(text) > 0 && text[0] == '<' && text[len(text)-1] == '>'
}

func IsSingleTag(tagName string) bool {
	if _, ok := TagsWithoutClose[tagName]; ok {
		return true
	}
	return false
}

func IsOpenTag(text []byte) bool {
	return len(text) > 0 && text[0] == '<' && text[len(text)-1] == '>' && text[1] != '/'
}

func IsCloseTag(text []byte) bool {
	return len(text) > 0 && text[0] == '<' && text[len(text)-1] == '>' && text[1] == '/'
}

func WrappedBy(str, wrap string) bool {
	if len(str)-2*len(wrap) < 0 {
		return false
	}
	if strings.HasPrefix(str, wrap) && strings.HasSuffix(str, wrap) {
		return true
	}
	return false
}
