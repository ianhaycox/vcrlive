package api

import (
	"time"
)

const (
	GracePeriod int = 5
)

type AccessToken struct {
	Token          string `json:"access_token,omitempty"`
	ExpiresSeconds int    `json:"expires_in,omitempty"`
	Type           string `json:"token_type,omitempty"`
}

type AccessTokenCache struct {
	accessToken *AccessToken
	expires     time.Time
}

func (a *AccessTokenCache) Set(accessToken AccessToken) {
	a.accessToken = &accessToken
	a.expires = time.Now().UTC().Add(time.Duration(accessToken.ExpiresSeconds-GracePeriod) * time.Second)
}

func (a *AccessTokenCache) Get() *AccessToken {
	if a.accessToken == nil {
		return nil
	}

	if time.Now().UTC().After(a.expires) {
		return nil
	}

	return a.accessToken
}
