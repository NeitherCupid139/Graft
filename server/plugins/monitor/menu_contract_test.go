package monitor

import (
	"testing"

	"graft/server/internal/menu"
	monitorcontract "graft/server/plugins/monitor/contract"
)

func TestRegisterMonitorMenuIncludesThreeLevelEntries(t *testing.T) {
	t.Parallel()

	registry := menu.NewRegistry()
	registerMonitorMenu(registry, pluginID)

	menus := registry.Items()
	if len(menus) != 3 {
		t.Fatalf("expected 3 registered monitor menus, got %#v", menus)
	}

	sectionMenu := menus[0]
	if sectionMenu.Code != "monitor.section" ||
		sectionMenu.TitleKey != monitorcontract.MonitorSectionTitle.String() ||
		sectionMenu.Path != monitorcontract.MonitorGroup ||
		sectionMenu.Permission != "" {
		t.Fatalf("expected canonical monitor section menu contract, got %#v", sectionMenu)
	}

	serverStatusMenu := menus[1]
	if serverStatusMenu.Code != "monitor.server-status" ||
		serverStatusMenu.TitleKey != monitorcontract.ServerStatusMenuTitle.String() ||
		serverStatusMenu.Path != monitorcontract.ServerStatusMenuPath ||
		serverStatusMenu.Permission != "" {
		t.Fatalf("expected canonical monitor server-status menu contract, got %#v", serverStatusMenu)
	}

	overviewMenu := menus[2]
	if overviewMenu.Code != "monitor.server-status.overview" ||
		overviewMenu.TitleKey != monitorcontract.ServerStatusOverviewMenuTitle.String() ||
		overviewMenu.Path != monitorcontract.ServerStatusOverviewMenuPath ||
		overviewMenu.Permission != monitorcontract.ServerStatusReadPermission.String() {
		t.Fatalf("expected canonical monitor overview menu contract, got %#v", overviewMenu)
	}
}
