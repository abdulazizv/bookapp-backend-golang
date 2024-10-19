package util

import (
	"fmt"
	"sync"
	"time"
)

// UniqueIDGenerator generates unique IDs.
type UniqueIDGenerator struct {
	counter uint64
	mu      sync.Mutex
}

// NewUniqueIDGenerator creates a new instance of UniqueIDGenerator.
func NewUniqueIDGenerator() *UniqueIDGenerator {
	return &UniqueIDGenerator{}
}

// GenerateID generates a unique six-digit ID based on the current timestamp and a counter.
func (u *UniqueIDGenerator) GenerateID() string {
	u.mu.Lock()
	defer u.mu.Unlock()

	// Get the current timestamp in milliseconds.
	timestamp := uint64(time.Now().UnixNano() / 1000000)

	// Combine timestamp and counter to create a unique ID.
	id := (timestamp<<16 | (u.counter & 0xFFFFF)) % 1000000

	u.counter++

	return fmt.Sprintf("%06d", id)
}

func GeneratorID() string {
	generator := NewUniqueIDGenerator()
	var id string
	for i := 0; i < 100; i++ {
		id = generator.GenerateID()
	}
	return id
}
