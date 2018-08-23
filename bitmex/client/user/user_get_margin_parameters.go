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

// NewUserGetMarginParams creates a new UserGetMarginParams object
// with the default values initialized.
func NewUserGetMarginParams() *UserGetMarginParams {
	var (
		currencyDefault = string("XBt")
	)
	return &UserGetMarginParams{
		Currency: &currencyDefault,

		timeout: cr.DefaultTimeout,
	}
}

// NewUserGetMarginParamsWithTimeout creates a new UserGetMarginParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewUserGetMarginParamsWithTimeout(timeout time.Duration) *UserGetMarginParams {
	var (
		currencyDefault = string("XBt")
	)
	return &UserGetMarginParams{
		Currency: &currencyDefault,

		timeout: timeout,
	}
}

// NewUserGetMarginParamsWithContext creates a new UserGetMarginParams object
// with the default values initialized, and the ability to set a context for a request
func NewUserGetMarginParamsWithContext(ctx context.Context) *UserGetMarginParams {
	var (
		currencyDefault = string("XBt")
	)
	return &UserGetMarginParams{
		Currency: &currencyDefault,

		Context: ctx,
	}
}

// NewUserGetMarginParamsWithHTTPClient creates a new UserGetMarginParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewUserGetMarginParamsWithHTTPClient(client *http.Client) *UserGetMarginParams {
	var (
		currencyDefault = string("XBt")
	)
	return &UserGetMarginParams{
		Currency:   &currencyDefault,
		HTTPClient: client,
	}
}

/*UserGetMarginParams contains all the parameters to send to the API endpoint
for the user get margin operation typically these are written to a http.Request
*/
type UserGetMarginParams struct {

	/*Currency*/
	Currency *string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the user get margin params
func (o *UserGetMarginParams) WithTimeout(timeout time.Duration) *UserGetMarginParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the user get margin params
func (o *UserGetMarginParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the user get margin params
func (o *UserGetMarginParams) WithContext(ctx context.Context) *UserGetMarginParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the user get margin params
func (o *UserGetMarginParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the user get margin params
func (o *UserGetMarginParams) WithHTTPClient(client *http.Client) *UserGetMarginParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the user get margin params
func (o *UserGetMarginParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithCurrency adds the currency to the user get margin params
func (o *UserGetMarginParams) WithCurrency(currency *string) *UserGetMarginParams {
	o.SetCurrency(currency)
	return o
}

// SetCurrency adds the currency to the user get margin params
func (o *UserGetMarginParams) SetCurrency(currency *string) {
	o.Currency = currency
}

// WriteToRequest writes these params to a swagger request
func (o *UserGetMarginParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if o.Currency != nil {

		// query param currency
		var qrCurrency string
		if o.Currency != nil {
			qrCurrency = *o.Currency
		}
		qCurrency := qrCurrency
		if qCurrency != "" {
			if err := r.SetQueryParam("currency", qCurrency); err != nil {
				return err
			}
		}

	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
