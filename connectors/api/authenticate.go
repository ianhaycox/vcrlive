//go:generate mockgen -package api -destination authenticate_mock.go -source authenticate.go
package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/google/go-querystring/query"
	"github.com/ianhaycox/vcrlive/connectors/api/secretsstore"
)

const BaseURLAuthenticationEnv string = "BASE_URL_AUTHENTICATION" // Cognito endpoint to get an Access Token. No trailing /

type AccessDetails struct {
	ClientID     string `url:"client_id" json:"client_id"`
	ClientSecret string `url:"client_secret" json:"client_secret"`
	Scope        string `url:"scope,omitempty" json:"scope,omitempty"`
}

type AuthenticationService struct {
	client    APIClientInterface
	secrets   secretsstore.SecretsStorer
	secretKey string
	cache     AccessTokenCache
	basicAuth *BasicAuth
	secret    string
}

type Authenticator interface {
	GetAccessToken() (*AccessToken, error)
	BasicAuth() (*BasicAuth, error)
	BasicAPIKey() (string, error)
}

func NewAuthenticatorConfiguration(basePath string) *Configuration {
	config := NewConfiguration(basePath)
	config.AddDefaultHeader("Content-Type", "application/x-www-form-urlencoded")

	return config
}

func NewAuthenticationService(client APIClientInterface, secrets secretsstore.SecretsStorer, secretKey string) *AuthenticationService {
	return &AuthenticationService{
		client:    client,
		secrets:   secrets,
		secretKey: secretKey,
		cache:     AccessTokenCache{},
	}
}

// GetAccessToken returns an Access Token from the Authentication service using the client_id and client_secret found via the secretKey
func (a *AuthenticationService) GetAccessToken(authResponse any) (*AccessToken, error) {
	var accessDetails AccessDetails

	if accessToken := a.cache.Get(); accessToken != nil {
		return accessToken, nil
	}

	// Retrieve the client_id/client_secret/scope from the secrets manager
	secret, err := a.secrets.Get(a.secretKey)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(secret), &accessDetails)
	if err != nil {
		return nil, err
	}

	if accessDetails.ClientID == "" || accessDetails.ClientSecret == "" {
		return nil, fmt.Errorf("client credentials can not be blank")
	}

	form, _ := query.Values(accessDetails)
	form.Set("grant_type", "client_credentials")

	/*
	 * Set the base path to be the full URL, e.g.
	 *
	 *     config := NewAuthenticatorConfiguration("https://domain.amazoncognito.com/oauth2/token")
	 *     client := NewAPIClient(config)
	 *     auth := NewAuthenticationService(client, secretStore, "secret-key-in-store")
	 *
	 * then this is re-usable for different OAuth2 servers.
	 */
	request, err := a.client.PrepareRequest(context.TODO(), "", http.MethodPost, url.Values{}, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}

	response, err := a.client.CallAPI(request)
	if err != nil || response == nil {
		return nil, err
	}

	defer BodyClose(response)

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, a.client.ReportError(response, body)
	}

	var accessToken AccessToken

	err = a.client.Decode(&accessToken, body, response.Header.Get("Content-Type"))
	if err != nil {
		return nil, err
	}

	a.cache.Set(accessToken)

	return &accessToken, nil
}

func (a *AuthenticationService) BasicAuth() (*BasicAuth, error) {
	if a.basicAuth != nil {
		return a.basicAuth, nil
	}

	// Retrieve the username:password pair from the secrets manager
	secret, err := a.secrets.Get(a.secretKey)
	if err != nil {
		return nil, err
	}

	var basicAuth BasicAuth

	err = json.Unmarshal([]byte(secret), &basicAuth)
	if err != nil {
		return nil, err
	}

	if basicAuth.UserName == "" || basicAuth.Password == "" {
		return nil, fmt.Errorf("username:password combo can not be blank")
	}

	a.basicAuth = &basicAuth

	return a.basicAuth, nil
}

func (a *AuthenticationService) BasicAPIKey() (string, error) {
	if a.secret != "" {
		return a.secret, nil
	}

	// Retrieve the token from the secrets manager
	secret, err := a.secrets.Get(a.secretKey)
	if err != nil {
		return "", err
	}

	if secret == "" {
		return "", fmt.Errorf("secret can not be blank")
	}

	a.secret = secret

	return a.secret, nil
}
