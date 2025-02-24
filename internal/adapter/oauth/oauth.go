package oauth

import (
	"backend-layout/internal/config"
	"context"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type Oauth struct {
	googleOAuthConfig *oauth2.Config
}

func NewOauth(conf config.OauthConfig) *Oauth {

	googleAuthConfig := &oauth2.Config{
		ClientID:     conf.GoogleClientID,
		ClientSecret: conf.GoogleClientSecret,
		RedirectURL:  conf.GoogleRedirectUrl,
		Endpoint:     google.Endpoint,
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
	}

	return &Oauth{
		googleOAuthConfig: googleAuthConfig,
	}
}

func (o *Oauth) ExhangeCodeForToken(ctx context.Context, code string) (*oauth2.Token, error) {
	return o.googleOAuthConfig.Exchange(ctx, code)
}

func (o *Oauth) AuthUrlGoogleLogin(state string) (string, error) {
	return o.googleOAuthConfig.AuthCodeURL(state), nil
}
