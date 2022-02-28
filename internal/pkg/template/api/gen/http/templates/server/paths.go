// Code generated by goa v3.5.5, DO NOT EDIT.
//
// HTTP request path constructors for the templates service.
//
// Command:
// $ goa gen github.com/keboola/keboola-as-code/api/templates --output
// ./internal/pkg/template/api

package server

// IndexRootTemplatesPath returns the URL path to the templates service index-root HTTP endpoint.
func IndexRootTemplatesPath() string {
	return "/"
}

// IndexEndpointTemplatesPath returns the URL path to the templates service index HTTP endpoint.
func IndexEndpointTemplatesPath() string {
	return "/v1"
}

// HealthCheckTemplatesPath returns the URL path to the templates service health-check HTTP endpoint.
func HealthCheckTemplatesPath() string {
	return "/v1/health-check"
}
