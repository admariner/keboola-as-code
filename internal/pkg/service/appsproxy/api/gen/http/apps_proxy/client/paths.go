// Code generated by goa v3.19.1, DO NOT EDIT.
//
// HTTP request path constructors for the apps-proxy service.
//
// Command:
// $ goa gen github.com/keboola/keboola-as-code/api/appsproxy --output
// ./internal/pkg/service/appsproxy/api

package client

// APIRootIndexAppsProxyPath returns the URL path to the apps-proxy service ApiRootIndex HTTP endpoint.
func APIRootIndexAppsProxyPath() string {
	return "/_proxy/api/"
}

// APIVersionIndexAppsProxyPath returns the URL path to the apps-proxy service ApiVersionIndex HTTP endpoint.
func APIVersionIndexAppsProxyPath() string {
	return "/_proxy/api/v1"
}

// HealthCheckAppsProxyPath returns the URL path to the apps-proxy service HealthCheck HTTP endpoint.
func HealthCheckAppsProxyPath() string {
	return "/_proxy/api/v1/health-check"
}

// ValidateAppsProxyPath returns the URL path to the apps-proxy service Validate HTTP endpoint.
func ValidateAppsProxyPath() string {
	return "/_proxy/api/v1/validate"
}
