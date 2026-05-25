package authopenapi

// ServerInterface is the minimal generated handler contract for guarded auth login/bootstrap migration.
type ServerInterface interface {
	PostAuthLogin(params PostAuthLoginParams, body PostAuthLoginJSONRequestBody)
	PostAuthRefresh(params PostAuthRefreshParams)
	PostAuthLogout(params PostAuthLogoutParams)
	GetAuthBootstrap(params GetAuthBootstrapParams)
}
