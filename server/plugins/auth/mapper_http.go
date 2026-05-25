package auth

import "graft/server/internal/pluginapi"

func toLoginResponse(result pluginapi.AuthRefreshResult) loginResponse {
	return loginResponse{
		AccessToken:        result.AccessToken,
		ExpiresAt:          result.AccessExpiry,
		MustChangePassword: result.MustChangePassword,
		User: loginUserResponse{
			ID:          result.User.ID,
			Username:    result.User.Username,
			DisplayName: result.User.DisplayName,
		},
	}
}

func toBootstrapResponse(payload pluginapi.AuthBootstrapPayload) bootstrapResponse {
	menus := make([]bootstrapMenuResponse, 0, len(payload.Menus))
	for _, item := range payload.Menus {
		menus = append(menus, bootstrapMenuResponse{
			Code:       item.Code,
			Title:      item.Title,
			TitleKey:   optionalStringPointer(item.TitleKey),
			Path:       item.Path,
			Icon:       item.Icon,
			Permission: item.Permission,
		})
	}

	return bootstrapResponse{
		User: loginUserResponse{
			ID:          payload.User.ID,
			Username:    payload.User.Username,
			DisplayName: payload.User.DisplayName,
		},
		MustChangePassword: payload.MustChangePassword,
		Roles:              append([]string(nil), payload.Roles...),
		Permissions:        append([]string(nil), payload.Permissions...),
		Menus:              menus,
		Locale: bootstrapLocaleSnapshot{
			CurrentLocale:    payload.Locale.CurrentLocale,
			DefaultLocale:    payload.Locale.DefaultLocale,
			FallbackLocale:   payload.Locale.FallbackLocale,
			SupportedLocales: append([]string(nil), payload.Locale.SupportedLocales...),
		},
	}
}

func optionalStringPointer(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}

func toSessionSummaries(items []pluginapi.AuthSessionSummary) []sessionSummary {
	summaries := make([]sessionSummary, 0, len(items))
	for _, item := range items {
		summaries = append(summaries, sessionSummary{
			SessionID: item.SessionID,
			CreatedAt: item.CreatedAt,
			ExpiresAt: item.ExpiresAt,
			Current:   item.Current,
		})
	}

	return summaries
}
