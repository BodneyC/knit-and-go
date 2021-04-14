package util

import (
	"encoding/json"
	"fmt"
	"runtime"
	"unicode"
)

func StripLeadingWhitespace(s string) string {
	i := len(s)
	for i > 0 && unicode.IsSpace(rune(s[i-1])) {
		i--
	}
	return s[0:i]
}

func JsonEscape(i string) string {
	b, err := json.Marshal(i)
	if err != nil {
		panic(err)
	}
	s := string(b)
	return s[1 : len(s)-1]
}

func stackLine(pc uintptr, file string, line int, ok bool) string {
	if !ok {
		panic("Error retriving stack information")
	}
	return fmt.Sprintf("\n %s:%d\n  @%s", runtime.FuncForPC(pc).Name(), line, file)
}

func StackLine() string {
	return stackLine(runtime.Caller(1))
}

func StackLineN(n int) string {
	return stackLine(runtime.Caller(n))
}
