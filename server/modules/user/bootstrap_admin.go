// Copyright (c) 2025-2026 GeWuYou
// SPDX-License-Identifier: Apache-2.0

package user

import (
	"context"
	"errors"
	"fmt"

	"graft/server/internal/i18n"
	"golang.org/x/crypto/bcrypt"

	"graft/server/internal/moduleapi"
	"graft/server/internal/permission"
	userstore "graft/server/modules/user/store"
)

// ensureDefaultAdmin 幂等确保默认管理员存在且具备当前 MVP 所需的最小后台可见性。
func (s authService) ensureDefaultAdmin(
	ctx context.Context,
	localizer *i18n.Service,
	rbac moduleapi.RBACBootstrapService,
	permissions []permission.Item,
) error {
	if s.auth == nil {
		return fmt.Errorf("auth repository is unavailable")
	}
	if rbac == nil {
		return fmt.Errorf("rbac bootstrap service is unavailable")
	}

	credential, err := s.ensureAdminCredential(ctx)
	if err != nil {
		return err
	}
	return rbac.EnsureDefaultAdminAccess(ctx, credential.UserID, permissionSeedsFromItems(localizer, permissions))
}

func (s authService) ensureAdminCredential(ctx context.Context) (userstore.UserCredential, error) {
	credential, err := s.auth.GetUserCredentialByUsername(ctx, defaultAdminUsername)
	if err == nil {
		return s.reconcileDefaultAdminCredential(ctx, credential)
	}
	if !errors.Is(err, userstore.ErrUserNotFound) {
		return userstore.UserCredential{}, fmt.Errorf("get default admin credential: %w", err)
	}

	hash, hashErr := s.passwords.Hash(defaultAdminPassword)
	if hashErr != nil {
		return userstore.UserCredential{}, fmt.Errorf("hash default admin password: %w", hashErr)
	}

	credential, err = s.auth.EnsureUserCredential(ctx, userstore.EnsureUserCredentialInput{
		Username:           defaultAdminUsername,
		Display:            defaultAdminDisplay,
		PasswordHash:       hash,
		MustChangePassword: true,
	})
	if err != nil {
		return userstore.UserCredential{}, fmt.Errorf("ensure default admin credential: %w", err)
	}

	return credential, nil
}

func (s authService) reconcileDefaultAdminCredential(
	ctx context.Context,
	credential userstore.UserCredential,
) (userstore.UserCredential, error) {
	if credential.MustChangePassword || credential.PasswordHash == nil || *credential.PasswordHash == "" {
		return credential, nil
	}

	if err := s.passwords.Compare(*credential.PasswordHash, defaultAdminPassword); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return credential, nil
		}
		return userstore.UserCredential{}, fmt.Errorf("compare default admin password hash: %w", err)
	}

	if err := s.auth.SetPasswordHash(ctx, userstore.SetPasswordHashInput{
		UserID:             credential.UserID,
		PasswordHash:       *credential.PasswordHash,
		MustChangePassword: true,
		ChangedAt:          credential.PasswordChangedAt,
	}); err != nil {
		return userstore.UserCredential{}, fmt.Errorf("mark default admin credential for password change: %w", err)
	}

	credential.MustChangePassword = true
	return credential, nil
}

func permissionSeedsFromItems(localizer *i18n.Service, items []permission.Item) []moduleapi.PermissionSeed {
	seeds := make([]moduleapi.PermissionSeed, 0, len(items))
	for _, item := range items {
		seeds = append(seeds, moduleapi.PermissionSeed{
			Code:           item.Code,
			Display:        lookupPermissionText(localizer, item.DisplayKey, item.Code),
			DisplayKey:     item.DisplayKey,
			Description:    lookupPermissionText(localizer, item.DescriptionKey, item.Code),
			DescriptionKey: item.DescriptionKey,
			Category:       item.Category,
		})
	}

	return seeds
}

func lookupPermissionText(localizer *i18n.Service, key string, permissionCode string) string {
	if localizer == nil {
		panic("permission seed localization requires i18n service")
	}
	if key == "" {
		panic("permission seed localization requires stable locale key for " + permissionCode)
	}
	if len(localizer.RegisteredMessageResources(i18n.LocaleTag(localizer.DefaultLocale()), i18n.MessageKey(key))) == 0 {
		panic("permission seed localization key missing for " + permissionCode + ": " + key)
	}

	return localizer.Lookup(i18n.LookupRequest{
		Locale: i18n.LocaleTag(localizer.DefaultLocale()),
		Key:    i18n.MessageKey(key),
	})
}
