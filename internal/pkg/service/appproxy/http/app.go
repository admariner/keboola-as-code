package http

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/oauth2-proxy/oauth2-proxy/v7/pkg/apis/options"
	"go.opentelemetry.io/otel/attribute"

	"github.com/keboola/keboola-as-code/internal/pkg/service/common/ctxattr"
	"github.com/keboola/keboola-as-code/internal/pkg/service/common/httpserver/middleware"
)

type DataApp struct {
	ID           AppID              `json:"id" validator:"required"`
	Name         string             `json:"name" validator:"required"`
	UpstreamHost string             `json:"upstreamUrl" validator:"required"`
	Providers    []options.Provider `json:"providers"`
	Rules        []Rule             `json:"rules"`
}

type Rule struct {
	Type      RuleType `json:"type"`
	Value     string   `json:"value"`
	Providers []string `json:"providers"`
}
type RuleType string

const PathPrefix = RuleType("pathPrefix")

type AppID string

func (v AppID) String() string {
	return string(v)
}

const attrAppID = "proxy.appid"

func appIDMiddleware(publicURL *url.URL) middleware.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			appID, ok := parseAppID(publicURL, req.Host)

			if ok {
				ctx := req.Context()
				ctx = ctxattr.ContextWith(ctx, attribute.String(attrAppID, string(appID)))
				req = req.WithContext(ctx)
			}

			next.ServeHTTP(w, req)
		})
	}
}

func parseAppID(publicURL *url.URL, host string) (AppID, bool) {
	if !strings.HasSuffix(host, "."+publicURL.Host) {
		return "", false
	}

	if strings.Count(host, ".") != strings.Count(publicURL.Host, ".")+1 {
		return "", false
	}

	idx := strings.IndexByte(host, '.')
	if idx < 0 {
		return "", false
	}

	subdomain := host[:idx]
	idx = strings.LastIndexByte(subdomain, '-')
	if idx < 0 {
		return AppID(subdomain), true
	}

	return AppID(subdomain[idx+1:]), true
}
