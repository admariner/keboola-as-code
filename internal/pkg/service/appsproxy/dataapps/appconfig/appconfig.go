package appconfig

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/benbjohnson/clock"

	"github.com/keboola/keboola-as-code/internal/pkg/log"
	"github.com/keboola/keboola-as-code/internal/pkg/service/appsproxy/dataapps/api"
	"github.com/keboola/keboola-as-code/internal/pkg/service/appsproxy/syncmap"
	"github.com/keboola/keboola-as-code/internal/pkg/utils/errors"
)

// staleCacheFallbackDuration is the maximum duration for which the old configuration of an application is used if loading new configuration is not possible.
const staleCacheFallbackDuration = time.Hour

type Loader struct {
	clock  clock.Clock
	logger log.Logger
	api    *api.API
	cache  *syncmap.SyncMap[api.AppID, cachedAppProxyConfig]
}

type cachedAppProxyConfig struct {
	lock      *sync.Mutex
	config    api.AppConfig
	expiresAt time.Time
}

type dependencies interface {
	Clock() clock.Clock
	Logger() log.Logger
	AppsAPI() *api.API
}

func NewLoader(d dependencies) *Loader {
	return &Loader{
		clock:  d.Clock(),
		logger: d.Logger(),
		api:    d.AppsAPI(),
		cache: syncmap.New[api.AppID, cachedAppProxyConfig](func() *cachedAppProxyConfig {
			return &cachedAppProxyConfig{
				lock: &sync.Mutex{},
			}
		}),
	}
}

// GetConfig gets the AppConfig by the ID from Sandboxes Service.
// It handles local caching based on the Cache-Control and ETag headers.
func (l *Loader) GetConfig(ctx context.Context, appID api.AppID) (out api.AppConfig, modified bool, err error) {
	// Get cache item or init an empty item
	item := l.cache.GetOrInit(appID)

	// Only one update runs in parallel.
	// If there is an in-flight update, we are waiting for its results.
	item.lock.Lock()
	defer item.lock.Unlock()

	// Return config from cache if it is still valid.
	// At first, the item.expiresAt is zero, so the condition is skipped.
	now := l.clock.Now()
	if now.Before(item.expiresAt) {
		return item.config, false, nil
	}

	// Send API request with cached eTag.
	// At first, the item.config.ETag() is empty string.
	newConfig, err := l.api.GetAppConfig(appID, item.config.ETag()).Send(ctx)
	if err != nil {
		// The config hasn't been modified, extend expiration, return cached version
		notModifierErr := api.NotModifiedError{}
		if errors.As(err, &notModifierErr) {
			item.ExtendExpiration(now, notModifierErr.MaxAge)
			return item.config, false, nil
		}

		// Only the not found error is expected
		var apiErr *api.Error
		if errors.As(err, &apiErr) && apiErr.StatusCode() != http.StatusNotFound {
			// Log other errors
			l.logger.Errorf(ctx, `failed loading config for app "%s": %s`, appID, err.Error())

			// Keep the proxy running for a limited time in case of an API outage.
			// The item.expiresAt may be zero, if there is no cached version, then the condition is skipped.
			if now.Before(item.expiresAt.Add(staleCacheFallbackDuration)) {
				l.logger.Warnf(ctx, `using stale cache for app "%s": %s`, appID, err.Error())
				return item.config, false, nil
			}
		}

		// Return the error if:
		//  - It is not found error.
		//  - There is no cached version.
		//  - The staleCacheFallbackDuration has been exceeded.
		return api.AppConfig{}, false, err
	}

	// Cache the loaded configuration
	item.config = *newConfig
	item.ExtendExpiration(now, item.config.MaxAge())
	return item.config, true, nil
}

func (v *cachedAppProxyConfig) ExtendExpiration(now time.Time, maxAge time.Duration) {
	v.expiresAt = now.Add(maxAge)
}
