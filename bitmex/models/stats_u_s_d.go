// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// StatsUSD stats u s d
// swagger:model StatsUSD
type StatsUSD struct {

	// currency
	Currency string `json:"currency,omitempty"`

	// root symbol
	// Required: true
	RootSymbol *string `json:"rootSymbol"`

	// turnover
	Turnover int64 `json:"turnover,omitempty"`

	// turnover24h
	Turnover24h int64 `json:"turnover24h,omitempty"`

	// turnover30d
	Turnover30d int64 `json:"turnover30d,omitempty"`

	// turnover365d
	Turnover365d int64 `json:"turnover365d,omitempty"`
}

// Validate validates this stats u s d
func (m *StatsUSD) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateRootSymbol(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *StatsUSD) validateRootSymbol(formats strfmt.Registry) error {

	if err := validate.Required("rootSymbol", "body", m.RootSymbol); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *StatsUSD) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *StatsUSD) UnmarshalBinary(b []byte) error {
	var res StatsUSD
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
