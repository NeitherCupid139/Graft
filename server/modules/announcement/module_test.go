// Copyright (c) 2025-2026 GeWuYou
// SPDX-License-Identifier: Apache-2.0

package announcement

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"graft/server/internal/config"
	"graft/server/internal/configregistry"
	"graft/server/internal/container"
	"graft/server/internal/cronx"
	"graft/server/internal/dashboard"
	"graft/server/internal/i18n"
	"graft/server/internal/menu"
	"graft/server/internal/module"
	"graft/server/internal/moduleapi"
	"graft/server/internal/permission"
	announcementcontract "graft/server/modules/announcement/contract"
	announcementstore "graft/server/modules/announcement/store"
)

func TestModuleRegistersAnnouncementMetadata(t *testing.T) {
	service, err := NewService(testAnnouncementRepository{})
	if err != nil {
		t.Fatalf("new service: %v", err)
	}
	ctx := newAnnouncementTestContext(nil)
	moduleInstance := NewModule(service)
	if err := moduleInstance.Register(ctx); err != nil {
		t.Fatalf("register announcement module: %v", err)
	}

	assertAnnouncementPermissionsRegistered(t, ctx.PermissionRegistry)
	assertAnnouncementMenuRegistered(t, ctx.MenuRegistry)
	assertAnnouncementMessageRegistered(t, ctx.I18n, i18n.LocaleZHCN, "公告管理")
	assertAnnouncementMessageRegistered(t, ctx.I18n, i18n.LocaleENUS, "Announcements")
}

func TestNewModuleSpecDeclaresMigrationAndDependencies(t *testing.T) {
	spec := NewModuleSpec()
	if spec.ID != moduleID {
		t.Fatalf("unexpected module id %q", spec.ID)
	}
	if len(spec.Dependencies) != 2 || spec.Dependencies[0] != "user" || spec.Dependencies[1] != "rbac" {
		t.Fatalf("unexpected dependencies %#v", spec.Dependencies)
	}
	if len(spec.MigrationPath) != 1 || spec.MigrationPath[0] != "modules/announcement/migrations" {
		t.Fatalf("unexpected migration paths %#v", spec.MigrationPath)
	}
}

func TestAnnouncementContractValidators(t *testing.T) {
	if !announcementcontract.ValidAnnouncementStatus(announcementcontract.AnnouncementStatusPublished) {
		t.Fatal("expected published status to be valid")
	}
	if announcementcontract.ValidAnnouncementStatus(announcementcontract.AnnouncementStatus("visible")) {
		t.Fatal("unexpected ad-hoc announcement status accepted")
	}
	if !announcementcontract.ValidAnnouncementLevel(announcementcontract.AnnouncementLevelWarning) {
		t.Fatal("expected warning level to be valid")
	}
	if announcementcontract.ValidAnnouncementLevel(announcementcontract.AnnouncementLevel("critical")) {
		t.Fatal("unexpected ad-hoc announcement level accepted")
	}
}

func TestAnnouncementRoutesReturnExplicitPhaseOneNotImplemented(t *testing.T) {
	gin.SetMode(gin.TestMode)
	service, err := NewService(testAnnouncementRepository{})
	if err != nil {
		t.Fatalf("new service: %v", err)
	}
	engine := gin.New()
	ctx := newAnnouncementTestContext(engine)
	if err := registerAnnouncementRoutes(ctx, service, announcementGuards{
		authenticated: announcementRouteTestAuth(42),
		read:          announcementRouteTestAuth(42),
		create:        announcementRouteTestAuth(42),
		update:        announcementRouteTestAuth(42),
		publish:       announcementRouteTestAuth(42),
		delete:        announcementRouteTestAuth(42),
	}); err != nil {
		t.Fatalf("register announcement routes: %v", err)
	}

	request := httptest.NewRequest(http.MethodGet, "/api/announcements", nil)
	recorder := httptest.NewRecorder()
	engine.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusNotImplemented {
		t.Fatalf("expected 501 phase-one route response, got %d body=%s", recorder.Code, recorder.Body.String())
	}
}

func assertAnnouncementPermissionsRegistered(t *testing.T, registry *permission.Registry) {
	t.Helper()
	registered := make(map[string]struct{}, len(registry.Items()))
	for _, item := range registry.Items() {
		registered[item.Code] = struct{}{}
	}
	for _, code := range []string{
		announcementcontract.AnnouncementReadPermission.String(),
		announcementcontract.AnnouncementCreatePermission.String(),
		announcementcontract.AnnouncementUpdatePermission.String(),
		announcementcontract.AnnouncementPublishPermission.String(),
		announcementcontract.AnnouncementDeletePermission.String(),
	} {
		if _, ok := registered[code]; !ok {
			t.Fatalf("expected announcement permission %s to be registered", code)
		}
	}
}

func assertAnnouncementMenuRegistered(t *testing.T, registry *menu.Registry) {
	t.Helper()
	for _, item := range registry.Items() {
		if item.Code == "announcement.list" &&
			item.Path == announcementcontract.AnnouncementMenuPath &&
			item.Permission == announcementcontract.AnnouncementReadPermission.String() &&
			item.TitleKey == announcementcontract.AnnouncementMenuTitle.String() {
			return
		}
	}
	t.Fatalf("expected announcement management menu, got %#v", registry.Items())
}

func assertAnnouncementMessageRegistered(t *testing.T, localizer *i18n.Service, locale i18n.LocaleTag, expected string) {
	t.Helper()
	matches := localizer.RegisteredMessageResources(locale, i18n.MessageKey(announcementcontract.AnnouncementMenuTitle.String()))
	if len(matches) != 1 || matches[0].Text != expected {
		t.Fatalf("expected announcement menu title %q for %s, got %#v", expected, locale, matches)
	}
}

func newAnnouncementTestContext(engine *gin.Engine) *module.Context {
	localizer := i18n.MustNew(config.I18nConfig{
		DefaultLocale:  "zh-CN",
		FallbackLocale: "zh-CN",
		SupportedLocales: []string{
			"zh-CN",
			"en-US",
		},
	})
	var router gin.IRouter
	if engine != nil {
		router = engine.Group("/api")
	}
	return &module.Context{
		Config:             &config.Config{},
		Router:             router,
		I18n:               localizer,
		Services:           container.New(),
		MenuRegistry:       menu.NewRegistry(),
		PermissionRegistry: permission.NewRegistry(),
		CronRegistry:       cronx.NewRegistry(),
		ConfigRegistry:     configregistry.NewRegistry(),
		DashboardRegistry:  dashboard.NewRegistry(),
	}
}

func announcementRouteTestAuth(userID uint64) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Request = ctx.Request.WithContext(moduleapi.WithRequestAuthContext(ctx.Request.Context(), moduleapi.RequestAuthContext{
			User: &moduleapi.CurrentUser{ID: userID, Username: "alice"},
		}))
		ctx.Next()
	}
}

type testAnnouncementRepository struct{}

func (testAnnouncementRepository) Ping(context.Context) error {
	return nil
}

var _ announcementstore.Repository = testAnnouncementRepository{}
