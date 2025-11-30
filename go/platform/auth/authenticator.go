package auth

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

// Authenticator is used to authenticate our users.
type Authenticator struct {
	*oidc.Provider
	oauth2.Config
	HTTPClient *http.Client
}

// New instantiates the *Authenticator.
func New() (*Authenticator, error) {
	transport := &loggingTransport{base: http.DefaultTransport}
	httpClient := &http.Client{Transport: transport}
	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, httpClient)

	provider, err := oidc.NewProvider(
		ctx,
		"https://"+os.Getenv("AUTH0_DOMAIN")+"/",
	)
	if err != nil {
		return nil, err
	}

	conf := oauth2.Config{
		ClientID:     os.Getenv("AUTH0_CLIENT_ID"),
		ClientSecret: os.Getenv("AUTH0_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("AUTH0_CALLBACK_URL"),
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "profile"},
	}

	return &Authenticator{
		Provider:   provider,
		Config:     conf,
		HTTPClient: httpClient,
	}, nil
}

// VerifyIDToken verifies that an *oauth2.Token is a valid *oidc.IDToken.
func (a *Authenticator) VerifyIDToken(ctx context.Context, token *oauth2.Token) (*oidc.IDToken, error) {
	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		return nil, errors.New("no id_token field in oauth2 token")
	}

	oidcConfig := &oidc.Config{
		ClientID: a.ClientID,
	}

	return a.Verifier(oidcConfig).Verify(ctx, rawIDToken)
}

type loggingTransport struct {
	base http.RoundTripper
}

func (lt *loggingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	base := lt.base
	if base == nil {
		base = http.DefaultTransport
	}

	start := time.Now()
	log.Printf("[auth-client] -> %s %s", req.Method, req.URL)

	resp, err := base.RoundTrip(req)
	if err != nil {
		log.Printf("[auth-client] <- error %v (%s)", err, time.Since(start))
		return nil, err
	}

	log.Printf("[auth-client] <- %d %s (%s)", resp.StatusCode, resp.Status, time.Since(start))
	return resp, nil
}
