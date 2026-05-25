package auth

import (
	"time"

	generated "graft/server/internal/contract/openapi/generated"
)

type loginUserResponse struct {
	ID          uint64 `json:"id"`
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
}

type loginResponse struct {
	AccessToken        string            `json:"access_token"`
	ExpiresAt          time.Time         `json:"expires_at"`
	MustChangePassword bool              `json:"must_change_password"`
	User               loginUserResponse `json:"user"`
}

type bootstrapResponse struct {
	User               loginUserResponse       `json:"user"`
	MustChangePassword bool                    `json:"must_change_password"`
	Roles              []string                `json:"roles"`
	Permissions        []string                `json:"permissions"`
	Menus              []bootstrapMenuResponse `json:"menus"`
	Locale             bootstrapLocaleSnapshot `json:"locale"`
}

type bootstrapMenuResponse = generated.BootstrapMenu
type bootstrapLocaleSnapshot = generated.BootstrapLocale
type sessionSummary struct {
	SessionID string    `json:"session_id"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
	Current   bool      `json:"current"`
}
