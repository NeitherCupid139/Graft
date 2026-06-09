package dashboard

import (
	"testing"

	"graft/server/internal/configregistry"
)

func TestRegisterQuickActionsConfigDefinitionsUsesDomainGroupItemMetadata(t *testing.T) {
	t.Parallel()

	registry := configregistry.NewRegistry()
	if err := RegisterQuickActionsConfigDefinitions(registry); err != nil {
		t.Fatalf("register quick-actions config definitions: %v", err)
	}

	items := registry.Items()
	if len(items) != 3 {
		t.Fatalf("expected three quick-action config items, got %#v", items)
	}
	for _, item := range items {
		if item.Domain != quickActionsConfigDomain ||
			item.DomainKey != quickActionsConfigDomainKey ||
			item.Group != quickActionsConfigGroup ||
			item.GroupKey != quickActionsConfigGroupKey ||
			item.GroupDescriptionKey != quickActionsConfigGroupDescKey {
			t.Fatalf("expected dashboard quick-action hierarchy metadata, got %#v", item)
		}
		if item.GroupLabel == "core / dashboard.quick_actions" {
			t.Fatalf("group label must be product-facing fallback, got %q", item.GroupLabel)
		}
	}
}
