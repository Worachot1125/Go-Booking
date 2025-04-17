package provider

import (
	"sync"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	googleOauthConfig *oauth2.Config
	oauthOnce         sync.Once
)

// OAuthOption contains the configuration details for OAuth
type OAuthOption struct {
	RedirectURL  string
	ClientID     string
	ClientSecret string
	Scopes       []string
}

// RegisterOAuth initializes and registers the Google OAuth2 configuration
func RegisterOAuth(opt *OAuthOption) {
	oauthOnce.Do(func() {
		googleOauthConfig = &oauth2.Config{
			RedirectURL:  opt.RedirectURL,
			ClientID:     opt.ClientID,
			ClientSecret: opt.ClientSecret,
			Scopes:       opt.Scopes,
			Endpoint:     google.Endpoint,
		}
	})
}

// GetGoogleOAuthConfig returns the registered Google OAuth2 configuration
func GetGoogleOAuthConfig() *oauth2.Config {
	if googleOauthConfig == nil {
		panic("OAuth configuration not initialized")
	}
	return googleOauthConfig
}
