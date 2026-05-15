package user

import (
	"errors"
	"strings"
	"unicode"
)

const (
	defaultAdminUsername = "graft"
	defaultAdminDisplay  = "Graft Admin"
	defaultAdminPassword = "graft-admin"
	defaultAdminRoleName = "admin"
)

type passwordPolicy struct{}

func newPasswordPolicy() passwordPolicy {
	return passwordPolicy{}
}

// ValidateNewPassword 校验一次新密码是否满足当前 MVP 固定策略。
func (passwordPolicy) ValidateNewPassword(currentPassword string, newPassword string) error {
	if newPassword == currentPassword {
		return errPasswordReuseForbidden
	}
	if newPassword == defaultAdminPassword {
		return errPasswordReuseForbidden
	}
	if len(newPassword) < 12 {
		return errPasswordPolicyViolation
	}

	var hasLetter bool
	var hasDigit bool
	for _, r := range newPassword {
		if unicode.IsLetter(r) {
			hasLetter = true
		}
		if unicode.IsDigit(r) {
			hasDigit = true
		}
	}
	if !hasLetter || !hasDigit {
		return errPasswordPolicyViolation
	}

	return nil
}

// ValidateDefaultAdminPasswordGuard 验证初始化路径之外是否误用默认管理员例外密码。
func (passwordPolicy) ValidateDefaultAdminPasswordGuard(password string) error {
	if strings.TrimSpace(password) == "" {
		return errPasswordRequired
	}
	if password == defaultAdminPassword {
		return errors.New("default admin password is reserved for bootstrap only")
	}
	return nil
}
