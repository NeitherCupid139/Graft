package realtime

import (
	"strings"
	"sync"
	"time"
)

const defaultSubscriberBuffer = 8

type memoryHub struct {
	mu     sync.RWMutex
	topics map[string]map[uint64]chan Event
	nextID uint64
}

// NewHub creates an in-memory realtime topic hub.
func NewHub() Hub {
	return &memoryHub{
		topics: make(map[string]map[uint64]chan Event),
	}
}

func (h *memoryHub) Publish(topic string, payload any) {
	normalized := strings.TrimSpace(topic)
	if h == nil || normalized == "" {
		return
	}

	event := Event{
		Topic:      normalized,
		Data:       payload,
		OccurredAt: time.Now().UTC(),
	}

	h.mu.RLock()
	subscribers := h.topics[normalized]
	channels := make([]chan Event, 0, len(subscribers))
	for _, ch := range subscribers {
		channels = append(channels, ch)
	}
	h.mu.RUnlock()

	for _, ch := range channels {
		publishLatestEvent(ch, event)
	}
}

func publishLatestEvent(ch chan Event, event Event) {
	select {
	case ch <- event:
	default:
		drainStaleEvent(ch)
		select {
		case ch <- event:
		default:
		}
	}
}

func drainStaleEvent(ch chan Event) {
	select {
	case <-ch:
	default:
	}
}

func (h *memoryHub) Subscribe(topic string) (<-chan Event, func()) {
	normalized := strings.TrimSpace(topic)
	ch := make(chan Event, defaultSubscriberBuffer)
	if h == nil || normalized == "" {
		close(ch)
		return ch, func() {}
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	h.nextID++
	id := h.nextID
	if h.topics[normalized] == nil {
		h.topics[normalized] = make(map[uint64]chan Event)
	}
	h.topics[normalized][id] = ch

	return ch, func() {
		h.mu.Lock()
		defer h.mu.Unlock()

		subscribers := h.topics[normalized]
		if subscribers == nil {
			return
		}
		if existing, ok := subscribers[id]; ok {
			delete(subscribers, id)
			close(existing)
		}
		if len(subscribers) == 0 {
			delete(h.topics, normalized)
		}
	}
}
