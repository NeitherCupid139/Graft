// Copyright (c) 2025-2026 GeWuYou
// SPDX-License-Identifier: Apache-2.0

package announcement

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	"graft/server/internal/config"
	"graft/server/internal/configregistry"
	"graft/server/internal/container"
	messagecontract "graft/server/internal/contract/message"
	"graft/server/internal/cronx"
	"graft/server/internal/dashboard"
	"graft/server/internal/httpx"
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

func TestAnnouncementUserRoutesRemainExplicitPhaseThreeNotImplemented(t *testing.T) {
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

	request := httptest.NewRequest(http.MethodGet, "/api/my/announcements", nil)
	recorder := httptest.NewRecorder()
	engine.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusNotImplemented {
		t.Fatalf("expected 501 phase-three route response, got %d body=%s", recorder.Code, recorder.Body.String())
	}
}

func TestAnnouncementManagementServiceLifecycle(t *testing.T) {
	repository := newMemoryAnnouncementRepository()
	service, err := NewService(repository)
	if err != nil {
		t.Fatalf("new service: %v", err)
	}
	ctx := context.Background()
	actorID := uint64(7)
	publishAt := time.Date(2026, 6, 12, 8, 0, 0, 0, time.FixedZone("cst", 8*60*60))
	expireAt := publishAt.Add(2 * time.Hour)

	created, err := service.Create(ctx, announcementstore.CreateInput{
		Title:     " Maintenance ",
		Content:   "Window",
		Level:     announcementcontract.AnnouncementLevelWarning.String(),
		Pinned:    true,
		PublishAt: &publishAt,
		ExpireAt:  &expireAt,
		ActorID:   &actorID,
	})
	if err != nil {
		t.Fatalf("create announcement: %v", err)
	}
	if created.Status != announcementcontract.AnnouncementStatusDraft.String() {
		t.Fatalf("expected draft create status, got %q", created.Status)
	}

	published, err := service.Publish(ctx, created.ID, nil, &actorID)
	if err != nil {
		t.Fatalf("publish announcement: %v", err)
	}
	if published.Status != announcementcontract.AnnouncementStatusPublished.String() {
		t.Fatalf("expected published status, got %q", published.Status)
	}
	if published.PublishAt == nil || !published.PublishAt.Equal(publishAt.UTC()) {
		t.Fatalf("expected publish_at to keep UTC input, got %#v", published.PublishAt)
	}
	if err := service.Delete(ctx, created.ID, actorID); !errors.Is(err, errAnnouncementPublishedDelete) {
		t.Fatalf("expected published delete guard, got %v", err)
	}

	archived, err := service.Archive(ctx, created.ID, &actorID)
	if err != nil {
		t.Fatalf("archive announcement: %v", err)
	}
	if archived.Status != announcementcontract.AnnouncementStatusArchived.String() {
		t.Fatalf("expected archived status, got %q", archived.Status)
	}
	if _, err := service.Update(ctx, created.ID, announcementstore.UpdateInput{
		Title:   "New",
		Content: "New",
		Level:   announcementcontract.AnnouncementLevelInfo.String(),
	}); !errors.Is(err, errAnnouncementInvalidTransition) {
		t.Fatalf("expected archived update guard, got %v", err)
	}
	if err := service.Delete(ctx, created.ID, actorID); err != nil {
		t.Fatalf("delete archived announcement: %v", err)
	}
	if _, err := service.GetAdmin(ctx, created.ID); !errors.Is(err, errAnnouncementNotFound) {
		t.Fatalf("expected deleted announcement not found, got %v", err)
	}
}

func TestAnnouncementManagementServiceDeleteDraftAndInvalidPublishWindow(t *testing.T) {
	repository := newMemoryAnnouncementRepository()
	service, err := NewService(repository)
	if err != nil {
		t.Fatalf("new service: %v", err)
	}
	ctx := context.Background()
	actorID := uint64(7)
	publishAt := time.Date(2026, 6, 12, 8, 0, 0, 0, time.UTC)
	expireAt := publishAt

	if _, err := service.Create(ctx, announcementstore.CreateInput{
		Title:     "Bad window",
		Content:   "Bad window",
		Level:     announcementcontract.AnnouncementLevelInfo.String(),
		PublishAt: &publishAt,
		ExpireAt:  &expireAt,
		ActorID:   &actorID,
	}); !errors.Is(err, errAnnouncementInvalidInput) {
		t.Fatalf("expected invalid expire_at guard, got %v", err)
	}

	created, err := service.Create(ctx, announcementstore.CreateInput{
		Title:   "Draft",
		Content: "Draft",
		Level:   announcementcontract.AnnouncementLevelInfo.String(),
		ActorID: &actorID,
	})
	if err != nil {
		t.Fatalf("create draft: %v", err)
	}
	if err := service.Delete(ctx, created.ID, actorID); err != nil {
		t.Fatalf("delete draft: %v", err)
	}
}

func TestAnnouncementManagementServiceListFiltersAndSort(t *testing.T) {
	repository := newMemoryAnnouncementRepository()
	service, err := NewService(repository)
	if err != nil {
		t.Fatalf("new service: %v", err)
	}
	ctx := context.Background()
	create := func(title string, level announcementcontract.AnnouncementLevel, pinned bool) announcementstore.Announcement {
		item, createErr := service.Create(ctx, announcementstore.CreateInput{
			Title:   title,
			Content: "content " + title,
			Level:   level.String(),
			Pinned:  pinned,
		})
		if createErr != nil {
			t.Fatalf("create %s: %v", title, createErr)
		}
		return item
	}
	create("Alpha maintenance", announcementcontract.AnnouncementLevelWarning, true)
	create("Beta release", announcementcontract.AnnouncementLevelInfo, false)
	create("Gamma maintenance", announcementcontract.AnnouncementLevelWarning, false)

	pinned := true
	result, err := service.ListAdmin(ctx, AdminListQuery{
		Level:    announcementcontract.AnnouncementLevelWarning.String(),
		Pinned:   &pinned,
		Keyword:  "maintenance",
		Page:     1,
		PageSize: 10,
		Sort:     "updated_desc",
	})
	if err != nil {
		t.Fatalf("list filtered announcements: %v", err)
	}
	if result.Total != 1 || len(result.Items) != 1 || result.Items[0].Title != "Alpha maintenance" {
		t.Fatalf("unexpected filtered result: total=%d items=%#v", result.Total, result.Items)
	}
}

func TestAnnouncementManagementRoutePermissionDenied(t *testing.T) {
	gin.SetMode(gin.TestMode)
	service, err := NewService(newMemoryAnnouncementRepository())
	if err != nil {
		t.Fatalf("new service: %v", err)
	}
	engine := gin.New()
	ctx := newAnnouncementTestContext(engine)
	if err := registerAnnouncementRoutes(ctx, service, announcementGuards{
		read: func(ginCtx *gin.Context) {
			httpx.AbortLocalizedError(ginCtx, ctx.I18n, http.StatusForbidden, messagecontract.AuthForbidden.String(), nil)
		},
	}); err != nil {
		t.Fatalf("register announcement routes: %v", err)
	}

	request := httptest.NewRequest(http.MethodGet, "/api/announcements", nil)
	recorder := httptest.NewRecorder()
	engine.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusForbidden {
		t.Fatalf("expected 403 permission denial, got %d body=%s", recorder.Code, recorder.Body.String())
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

func (testAnnouncementRepository) ListAdmin(context.Context, announcementstore.ListQuery) (announcementstore.ListResult, error) {
	return announcementstore.ListResult{}, announcementstore.ErrAnnouncementNotFound
}

func (testAnnouncementRepository) Create(context.Context, announcementstore.CreateInput) (announcementstore.Announcement, error) {
	return announcementstore.Announcement{}, announcementstore.ErrAnnouncementNotFound
}

func (testAnnouncementRepository) GetAdmin(context.Context, uint64) (announcementstore.Announcement, error) {
	return announcementstore.Announcement{}, announcementstore.ErrAnnouncementNotFound
}

func (testAnnouncementRepository) Update(context.Context, uint64, announcementstore.UpdateInput) (announcementstore.Announcement, error) {
	return announcementstore.Announcement{}, announcementstore.ErrAnnouncementNotFound
}

func (testAnnouncementRepository) Publish(context.Context, uint64, time.Time, *uint64) (announcementstore.Announcement, error) {
	return announcementstore.Announcement{}, announcementstore.ErrAnnouncementNotFound
}

func (testAnnouncementRepository) Archive(context.Context, uint64, *uint64) (announcementstore.Announcement, error) {
	return announcementstore.Announcement{}, announcementstore.ErrAnnouncementNotFound
}

func (testAnnouncementRepository) Delete(context.Context, uint64, uint64, time.Time) error {
	return announcementstore.ErrAnnouncementNotFound
}

var _ announcementstore.Repository = testAnnouncementRepository{}
