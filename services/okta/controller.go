package okta

import (
	"time"

	"github.com/go-redis/redis"
	"golang.org/x/sync/singleflight"
)

var requestGroup singleflight.Group

// GetUser returns an Okta user by email. It first tries to get it from cache,
// and if not present there, it will fetch it from Okta API.
func GetUser(email string) (User, error) {
	user, err := CacheGet(email)
	if err == nil {
		// Cache hit
		return user, nil
	}

	if err != redis.Nil {
		// Not a cache hit, not a cache miss, something went wrong
		return User{}, err
	}

	// Cache miss
	// Deduplicate network calls and cache writes if this controller is called
	// multiple times concurrently.
	val, err, _ := requestGroup.Do(email, func() (interface{}, error) {
		user, err := FetchUser(email)
		if err != nil {
			return User{}, err
		}

		err = CacheSet(user.Email, user, time.Minute*10)
		return user, err
	})

	if err != nil {
		return User{}, err
	}
	return val.(User), nil
}
