package util

import (
	"fmt"
	"runtime"
)

func Fname() string {
	fpcs := make([]uintptr, 1)
	i := runtime.Callers(2, fpcs)
	if i == 0 {
		return fmt.Sprintln("No caller")
	}
	caller := runtime.FuncForPC(fpcs[0] - 1)
	if caller == nil {
		return fmt.Sprintln("Caller was nil")
	}
	return fmt.Sprintln(caller.Name())
}

func FnameN(n int) string {
	fpcs := make([]uintptr, 1)
	i := runtime.Callers(n, fpcs)
	if i == 0 {
		return fmt.Sprintln("No caller")
	}
	caller := runtime.FuncForPC(fpcs[0] - 1)
	if caller == nil {
		return fmt.Sprintln("Caller was nil")
	}
	return fmt.Sprintln(caller.Name())
}

