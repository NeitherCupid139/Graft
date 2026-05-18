package user

import (
	"context"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"graft/server/internal/permission"
	"graft/server/internal/pluginapi"
	internalstore "graft/server/internal/store"
	userstore "graft/server/plugins/user/store"
)

// ensureDefaultAdmin 幂等确保默认管理员存在且具备当前 MVP 所需的最小后台可见性。
func (s authService) ensureDefaultAdmin(
	ctx context.Context,
	rbac pluginapi.RBACBootstrapService,
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
	return rbac.EnsureDefaultAdminAccess(ctx, credential.UserID, permissionSeedsFromItems(permissions))
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

func ensureRolePermissions(
	ctx context.Context,
	rbac internalstore.RBACRepository,
	roleID uint64,
	permissions []permission.Item,
) error {
	permissionIDs := make([]uint64, 0, len(permissions))
	for _, item := range permissions {
		record, err := rbac.EnsurePermission(ctx, internalstore.EnsurePermissionInput{
			Code:        item.Code,
			Display:     item.Name,
			Description: stringPtrOrNil(item.Description),
			Category:    item.Category,
		})
		if err != nil {
			return fmt.Errorf("ensure permission %s: %w", item.Code, err)
		}
		permissionIDs = append(permissionIDs, record.ID)
	}
	if len(permissionIDs) == 0 {
		return nil
	}

	if err := rbac.AssignPermissionsToRole(ctx, internalstore.AssignPermissionsToRoleInput{
		RoleID:        roleID,
		PermissionIDs: permissionIDs,
	}); err != nil {
		return fmt.Errorf("assign permissions to default admin role: %w", err)
	}

	return nil
}

func ensureDefaultAdminAccess(
	ctx context.Context,
	rbac internalstore.RBACRepository,
	userID uint64,
	permissions []permission.Item,
) error {
	role, err := rbac.EnsureRole(ctx, internalstore.EnsureRoleInput{
		Name:    defaultAdminRoleName,
		Display: "管理员",
		Builtin: true,
	})
	if err != nil {
		return fmt.Errorf("ensure default admin role: %w", err)
	}

	if err := ensureRolePermissions(ctx, rbac, role.ID, permissions); err != nil {
		return err
	}
	if err := rbac.AssignRoleToUser(ctx, internalstore.AssignRoleToUserInput{
		UserID: userID,
		RoleID: role.ID,
	}); err != nil {
		return fmt.Errorf("assign default admin role to user: %w", err)
	}

	return nil
}

func stringPtrOrNil(value string) *string {
	if value == "" {
		return nil
	}
	result := value
	return &result
}

func permissionSeedsFromItems(items []permission.Item) []pluginapi.PermissionSeed {
	seeds := make([]pluginapi.PermissionSeed, 0, len(items))
	for _, item := range items {
		seeds = append(seeds, pluginapi.PermissionSeed{
			Code:        item.Code,
			Display:     item.Name,
			Description: item.Description,
			Category:    item.Category,
		})
	}

	return seeds
}

type repositoryBackedRBACBootstrapService struct {
	rbac internalstore.RBACRepository
}

func (s repositoryBackedRBACBootstrapService) EnsureDefaultAdminAccess(
	ctx context.Context,
	userID uint64,
	permissions []pluginapi.PermissionSeed,
) error {
	if s.rbac == nil {
		return errors.New("rbac repository is unavailable")
	}

	items := make([]permission.Item, 0, len(permissions))
	for _, permissionSeed := range permissions {
		items = append(items, permission.Item{
			Code:        permissionSeed.Code,
			Name:        permissionSeed.Display,
			Description: permissionSeed.Description,
			Category:    permissionSeed.Category,
		})
	}

	return ensureDefaultAdminAccess(ctx, s.rbac, userID, items)
}

var _ pluginapi.RBACBootstrapService = repositoryBackedRBACBootstrapService{}
