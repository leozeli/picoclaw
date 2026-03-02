// PicoClaw - Ultra-lightweight personal AI agent
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors

package agent

import (
	"sync"

	"github.com/sipeed/picoclaw/pkg/bus"
)

// InterruptionChecker manages pending interruption messages for a session.
// This is a simplified, nanobot-inspired approach that uses message injection
// instead of task cancellation.
//
// Design Philosophy:
// - Per-session queue for isolation
// - Thread-safe for concurrent access
// - Simple API: Signal, DrainAll, HasPending
// - Zero overhead when not in use
type InterruptionChecker struct {
	queue []bus.InboundMessage
	mu    sync.Mutex
}

// NewInterruptionChecker creates a new checker for a session
func NewInterruptionChecker() *InterruptionChecker {
	return &InterruptionChecker{
		queue: make([]bus.InboundMessage, 0, 10), // Pre-allocate for common case
	}
}

// Signal pushes a new interrupting message into the queue.
// This is called when a new message arrives for an already-active session.
func (ic *InterruptionChecker) Signal(msg bus.InboundMessage) {
	ic.mu.Lock()
	defer ic.mu.Unlock()
	ic.queue = append(ic.queue, msg)
}

// DrainAll returns and clears all pending messages.
// This is called after tool execution to inject pending interruptions
// into the conversation.
func (ic *InterruptionChecker) DrainAll() []bus.InboundMessage {
	ic.mu.Lock()
	defer ic.mu.Unlock()

	if len(ic.queue) == 0 {
		return nil
	}

	// Copy messages to return
	msgs := make([]bus.InboundMessage, len(ic.queue))
	copy(msgs, ic.queue)

	// Clear queue but keep capacity to avoid reallocation
	ic.queue = ic.queue[:0]

	return msgs
}

// HasPending returns true if there are pending interruptions
func (ic *InterruptionChecker) HasPending() bool {
	ic.mu.Lock()
	defer ic.mu.Unlock()
	return len(ic.queue) > 0
}

// Peek returns the next message without removing it.
// Returns nil if queue is empty.
func (ic *InterruptionChecker) Peek() *bus.InboundMessage {
	ic.mu.Lock()
	defer ic.mu.Unlock()

	if len(ic.queue) == 0 {
		return nil
	}
	return &ic.queue[0]
}

// Len returns the number of pending messages
func (ic *InterruptionChecker) Len() int {
	ic.mu.Lock()
	defer ic.mu.Unlock()
	return len(ic.queue)
}

// Clear removes all pending messages without returning them
func (ic *InterruptionChecker) Clear() {
	ic.mu.Lock()
	defer ic.mu.Unlock()
	ic.queue = ic.queue[:0]
}
