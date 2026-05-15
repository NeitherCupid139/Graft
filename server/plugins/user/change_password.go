package user

import (
	"context"
	"errors"
	"fmt"

	"graft/server/internal/pluginapi"
	"graft/server/internal/store"
)

type changePasswordRequest struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}

// ChangeCurrentUserPassword 在当前登录态下完成一次自助改密，并保留当前会话。
func (s authService) ChangeCurrentUserPassword(ctx context.Context, currentPassword string, newPassword string) error {
	if s.auth == nil {
		return errors.New("auth repository is unavailable")
	}

	requestAuth, ok := pluginapi.RequestAuthContextFromContext(ctx)
	if !ok || requestAuth.User == nil || requestAuth.Claims == nil {
		return pluginapi.ErrUnauthenticated
	}

	credential, err := s.auth.GetUserCredentialByUsername(ctx, requestAuth.User.Username)
	if err != nil {
		if errors.Is(err, store.ErrUserNotFound) {
			return pluginapi.ErrUnauthenticated
		}
		return fmt.Errorf("get current user credential: %w", err)
	}
	if credential.PasswordHash == nil || *credential.PasswordHash == "" {
		return errCurrentPasswordInvalid
	}

	if err := s.passwords.Compare(*credential.PasswordHash, currentPassword); err != nil {
		return errCurrentPasswordInvalid
	}
	if err := s.policy.ValidateNewPassword(currentPassword, newPassword); err != nil {
		return err
	}

	hash, err := s.passwords.Hash(newPassword)
	if err != nil {
		return fmt.Errorf("hash new password: %w", err)
	}
	changedAt := s.nowUTC()
	if err := s.auth.SetPasswordHash(ctx, store.SetPasswordHashInput{
		UserID:             credential.UserID,
		PasswordHash:       hash,
		MustChangePassword: false,
		ChangedAt:          &changedAt,
	}); err != nil {
		return fmt.Errorf("set current user password hash: %w", err)
	}
	if err := s.auth.RevokeOtherRefreshSessionsByUserID(ctx, store.RevokeOtherRefreshSessionsInput{
		UserID:         credential.UserID,
		CurrentTokenID: requestAuth.Claims.SessionID,
		RevokedAt:      changedAt,
	}); err != nil {
		return fmt.Errorf("revoke other refresh sessions after password change: %w", err)
	}

	return nil
}
