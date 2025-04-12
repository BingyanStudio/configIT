package auth

import (
	"context"
	"fmt"
	"strings"

	"github.com/BingyanStudio/configIT/internal/model"
	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

var (
	OAuthConfig  *oauth2.Config
	OidcProvider *oidc.Provider
	OidcVerifier *oidc.IDTokenVerifier
)

func InitOIDCConfig(providerUrl string, clientId string, clientSecret string, redirectUrl string, scopes string) error {
	OidcProvider, err := oidc.NewProvider(context.Background(), providerUrl)
	if err != nil {
		return err
	}

	scopeSlice := []string{oidc.ScopeOpenID}
	if scopes != "" {
		for _, scope := range strings.Split(strings.TrimSpace(scopes), " ") {
			if scope != "" {
				scopeSlice = append(scopeSlice, scope)
			}
		}
	}

	OidcVerifier = OidcProvider.Verifier(&oidc.Config{ClientID: clientId})

	OAuthConfig = &oauth2.Config{
		ClientID:     clientId,
		ClientSecret: clientSecret,
		RedirectURL:  redirectUrl,
		Endpoint:     OidcProvider.Endpoint(),
		Scopes:       scopeSlice,
	}
	return nil
}

func handleOAuth2Callback(code string) (string, string, error) {

	oauth2Token, err := OAuthConfig.Exchange(context.Background(), code)
	if err != nil {
		return "", "", err
	}

	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		return "", "", fmt.Errorf("No id_token field in oauth2 token")
	}

	idToken, err := OidcVerifier.Verify(context.Background(), rawIDToken)
	if err != nil {
		return "", "", fmt.Errorf("Failed to verify ID token: %v", err)
	}

	var claims map[string]interface{}
	if err := idToken.Claims(&claims); err != nil {
		return "", "", fmt.Errorf("Failed to parse ID token claims: %v", err)
	}

	subKey, err := model.GetSettings(context.TODO(), "OidcClaimSub")
	if err != nil {
		return "", "", fmt.Errorf("Failed to get OIDC claim sub key: %v", err)
	}
	deptKey, err := model.GetSettings(context.TODO(), "OidcClaimDepartment")
	if err != nil {
		return "", "", fmt.Errorf("Failed to get OIDC claim department key: %v", err)
	}

	return claims[subKey].(string), claims[deptKey].(string), nil
}
