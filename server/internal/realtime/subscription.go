package realtime

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"sync"
	"time"

	"graft/server/internal/httpx"
	"graft/server/internal/moduleapi"
	"graft/server/internal/realtimeauth"
)

const defaultSubscriptionTicketTTL = 30 * time.Second
const initialTopicIssuerCapacity = 4

var (
	// ErrTopicRequired reports that the realtime topic is missing.
	ErrTopicRequired   = errors.New("realtime topic required")
	// ErrTopicNotFound reports that no issuer owns the requested topic.
	ErrTopicNotFound   = errors.New("realtime topic not found")
	// ErrTopicForbidden reports that the caller cannot subscribe to the requested topic.
	ErrTopicForbidden  = errors.New("realtime topic forbidden")
	// ErrTopicConflict reports a transient failure while preparing the topic subscription.
	ErrTopicConflict   = errors.New("realtime topic unavailable")
	// ErrIssuerRequired reports that a topic issuer dependency is missing.
	ErrIssuerRequired  = errors.New("realtime subscription issuer is required")
	// ErrDuplicateIssuer reports that the same topic prefix was registered more than once.
	ErrDuplicateIssuer = errors.New("realtime subscription issuer already registered")
)

// SubscriptionRequest carries normalized request context for topic subscription issuance.
type SubscriptionRequest struct {
	Topic       string
	RequestAuth moduleapi.RequestAuthContext
	ClientIP    string
	UserAgent   string
}

// SubscriptionResponse returns the websocket bootstrap data for a realtime topic.
type SubscriptionResponse struct {
	Topic        string    `json:"topic"`
	Ticket       string    `json:"ticket"`
	WebSocketURL string    `json:"websocket_url"`
	ExpiresAt    time.Time `json:"expires_at"`
}

// SubscriptionIssuer issues websocket bootstrap data for one bounded topic family.
type SubscriptionIssuer interface {
	IssueSubscription(ctx context.Context, request SubscriptionRequest) (SubscriptionResponse, error)
}

// TopicIssuerRegistry resolves topic prefixes to their owning issuer.
type TopicIssuerRegistry interface {
	Register(prefix string, issuer SubscriptionIssuer) error
	Resolve(topic string) (SubscriptionIssuer, bool)
}

type topicIssuerRegistry struct {
	mu      sync.RWMutex
	entries []topicIssuerEntry
}

type topicIssuerEntry struct {
	prefix string
	issuer SubscriptionIssuer
}

// NewTopicIssuerRegistry creates the in-memory registry used by the unified realtime gateway.
func NewTopicIssuerRegistry() TopicIssuerRegistry {
	return &topicIssuerRegistry{
		entries: make([]topicIssuerEntry, 0, initialTopicIssuerCapacity),
	}
}

func (r *topicIssuerRegistry) Register(prefix string, issuer SubscriptionIssuer) error {
	normalizedPrefix := NormalizeTopic(prefix)
	if normalizedPrefix == "" {
		return ErrTopicRequired
	}
	if issuer == nil {
		return ErrIssuerRequired
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	for _, entry := range r.entries {
		if entry.prefix == normalizedPrefix {
			return fmt.Errorf("%w: %s", ErrDuplicateIssuer, normalizedPrefix)
		}
	}

	r.entries = append(r.entries, topicIssuerEntry{
		prefix: normalizedPrefix,
		issuer: issuer,
	})
	return nil
}

func (r *topicIssuerRegistry) Resolve(topic string) (SubscriptionIssuer, bool) {
	normalizedTopic := NormalizeTopic(topic)
	if normalizedTopic == "" {
		return nil, false
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	var matched SubscriptionIssuer
	longestPrefix := 0
	for _, entry := range r.entries {
		if !strings.HasPrefix(normalizedTopic, entry.prefix) {
			continue
		}
		if len(entry.prefix) <= longestPrefix {
			continue
		}
		longestPrefix = len(entry.prefix)
		matched = entry.issuer
	}

	return matched, matched != nil
}

// TicketIssuer delegates bounded websocket ticket issuance to the realtime auth service.
type TicketIssuer struct {
	Tickets realtimeauth.Service
}

// IssueTopicTicket issues a bounded websocket ticket for the requested topic.
func (i TicketIssuer) IssueTopicTicket(
	ctx context.Context,
	request SubscriptionRequest,
) (realtimeauth.IssuedTicket, error) {
	if i.Tickets == nil {
		return realtimeauth.IssuedTicket{}, ErrIssuerRequired
	}
	if request.RequestAuth.User == nil {
		return realtimeauth.IssuedTicket{}, ErrTopicForbidden
	}

	issued, err := i.Tickets.Issue(ctx, realtimeauth.IssueRequest{
		UserID:       request.RequestAuth.User.ID,
		ResourceType: WebSocketTopicResourceType,
		ResourceID:   request.Topic,
		Scope:        WebSocketTopicScope,
		ClientIP:     request.ClientIP,
		UserAgent:    request.UserAgent,
		TTL:          defaultSubscriptionTicketTTL,
	})
	if err != nil {
		return realtimeauth.IssuedTicket{}, err
	}
	return issued, nil
}

// BuildTopicWebSocketURL returns the canonical websocket URL for a topic ticket pair.
func BuildTopicWebSocketURL(topic string, ticket string) string {
	values := url.Values{}
	values.Set("topic", topic)
	values.Set("ticket", ticket)
	return "/ws?" + values.Encode()
}

// BuildSubscriptionRequest extracts the normalized realtime request context from the current request.
func BuildSubscriptionRequest(ctx context.Context, topic string) SubscriptionRequest {
	request := SubscriptionRequest{
		Topic: NormalizeTopic(topic),
	}
	if requestAuth, ok := moduleapi.RequestAuthContextFromContext(ctx); ok {
		request.RequestAuth = requestAuth
	}
	request.ClientIP = strings.TrimSpace(currentRequestClientIP(ctx))
	request.UserAgent = strings.TrimSpace(currentRequestUserAgent(ctx))
	return request
}

func currentRequestClientIP(ctx context.Context) string {
	requestAudit, ok := httpx.RequestAuditContextFromContext(ctx)
	if !ok {
		return ""
	}
	return strings.TrimSpace(requestAudit.ClientIP)
}

func currentRequestUserAgent(ctx context.Context) string {
	requestAudit, ok := httpx.RequestAuditContextFromContext(ctx)
	if !ok {
		return ""
	}
	return strings.TrimSpace(requestAudit.UserAgent)
}
