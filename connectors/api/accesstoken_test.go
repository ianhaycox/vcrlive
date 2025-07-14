package api

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAccessTokenCache(t *testing.T) {
	t.Run("empty cache returns nil", func(t *testing.T) {
		cache := AccessTokenCache{}

		assert.Nil(t, cache.Get())
	})

	t.Run("get cache returns previous item", func(t *testing.T) {
		cache := AccessTokenCache{}

		cache.Set(AccessToken{Token: "foo", ExpiresSeconds: 10, Type: "bar"})

		token := cache.Get()
		assert.Equal(t, AccessToken{Token: "foo", ExpiresSeconds: 10, Type: "bar"}, *token)

		token2 := cache.Get()
		assert.Equal(t, AccessToken{Token: "foo", ExpiresSeconds: 10, Type: "bar"}, *token2)
	})

	t.Run("set cache returns new item", func(t *testing.T) {
		cache := AccessTokenCache{}

		cache.Set(AccessToken{Token: "foo", ExpiresSeconds: 10, Type: "bar"})

		token := cache.Get()
		assert.Equal(t, AccessToken{Token: "foo", ExpiresSeconds: 10, Type: "bar"}, *token)

		cache.Set(AccessToken{Token: "foo2", ExpiresSeconds: 10, Type: "bar2"})
		token2 := cache.Get()
		assert.Equal(t, AccessToken{Token: "foo2", ExpiresSeconds: 10, Type: "bar2"}, *token2)
	})

	t.Run("get cache returns nil if expired", func(t *testing.T) {
		cache := AccessTokenCache{}

		cache.Set(AccessToken{Token: "short", ExpiresSeconds: 1, Type: "bar"})
		tokenExpired := cache.Get()
		assert.Nil(t, tokenExpired, "because less than grace period")

		cache.Set(AccessToken{Token: "foo", ExpiresSeconds: 6, Type: "bar"})
		token := cache.Get()
		assert.Equal(t, AccessToken{Token: "foo", ExpiresSeconds: 6, Type: "bar"}, *token, "ok because expires 1 second after grace period")

		time.Sleep(time.Duration(2) * time.Second)
		token2 := cache.Get()
		assert.Nil(t, token2, "Even with a few seconds left, consider it expired")

		cache.Set(AccessToken{Token: "foo3", ExpiresSeconds: 10, Type: "bar3"})
		token3 := cache.Get()
		assert.Equal(t, AccessToken{Token: "foo3", ExpiresSeconds: 10, Type: "bar3"}, *token3)
	})
}
