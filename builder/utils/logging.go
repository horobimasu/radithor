package utils

import (
	"fmt"

	"radithor-builder/store"
)

func LogDebug(msg string, args ...any) {
	if !store.DEBUG {
		return
	}

	if len(args) == 0 {
		fmt.Printf("[DEBUG] %s\n", msg)
	} else {
		fmt.Printf("[DEBUG] "+msg+"\n", args...)
	}
}

func NewLogChunk() {
	fmt.Print("\n")
}

func LogInfo(msg string, args ...any) {
	if len(args) == 0 {
		fmt.Print("[INFO] %s\n", msg)
	} else {
		fmt.Print("[INFO] %s %s\n", msg, fmt.Sprint(args...))
	}
}

func LogWarning(msg string, args ...any) {
	if len(args) == 0 {
		fmt.Print("[WARN] %s\n", msg)
	} else {
		fmt.Print("[WARN] %s %s\n", msg, fmt.Sprint(args...))
	}
}

func LogError(msg string, args ...any) {
	if len(args) == 0 {
		fmt.Print("[ERROR] %s\n", msg)
	} else {
		fmt.Print("[ERROR] %s %s\n", msg, fmt.Sprint(args...))
	}
}
