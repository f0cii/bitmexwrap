// Code generated by go-swagger; DO NOT EDIT.

package order

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

// NewOrderAmendBulkParams creates a new OrderAmendBulkParams object
// with the default values initialized.
func NewOrderAmendBulkParams() *OrderAmendBulkParams {
	var ()
	return &OrderAmendBulkParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewOrderAmendBulkParamsWithTimeout creates a new OrderAmendBulkParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewOrderAmendBulkParamsWithTimeout(timeout time.Duration) *OrderAmendBulkParams {
	var ()
	return &OrderAmendBulkParams{

		timeout: timeout,
	}
}

// NewOrderAmendBulkParamsWithContext creates a new OrderAmendBulkParams object
// with the default values initialized, and the ability to set a context for a request
func NewOrderAmendBulkParamsWithContext(ctx context.Context) *OrderAmendBulkParams {
	var ()
	return &OrderAmendBulkParams{

		Context: ctx,
	}
}

// NewOrderAmendBulkParamsWithHTTPClient creates a new OrderAmendBulkParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewOrderAmendBulkParamsWithHTTPClient(client *http.Client) *OrderAmendBulkParams {
	var ()
	return &OrderAmendBulkParams{
		HTTPClient: client,
	}
}

/*OrderAmendBulkParams contains all the parameters to send to the API endpoint
for the order amend bulk operation typically these are written to a http.Request
*/
type OrderAmendBulkParams struct {

	/*Orders
	  An array of orders.

	*/
	Orders *string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the order amend bulk params
func (o *OrderAmendBulkParams) WithTimeout(timeout time.Duration) *OrderAmendBulkParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the order amend bulk params
func (o *OrderAmendBulkParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the order amend bulk params
func (o *OrderAmendBulkParams) WithContext(ctx context.Context) *OrderAmendBulkParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the order amend bulk params
func (o *OrderAmendBulkParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the order amend bulk params
func (o *OrderAmendBulkParams) WithHTTPClient(client *http.Client) *OrderAmendBulkParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the order amend bulk params
func (o *OrderAmendBulkParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithOrders adds the orders to the order amend bulk params
func (o *OrderAmendBulkParams) WithOrders(orders *string) *OrderAmendBulkParams {
	o.SetOrders(orders)
	return o
}

// SetOrders adds the orders to the order amend bulk params
func (o *OrderAmendBulkParams) SetOrders(orders *string) {
	o.Orders = orders
}

// WriteToRequest writes these params to a swagger request
func (o *OrderAmendBulkParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if o.Orders != nil {

		// form param orders
		var frOrders string
		if o.Orders != nil {
			frOrders = *o.Orders
		}
		fOrders := frOrders
		if fOrders != "" {
			if err := r.SetFormParam("orders", fOrders); err != nil {
				return err
			}
		}

	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
