package audit

import (
	"context"
	"fmt"
	"strings"

	auditcore "graft/server/internal/audit"
	"graft/server/internal/drilldown"
	auditstore "graft/server/plugins/audit/store"
)

type auditScopeResolver struct{}

func newAuditScopeResolver() auditScopeResolver {
	return auditScopeResolver{}
}

func (auditScopeResolver) SensitiveOperationKeywords() []string {
	return []string{"delete", "reset", "grant", "assign", "revoke", "remove", "replace"}
}

type scopeDefinition struct {
	category auditstore.AuditBusinessCategory
	label    string
}

func (r auditScopeResolver) Resolve(
	_ context.Context,
	metadata drilldown.ScopeMetadata,
	currentQuery auditcore.ListQuery,
) (drilldown.ResolvedScope[auditcore.ListQuery], error) {
	scope := strings.TrimSpace(metadata.Scope)
	definitions := map[string]scopeDefinition{
		string(auditstore.AuditBusinessCategoryFailedOperations):    {category: auditstore.AuditBusinessCategoryFailedOperations, label: "失败操作"},
		string(auditstore.AuditBusinessCategoryHighRiskOperations):  {category: auditstore.AuditBusinessCategoryHighRiskOperations, label: "高风险操作"},
		string(auditstore.AuditBusinessCategorySensitiveOperations): {category: auditstore.AuditBusinessCategorySensitiveOperations, label: "业务分类"},
		string(auditstore.AuditBusinessCategoryAuthFailures):        {category: auditstore.AuditBusinessCategoryAuthFailures, label: "认证失败"},
		string(auditstore.AuditBusinessCategoryPermissionDenials):   {category: auditstore.AuditBusinessCategoryPermissionDenials, label: "权限拒绝"},
		string(auditstore.AuditBusinessCategoryRBACChanges):         {category: auditstore.AuditBusinessCategoryRBACChanges, label: "权限配置变更"},
		string(auditstore.AuditBusinessCategoryCriticalSecurity):    {category: auditstore.AuditBusinessCategoryCriticalSecurity, label: "关键安全事件"},
	}

	definition, ok := definitions[scope]
	if !ok {
		return drilldown.ResolvedScope[auditcore.ListQuery]{}, drilldown.ErrScopeNotFound
	}

	if currentQuery.BusinessCategory != "" {
		return drilldown.ResolvedScope[auditcore.ListQuery]{}, fmt.Errorf("%w: business_category", drilldown.ErrScopeConflict)
	}

	return drilldown.ResolvedScope[auditcore.ListQuery]{
		Metadata: metadata,
		Applied: drilldown.AppliedScope{
			Module:      metadata.Module,
			Scope:       metadata.Scope,
			Name:        metadata.Name,
			Description: metadata.Description,
			OwnedFields: []string{"business_category"},
		},
		Projection: drilldown.ScopeProjection{
			Title:       metadata.Name,
			Description: metadata.Description,
			Items: []drilldown.ScopeProjectionItem{
				{
					Key:    "business_category",
					Label:  definition.label,
					Kind:   "enum",
					Values: []string{string(definition.category)},
					Locked: true,
				},
			},
		},
		ConvertibleFilters: drilldown.ConvertibleFilters{
			Preset:           strings.TrimSpace(string(currentQuery.TimePreset)),
			BusinessCategory: string(definition.category),
		},
		QueryPatch: auditcore.ListQuery{
			BusinessCategory: definition.category,
		},
	}, nil
}
