package api

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ianhaycox/vcrlive/connectors/api/secretsstore"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestAuthenticationServiceHappyPath(t *testing.T) {
	t.Parallel()

	t.Run("Returns an Access Token and no error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		svr := httptest.NewServer(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "/oauth2/token", r.URL.Path)

				form, err := io.ReadAll(r.Body)
				assert.NoError(t, err)
				assert.Equal(t, "client_id=cid&client_secret=secret&grant_type=client_credentials", string(form))

				w.Header().Add("Content-Type", "application/json")
				w.Write([]byte(`{"access_token":"xxxx","expires_in":3600,"token_type":"Bearer"}`))
			}))
		defer svr.Close()

		secrets := secretsstore.NewMockSecretsStorer(ctrl)
		secrets.EXPECT().Get("client-key").Return(`{"client_id":"cid","client_secret":"secret"}`, nil)

		config := NewConfiguration(svr.URL + "/oauth2/token")
		client := NewAPIClient(config)
		svc := NewAuthenticationService(client, secrets, "client-key")
		require.NotNil(t, svc)

		token, err := svc.GetAccessToken(nil)
		assert.NoError(t, err)
		assert.Equal(t, &AccessToken{Token: "xxxx", ExpiresSeconds: 3600, Type: "Bearer"}, token)
	})

	t.Run("Returns an Basic Auth username:password pair and no error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		secrets := secretsstore.NewMockSecretsStorer(ctrl)
		secrets.EXPECT().Get("creds").Return(`{"userName":"user","password":"pass"}`, nil)

		svc := NewAuthenticationService(nil, secrets, "creds")
		require.NotNil(t, svc)

		basicAuth, err := svc.BasicAuth()
		assert.NoError(t, err)
		assert.Equal(t, &BasicAuth{UserName: "user", Password: "pass"}, basicAuth)
	})
}

func TestErrorPaths(t *testing.T) {
	t.Parallel()

	t.Run("returns error if fails to get client credentials from the secrets manager", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		svr := httptest.NewServer(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Add("Content-Type", "application/json")
			}))
		defer svr.Close()

		secrets := secretsstore.NewMockSecretsStorer(ctrl)
		secrets.EXPECT().Get("key").Return("", fmt.Errorf("secret not found"))

		config := NewConfiguration("http://bad request")
		client := NewAPIClient(config)
		svc := NewAuthenticationService(client, secrets, "key")

		_, err := svc.GetAccessToken(nil)
		assert.ErrorContains(t, err, "secret not found")
	})

	t.Run("returns error if fails to unmarshal client credentials from the secrets manager", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		svr := httptest.NewServer(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Add("Content-Type", "application/json")
			}))
		defer svr.Close()

		secrets := secretsstore.NewMockSecretsStorer(ctrl)
		secrets.EXPECT().Get("key").Return("", nil)

		config := NewConfiguration("http://bad request")
		client := NewAPIClient(config)
		svc := NewAuthenticationService(client, secrets, "key")

		_, err := svc.GetAccessToken(nil)
		assert.ErrorContains(t, err, "unexpected end of JSON input")
	})

	t.Run("returns error if prepare request fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		svr := httptest.NewServer(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Add("Content-Type", "application/json")
			}))
		defer svr.Close()

		secrets := secretsstore.NewMockSecretsStorer(ctrl)
		secrets.EXPECT().Get("key").Return(`{"client_id":"cid","client_secret":"secret"}`, nil)

		config := NewConfiguration("http://bad request")
		client := NewAPIClient(config)
		svc := NewAuthenticationService(client, secrets, "key")

		_, err := svc.GetAccessToken(nil)
		assert.ErrorContains(t, err, "bad request")
	})

	t.Run("returns error if the API calls fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		svr := httptest.NewServer(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			}))
		defer svr.Close()

		secrets := secretsstore.NewMockSecretsStorer(ctrl)
		secrets.EXPECT().Get("key").Return(`{"client_id":"cid","client_secret":"secret"}`, nil)

		config := NewConfiguration("")
		client := NewAPIClient(config)
		svc := NewAuthenticationService(client, secrets, "key")

		_, err := svc.GetAccessToken(nil)
		assert.ErrorContains(t, err, "unsupported protocol scheme")
	})

	t.Run("returns error if the client credentials are empty", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		svr := httptest.NewServer(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			}))
		defer svr.Close()

		secrets := secretsstore.NewMockSecretsStorer(ctrl)
		secrets.EXPECT().Get("key").Return(`{"client_id":"","client_secret":""}`, nil)

		config := NewConfiguration("")
		client := NewAPIClient(config)
		svc := NewAuthenticationService(client, secrets, "key")

		_, err := svc.GetAccessToken(nil)
		assert.ErrorContains(t, err, "client credentials can not be blank")
	})

	t.Run("returns error if the http status != 200", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		svr := httptest.NewServer(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusForbidden)
			}))
		defer svr.Close()

		secrets := secretsstore.NewMockSecretsStorer(ctrl)
		secrets.EXPECT().Get("key").Return(`{"client_id":"cid","client_secret":"secret"}`, nil)

		config := NewConfiguration(svr.URL)
		client := NewAPIClient(config)
		svc := NewAuthenticationService(client, secrets, "key")

		_, err := svc.GetAccessToken(nil)
		assert.ErrorContains(t, err, "server returned non-200 http code: 403")
	})

	t.Run("returns error if the payload can not be unmarshalled", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		svr := httptest.NewServer(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Add("Content-Type", "application/json")
			}))
		defer svr.Close()

		secrets := secretsstore.NewMockSecretsStorer(ctrl)
		secrets.EXPECT().Get("key").Return(`{"client_id":"cid","client_secret":"secret"}`, nil)

		config := NewConfiguration(svr.URL)
		client := NewAPIClient(config)
		svc := NewAuthenticationService(client, secrets, "key")

		_, err := svc.GetAccessToken(nil)
		assert.ErrorContains(t, err, "unexpected end of JSON input")
	})

	t.Run("returns error if fails to get username:password from the secrets manager", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		secrets := secretsstore.NewMockSecretsStorer(ctrl)
		secrets.EXPECT().Get("key").Return("", fmt.Errorf("secret not found"))

		svc := NewAuthenticationService(nil, secrets, "key")

		_, err := svc.BasicAuth()
		assert.ErrorContains(t, err, "secret not found")
	})

	t.Run("returns error if fails to unmarshal the username:password", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		secrets := secretsstore.NewMockSecretsStorer(ctrl)
		secrets.EXPECT().Get("key").Return("", nil)

		svc := NewAuthenticationService(nil, secrets, "key")

		_, err := svc.BasicAuth()
		assert.ErrorContains(t, err, "unexpected end of JSON input")
	})

	t.Run("returns error if username or password blank", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		secrets := secretsstore.NewMockSecretsStorer(ctrl)
		secrets.EXPECT().Get("blankuser").Return(`{"userName":"","password":"pass"}`, nil)
		secrets.EXPECT().Get("blankpass").Return(`{"userName":"user","password":""}`, nil)
		secrets.EXPECT().Get("blankboth").Return(`{}`, nil)

		svc := NewAuthenticationService(nil, secrets, "blankuser")
		_, err := svc.BasicAuth()
		assert.ErrorContains(t, err, "username:password combo can not be blank")

		svc2 := NewAuthenticationService(nil, secrets, "blankpass")
		_, err = svc2.BasicAuth()
		assert.ErrorContains(t, err, "username:password combo can not be blank")

		svc3 := NewAuthenticationService(nil, secrets, "blankboth")
		_, err = svc3.BasicAuth()
		assert.ErrorContains(t, err, "username:password combo can not be blank")
	})
}

func TestAuthenticationConfiguration(t *testing.T) {
	t.Parallel()

	t.Run("Returns a form url encoded header", func(t *testing.T) {
		config := NewAuthenticatorConfiguration("https://example.com")

		assert.Equal(t, config.BasePath, "https://example.com")
		assert.Equal(t, config.DefaultHeader, map[string]string{"Content-Type": "application/x-www-form-urlencoded", "Accept": "application/json"})
	})
}
