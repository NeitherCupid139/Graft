package store

import (
	"context"
	"errors"
	"time"
)

// ErrRefreshSessionNotFound 表示请求的刷新会话不存在。
var ErrRefreshSessionNotFound = errors.New("refresh session not found")

// UserCredential 表示认证路径可见的最小用户口令 DTO。
//
// 该 DTO 与 User 分离，避免普通用户资料读取能力意外获得密码散列等敏感字段。
type UserCredential struct {
	UserID            uint64
	Username          string
	PasswordHash      *string
	PasswordChangedAt *time.Time
}

// SetPasswordHashInput 描述一次密码散列更新所需的最小输入。
type SetPasswordHashInput struct {
	UserID       uint64
	PasswordHash string
	ChangedAt    time.Time
}

// RefreshSession 表示 refresh token 生命周期对应的稳定持久化 DTO。
type RefreshSession struct {
	ID                uint64
	UserID            uint64
	TokenID           string
	ExpiresAt         time.Time
	RevokedAt         *time.Time
	ReplacedByTokenID *string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

// CreateRefreshSessionInput 描述一次刷新会话创建所需的最小输入。
type CreateRefreshSessionInput struct {
	UserID    uint64
	TokenID   string
	ExpiresAt time.Time
}

// RevokeRefreshSessionInput 描述一次刷新会话吊销所需的最小输入。
type RevokeRefreshSessionInput struct {
	TokenID           string
	RevokedAt         time.Time
	ReplacedByTokenID *string
}

// AuthRepository 暴露未来认证插件所需的最小持久化操作集合。
//
// 该接口只提供口令与 refresh session 的存储能力，不承载登录、签发或授权决策。
type AuthRepository interface {
	// GetUserCredentialByUsername 按用户名读取口令校验所需的最小用户信息。
	//
	// 未命中时统一返回 ErrUserNotFound。
	GetUserCredentialByUsername(ctx context.Context, username string) (UserCredential, error)

	// SetPasswordHash 为指定用户写入口令散列及其最近变更时间。
	//
	// 当用户不存在时统一返回 ErrUserNotFound。
	SetPasswordHash(ctx context.Context, input SetPasswordHashInput) error

	// CreateRefreshSession 持久化一条新的刷新会话记录。
	CreateRefreshSession(ctx context.Context, input CreateRefreshSessionInput) (RefreshSession, error)

	// GetRefreshSessionByTokenID 按 token 标识读取刷新会话状态。
	//
	// 未命中时统一返回 ErrRefreshSessionNotFound。
	GetRefreshSessionByTokenID(ctx context.Context, tokenID string) (RefreshSession, error)

	// RevokeRefreshSession 吊销一条刷新会话，并可选记录轮换后的新 token 标识。
	//
	// 未命中时统一返回 ErrRefreshSessionNotFound。
	RevokeRefreshSession(ctx context.Context, input RevokeRefreshSessionInput) error
}
