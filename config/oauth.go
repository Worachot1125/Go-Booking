package config

import (
	provider "app/app/provider/OAuth"
	"app/internal/logger"
)

// OAuth initializes the OAuth2 configuration for Google
func OAuth() {
	provider.RegisterOAuth(&provider.OAuthOption{
		RedirectURL:  confString("REDIRECT_URL", ""),
		ClientID:     confString("CLIENT_ID", ""),
		ClientSecret: confString("CLIENT_SECRET", ""),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.profile", "https://www.googleapis.com/auth/userinfo.email"},
	})
	logger.Infof("OAuth configuration initialized")

}
