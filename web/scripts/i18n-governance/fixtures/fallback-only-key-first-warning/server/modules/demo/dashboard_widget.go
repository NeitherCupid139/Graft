// Copyright (c) 2025-2026 GeWuYou
// SPDX-License-Identifier: Apache-2.0

package demo

type WidgetDefinition struct {
	Title       string
	TitleKey    string
	Description string
}

func demoWidget() WidgetDefinition {
	return WidgetDefinition{
		Title:       "Dashboard title",
		Description: "Dashboard description",
	}
}
