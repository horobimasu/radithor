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
		fmt.Printf("[INFO] %s\n", msg)
	} else {
		fmt.Printf("[INFO] "+msg+"\n", args...)
	}
}

func LogWarning(msg string, args ...any) {
	if len(args) == 0 {
		fmt.Printf("[WARN] %s\n", msg)
	} else {
		fmt.Printf("[WARN] "+msg+"\n", args...)
	}
}

func LogError(msg string, args ...any) {
	if len(args) == 0 {
		fmt.Printf("[ERROR] %s\n", msg)
	} else {
		fmt.Printf("[ERROR] "+msg+"\n", args...)
	}
}
