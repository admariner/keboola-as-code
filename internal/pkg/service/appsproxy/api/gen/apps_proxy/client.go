// Code generated by goa v3.20.1, DO NOT EDIT.
//
// apps-proxy client
//
// Command:
// $ goa gen github.com/keboola/keboola-as-code/api/appsproxy --output
// ./internal/pkg/service/appsproxy/api

package appsproxy

import (
	"context"

	goa "goa.design/goa/v3/pkg"
)

// Client is the "apps-proxy" service client.
type Client struct {
	APIRootIndexEndpoint    goa.Endpoint
	APIVersionIndexEndpoint goa.Endpoint
	HealthCheckEndpoint     goa.Endpoint
	ValidateEndpoint        goa.Endpoint
}

// NewClient initializes a "apps-proxy" service client given the endpoints.
func NewClient(aPIRootIndex, aPIVersionIndex, healthCheck, validate goa.Endpoint) *Client {
	return &Client{
		APIRootIndexEndpoint:    aPIRootIndex,
		APIVersionIndexEndpoint: aPIVersionIndex,
		HealthCheckEndpoint:     healthCheck,
		ValidateEndpoint:        validate,
	}
}

// APIRootIndex calls the "ApiRootIndex" endpoint of the "apps-proxy" service.
func (c *Client) APIRootIndex(ctx context.Context) (err error) {
	_, err = c.APIRootIndexEndpoint(ctx, nil)
	return
}

// APIVersionIndex calls the "ApiVersionIndex" endpoint of the "apps-proxy"
// service.
func (c *Client) APIVersionIndex(ctx context.Context) (res *ServiceDetail, err error) {
	var ires any
	ires, err = c.APIVersionIndexEndpoint(ctx, nil)
	if err != nil {
		return
	}
	return ires.(*ServiceDetail), nil
}

// HealthCheck calls the "HealthCheck" endpoint of the "apps-proxy" service.
func (c *Client) HealthCheck(ctx context.Context) (res string, err error) {
	var ires any
	ires, err = c.HealthCheckEndpoint(ctx, nil)
	if err != nil {
		return
	}
	return ires.(string), nil
}

// Validate calls the "Validate" endpoint of the "apps-proxy" service.
func (c *Client) Validate(ctx context.Context, p *ValidatePayload) (res *Validations, err error) {
	var ires any
	ires, err = c.ValidateEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.(*Validations), nil
}
