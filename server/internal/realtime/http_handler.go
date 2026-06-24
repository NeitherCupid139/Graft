// Package realtime provides the bounded HTTP and websocket surfaces for unified topic subscriptions.
package realtime

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	messagecontract "graft/server/internal/contract/message"
	openapigen "graft/server/internal/contract/openapi/generated"
	"graft/server/internal/httpx"
	"graft/server/internal/i18n"
)

// HTTPRegistration defines dependencies for issuing realtime subscription tickets over HTTP.
type HTTPRegistration struct {
	I18n     *i18n.Service
	Registry TopicIssuerRegistry
}

// RegisterSubscriptionRoutes mounts the canonical HTTP endpoint for issuing realtime topic subscriptions.
func RegisterSubscriptionRoutes(router gin.IRouter, registration HTTPRegistration) error {
	if router == nil {
		return errors.New("realtime router is unavailable")
	}
	if registration.Registry == nil {
		return errors.New("realtime subscription registry is unavailable")
	}

	router.POST("/realtime/subscriptions", func(ctx *gin.Context) {
		var request SubscriptionRequestPayload
		if !bindSubscriptionRequest(ctx, registration.I18n, &request) {
			return
		}

		topic := NormalizeTopic(request.Topic)
		if topic == "" {
			httpx.WriteLocalizedError(ctx, registration.I18n, http.StatusBadRequest, messagecontract.CommonInvalidArgument.String(), map[string]any{
				"field": "topic",
			})
			return
		}

		issuer, ok := registration.Registry.Resolve(topic)
		if !ok {
			httpx.WriteLocalizedError(ctx, registration.I18n, http.StatusNotFound, messagecontract.CommonInvalidArgument.String(), map[string]any{
				"topic": topic,
			})
			return
		}

		response, err := issuer.IssueSubscription(ctx.Request.Context(), BuildSubscriptionRequest(ctx.Request.Context(), topic))
		if err != nil {
			writeSubscriptionError(ctx, registration.I18n, err, topic)
			return
		}

		httpx.WriteSuccess(ctx, http.StatusOK, response)
	})

	return nil
}

// SubscriptionRequestPayload is the OpenAPI-generated request shape for subscription issuance.
type SubscriptionRequestPayload = openapigen.RealtimeSubscriptionRequest

func bindSubscriptionRequest(ctx *gin.Context, localizer *i18n.Service, request *SubscriptionRequestPayload) bool {
	if ctx == nil || request == nil {
		return false
	}
	if err := ctx.ShouldBindJSON(request); err != nil {
		httpx.WriteLocalizedError(ctx, localizer, http.StatusBadRequest, messagecontract.CommonInvalidArgument.String(), map[string]any{
			"field": "body",
		})
		return false
	}
	return true
}

func writeSubscriptionError(ctx *gin.Context, localizer *i18n.Service, err error, topic string) {
	status := http.StatusInternalServerError
	messageKey := messagecontract.CommonInternalError.String()

	switch {
	case errors.Is(err, ErrTopicRequired):
		status = http.StatusBadRequest
		messageKey = messagecontract.CommonInvalidArgument.String()
	case errors.Is(err, ErrTopicNotFound):
		status = http.StatusNotFound
		messageKey = messagecontract.CommonInvalidArgument.String()
	case errors.Is(err, ErrTopicForbidden):
		status = http.StatusForbidden
		messageKey = messagecontract.AuthForbidden.String()
	case errors.Is(err, ErrTopicConflict):
		status = http.StatusConflict
		messageKey = messagecontract.CommonInternalError.String()
	}

	httpx.WriteLocalizedError(ctx, localizer, status, messageKey, map[string]any{
		"topic": topic,
	})
}
