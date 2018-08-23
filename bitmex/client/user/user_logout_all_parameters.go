// Code generated by go-swagger; DO NOT EDIT.

package user

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"
	"time"

	"golang.org/x/net/context"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"

	strfmt "github.com/go-openapi/strfmt"
)

// NewUserLogoutAllParams creates a new UserLogoutAllParams object
// with the default values initialized.
func NewUserLogoutAllParams() *UserLogoutAllParams {

	return &UserLogoutAllParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewUserLogoutAllParamsWithTimeout creates a new UserLogoutAllParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewUserLogoutAllParamsWithTimeout(timeout time.Duration) *UserLogoutAllParams {

	return &UserLogoutAllParams{

		timeout: timeout,
	}
}

// NewUserLogoutAllParamsWithContext creates a new UserLogoutAllParams object
// with the default values initialized, and the ability to set a context for a request
func NewUserLogoutAllParamsWithContext(ctx context.Context) *UserLogoutAllParams {

	return &UserLogoutAllParams{

		Context: ctx,
	}
}

// NewUserLogoutAllParamsWithHTTPClient creates a new UserLogoutAllParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewUserLogoutAllParamsWithHTTPClient(client *http.Client) *UserLogoutAllParams {

	return &UserLogoutAllParams{
		HTTPClient: client,
	}
}

/*UserLogoutAllParams contains all the parameters to send to the API endpoint
for the user logout all operation typically these are written to a http.Request
*/
type UserLogoutAllParams struct {
	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the user logout all params
func (o *UserLogoutAllParams) WithTimeout(timeout time.Duration) *UserLogoutAllParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the user logout all params
func (o *UserLogoutAllParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the user logout all params
func (o *UserLogoutAllParams) WithContext(ctx context.Context) *UserLogoutAllParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the user logout all params
func (o *UserLogoutAllParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the user logout all params
func (o *UserLogoutAllParams) WithHTTPClient(client *http.Client) *UserLogoutAllParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the user logout all params
func (o *UserLogoutAllParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WriteToRequest writes these params to a swagger request
func (o *UserLogoutAllParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
