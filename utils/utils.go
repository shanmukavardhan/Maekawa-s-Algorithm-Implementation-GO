package utils

import (
	"log"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Log logs a message to the console.
func Log(message string) {
	log.Println(message)
}

// RandomDuration returns a random duration up to 1 second.
func RandomDuration() time.Duration {
	return time.Duration(rand.Intn(1000)) * time.Millisecond
}
