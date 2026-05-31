package audit

import (
	"context"
	"fmt"
	"strings"

	auditcore "graft/server/internal/audit"
	"graft/server/internal/drilldown"
)

type auditScopeResolver struct{}

func newAuditScopeResolver() auditScopeResolver {
	return auditScopeResolver{}
}

func (auditScopeResolver) SensitiveOperationKeywords() []string {
	return []string{"delete", "reset", "grant", "assign", "revoke", "remove", "replace"}
}

func (r auditScopeResolver) Resolve(
	_ context.Context,
	metadata drilldown.ScopeMetadata,
	currentQuery auditcore.ListQuery,
) (drilldown.ResolvedScope[auditcore.ListQuery], error) {
	switch strings.TrimSpace(metadata.Scope) {
	case "sensitive_operations":
		if len(currentQuery.ActionKeywords) > 0 {
			return drilldown.ResolvedScope[auditcore.ListQuery]{}, fmt.Errorf("%w: action_keywords", drilldown.ErrScopeConflict)
		}

		keywords := r.SensitiveOperationKeywords()
		return drilldown.ResolvedScope[auditcore.ListQuery]{
			Metadata: metadata,
			Applied: drilldown.AppliedScope{
				Module:      metadata.Module,
				Scope:       metadata.Scope,
				Name:        metadata.Name,
				Description: metadata.Description,
				OwnedFields: []string{"action_keywords"},
			},
			Projection: drilldown.ScopeProjection{
				Title:       metadata.Name,
				Description: metadata.Description,
				Items: []drilldown.ScopeProjectionItem{
					{
						Key:    "action_keywords",
						Label:  "操作关键词",
						Kind:   "string_list",
						Values: append([]string(nil), keywords...),
						Locked: true,
					},
				},
			},
			ConvertibleFilters: drilldown.ConvertibleFilters{
				ActionKeywords: append([]string(nil), keywords...),
				Preset:         strings.TrimSpace(string(currentQuery.TimePreset)),
			},
			QueryPatch: auditcore.ListQuery{
				ActionKeywords: keywords,
			},
		}, nil
	default:
		return drilldown.ResolvedScope[auditcore.ListQuery]{}, drilldown.ErrScopeNotFound
	}
}
