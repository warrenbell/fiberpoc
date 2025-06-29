package services

import (
	"context"
	"crypto/rand"
	"encoding/base64"

	oidc "github.com/coreos/go-oidc"
	"github.com/pkg/errors"
	"gitlab.com/sandstone2/fiberpoc/common/models"
	"go.uber.org/zap"
	oauth2 "golang.org/x/oauth2"
)

var (
	provider    *oidc.Provider
	oauthConfig *oauth2.Config
	verifier    *oidc.IDTokenVerifier
)

type AuthcServiceInterface interface {
	GetOauthConfig() *oauth2.Config
	GenerateState() (string, error)
	ProcessOauth(code string) (*models.Claims, *string, error)
}

type AuthcService struct {
	logger *zap.Logger
}

func NewAuthcService(logger *zap.Logger) (*AuthcService, error) {
	ctx := context.Background()

	// 1. Initialize OIDC Provider
	var err error
	provider, err = oidc.NewProvider(ctx, "https://accounts.google.com")
	if err != nil {
		return nil, errors.Wrap(err, "Fatal: 3JEUER - Getting oidc provider.")
	}

	// 2. Setup OAuth2 config
	oauthConfig = &oauth2.Config{
		ClientID:     *models.GlobalConfig.GetGoogleOidcClientId(),
		ClientSecret: *models.GlobalConfig.GetGoogleOidcClientSecret(),
		RedirectURL:  *models.GlobalConfig.GetRedirectUri(),
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "email", "profile"},
	}

	// 3. Verifier for the ID Token
	verifier = provider.Verifier(&oidc.Config{ClientID: *models.GlobalConfig.GetGoogleOidcClientId()})

	return &AuthcService{logger: logger}, nil
}

func (authcService *AuthcService) GetOauthConfig() *oauth2.Config {
	return oauthConfig
}

func (authcService *AuthcService) GetVerifier() *oidc.IDTokenVerifier {
	return verifier
}

func (authcService *AuthcService) GenerateState() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", errors.Wrap(err, "Error: Z34I1P - Generating state for oidc.")
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func (authcService *AuthcService) ProcessOauth(code string) (*models.Claims, *string, error) {
	token, err := oauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, nil, errors.Wrap(err, "Error: FN1SF9 - Exchanging the code for a token.")
	}

	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		return nil, nil, errors.New("Error: ZODLPM - Extracting the jwt.")
	}

	idToken, err := verifier.Verify(context.Background(), rawIDToken)
	if err != nil {
		return nil, nil, errors.Wrap(err, "Error: FX6ZJP - Verifying the jwt.")
	}

	// 5. Extract user claims
	claims := &models.Claims{}

	if err := idToken.Claims(claims); err != nil {
		return nil, nil, errors.Wrap(err, "Error: WTWOO1 - Extracting the claims.")
	}

	return claims, &rawIDToken, nil

}
