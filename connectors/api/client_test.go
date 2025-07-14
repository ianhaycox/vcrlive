package api

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"golang.org/x/oauth2"
)

type testResult struct {
	ID string `json:"id,omitempty" xml:"id,omitempty"`
}

type testToken struct{}

func (t *testToken) Token() (*oauth2.Token, error) {
	return &oauth2.Token{AccessToken: "test-token"}, nil
}

func TestPrepareRequest(t *testing.T) {
	t.Parallel()

	const testBasePath = "http://localhost"

	t.Run("return error if URL path cannot be parsed", func(t *testing.T) {
		api := NewAPIClient(NewConfiguration("http://a b"))

		_, err := api.PrepareRequest(context.TODO(), "/", http.MethodGet, url.Values{}, nil)
		assert.Error(t, err)
	})

	t.Run("return no error if GET request OK", func(t *testing.T) {
		api := NewAPIClient(NewConfiguration("http://bookings"))

		queryParams := url.Values{}
		queryParams.Add("foo", "bar")
		_, err := api.PrepareRequest(context.TODO(), "/", http.MethodGet, queryParams, nil)
		assert.NoError(t, err)
	})

	t.Run("return no error if POST request OK", func(t *testing.T) {
		api := NewAPIClient(NewConfiguration("http://bookings"))

		queryParams := url.Values{}
		queryParams.Add("foo", "bar")
		_, err := api.PrepareRequest(context.TODO(), "/", http.MethodPost, queryParams, testToken{})
		assert.NoError(t, err)
	})

	t.Run("Check OAuth2 authentication added", func(t *testing.T) {
		api := NewAPIClient(NewConfiguration(testBasePath))
		ctx := context.WithValue(context.Background(), ContextOAuth2, &testToken{})

		request, err := api.PrepareRequest(ctx, "/", http.MethodGet, url.Values{}, nil)
		assert.NoError(t, err)
		assert.Equal(t, "Bearer test-token", request.Header.Get("Authorization"))
	})

	t.Run("Check BasicAuth authentication", func(t *testing.T) {
		api := NewAPIClient(NewConfiguration(testBasePath))
		ctx := context.WithValue(context.Background(), ContextBasicAuth, BasicAuth{UserName: "test", Password: "pass"})

		request, err := api.PrepareRequest(ctx, "/", http.MethodGet, url.Values{}, nil)
		assert.NoError(t, err)
		assert.Equal(t, "Basic dGVzdDpwYXNz", request.Header.Get("Authorization"))
	})

	t.Run("Check BasicAPIKey authentication", func(t *testing.T) {
		api := NewAPIClient(NewConfiguration(testBasePath))
		ctx := context.WithValue(context.Background(), ContextBasicAPIKey, BasicAPIKey{Key: "dGVzdDpwYXNz"})

		request, err := api.PrepareRequest(ctx, "/", http.MethodGet, url.Values{}, nil)
		assert.NoError(t, err)
		assert.Equal(t, "Basic dGVzdDpwYXNz", request.Header.Get("Authorization"))
	})

	t.Run("Check AccessToken authentication added", func(t *testing.T) {
		api := NewAPIClient(NewConfiguration(testBasePath))
		ctx := context.WithValue(context.Background(), ContextAccessToken, "access-token")

		request, err := api.PrepareRequest(ctx, "/", http.MethodGet, url.Values{}, nil)
		assert.NoError(t, err)
		assert.Equal(t, "Bearer access-token", request.Header.Get("Authorization"))
	})

	t.Run("Check ContextAPIKey authentication added", func(t *testing.T) {
		api := NewAPIClient(NewConfiguration(testBasePath))
		ctx := context.WithValue(context.Background(), ContextAPIKey, APIKey{Key: "api-key", Prefix: "prefix"})

		request, err := api.PrepareRequest(ctx, "/", http.MethodGet, url.Values{}, nil)
		assert.NoError(t, err)
		assert.Equal(t, "prefix api-key", request.Header.Get("X-API-KEY"))
	})

	t.Run("Check ContextAPIKey without prefix authentication added", func(t *testing.T) {
		api := NewAPIClient(NewConfiguration(testBasePath))
		ctx := context.WithValue(context.Background(), ContextAPIKey, APIKey{Key: "api-key"})

		request, err := api.PrepareRequest(ctx, "/", http.MethodGet, url.Values{}, nil)
		assert.NoError(t, err)
		assert.Equal(t, "api-key", request.Header.Get("X-API-KEY"))
	})

	t.Run("Check ContextCustomAuth authentication", func(t *testing.T) {
		api := NewAPIClient(NewConfiguration(testBasePath))
		ctx := context.WithValue(context.Background(), ContextCustomAuth, "custom dGVzdDpwYXNz")

		request, err := api.PrepareRequest(ctx, "/", http.MethodGet, url.Values{}, nil)
		assert.NoError(t, err)
		assert.Equal(t, "custom dGVzdDpwYXNz", request.Header.Get("Authorization"))
	})
}

func TestDecode(t *testing.T) {
	t.Parallel()

	t.Run("decodes JSON successfully", func(t *testing.T) {
		api := NewAPIClient(NewConfiguration(""))

		result := testResult{}
		err := api.Decode(&result, []byte(`{"id":"test"}`), "application/json")

		assert.NoError(t, err)
		assert.Equal(t, testResult{ID: "test"}, result)
	})

	t.Run("returns error for unsuccessful JSON decode", func(t *testing.T) {
		api := NewAPIClient(NewConfiguration(""))

		result := testResult{}
		err := api.Decode(&result, []byte(`{"id":test}`), "application/json")

		assert.Error(t, err)
	})

	t.Run("decodes XML successfully", func(t *testing.T) {
		api := NewAPIClient(NewConfiguration(""))

		result := testResult{}
		err := api.Decode(&result, []byte(`<xml><id>test</id></xml>`), "application/xml")

		assert.NoError(t, err)
		assert.Equal(t, testResult{ID: "test"}, result)
	})

	t.Run("returns error for unsuccessful XML decode", func(t *testing.T) {
		api := NewAPIClient(NewConfiguration(""))

		result := testResult{}
		err := api.Decode(&result, []byte(`<xml><id>test</notid></xml>`), "application/xml")

		assert.Error(t, err)
	})

	t.Run("returns error for unknown content type", func(t *testing.T) {
		api := NewAPIClient(NewConfiguration(""))

		result := testResult{}
		err := api.Decode(&result, []byte(``), "unknown")

		assert.Error(t, err)
	})
}

func TestReportError(t *testing.T) {
	t.Run("should return error if the api call failed", func(t *testing.T) {
		api := NewAPIClient(NewConfiguration(""))

		err := api.ReportError(&http.Response{StatusCode: http.StatusBadRequest}, []byte(""))
		assert.ErrorContains(t, err, "server returned non-200")
		assert.ErrorContains(t, err, "400, response ''")
	})

	t.Run("should return error from response if the api call failed", func(t *testing.T) {
		api := NewAPIClient(NewConfiguration(""))

		err := api.ReportError(&http.Response{Header: http.Header{"Content-Type": []string{"application/json"}}, StatusCode: http.StatusBadRequest}, []byte(`{"message":"server error"}`))
		assert.ErrorContains(t, err, "server returned non-200 http code: 400")
		assert.ErrorContains(t, err, "server error")
	})
}

func TestApplyAccessToken(t *testing.T) {
	t.Run("should return an access token and no error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.TODO()
		authenticator := NewMockAuthenticator(ctrl)
		authenticator.EXPECT().GetAccessToken().Return(&AccessToken{Token: "foo", ExpiresSeconds: 10, Type: "Bearer"}, nil)

		ctx, err := ApplyAccessToken(ctx, authenticator)
		assert.NoError(t, err)
		token, ok := ctx.Value(ContextAccessToken).(string)
		assert.True(t, ok)
		assert.Equal(t, "foo", token)
	})

	t.Run("should return an error if can not get the access token", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.TODO()
		authenticator := NewMockAuthenticator(ctrl)
		authenticator.EXPECT().GetAccessToken().Return(nil, fmt.Errorf("can not get token"))

		ctx, err := ApplyAccessToken(ctx, authenticator)
		assert.ErrorContains(t, err, "can not get token")
		_, ok := ctx.Value(ContextAccessToken).(string)
		assert.False(t, ok)
	})

	t.Run("if the token is already present in the context don't get it again", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.TODO()
		ctx = context.WithValue(ctx, ContextAccessToken, "exists")

		ctx, err := ApplyAccessToken(ctx, nil)
		assert.NoError(t, err)
		token, ok := ctx.Value(ContextAccessToken).(string)
		assert.True(t, ok)
		assert.Equal(t, "exists", token)
	})
}

func TestApplyBasicAuth(t *testing.T) {
	t.Run("should return a username:password and no error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.TODO()
		authenticator := NewMockAuthenticator(ctrl)
		authenticator.EXPECT().BasicAuth().Return(&BasicAuth{UserName: "foo", Password: "bar"}, nil)

		ctx, err := ApplyBasicAuth(ctx, authenticator)
		assert.NoError(t, err)
		basicAuth, ok := ctx.Value(ContextBasicAuth).(BasicAuth)
		assert.True(t, ok)
		assert.Equal(t, BasicAuth{UserName: "foo", Password: "bar"}, basicAuth)
	})

	t.Run("should return an error if can not get the username:password", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.TODO()
		authenticator := NewMockAuthenticator(ctrl)
		authenticator.EXPECT().BasicAuth().Return(nil, fmt.Errorf("can not get secret"))

		ctx, err := ApplyBasicAuth(ctx, authenticator)
		assert.ErrorContains(t, err, "can not get secret")
		_, ok := ctx.Value(ContextBasicAuth).(BasicAuth)
		assert.False(t, ok)
	})

	t.Run("if the token is already present in the context don't get it again", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.TODO()
		ctx = context.WithValue(ctx, ContextBasicAuth, BasicAuth{UserName: "exists", Password: "x"})

		ctx, err := ApplyBasicAuth(ctx, nil)
		assert.NoError(t, err)
		basicAuth, ok := ctx.Value(ContextBasicAuth).(BasicAuth)
		assert.True(t, ok)
		assert.Equal(t, BasicAuth{UserName: "exists", Password: "x"}, basicAuth)
	})
}
