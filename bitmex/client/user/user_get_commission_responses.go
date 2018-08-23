// Code generated by go-swagger; DO NOT EDIT.

package user

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"

	models "github.com/SuperGod/coinex/bitmex/models"
)

// UserGetCommissionReader is a Reader for the UserGetCommission structure.
type UserGetCommissionReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *UserGetCommissionReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {

	case 200:
		result := NewUserGetCommissionOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewUserGetCommissionOK creates a UserGetCommissionOK with default headers values
func NewUserGetCommissionOK() *UserGetCommissionOK {
	return &UserGetCommissionOK{}
}

/*UserGetCommissionOK handles this case with default header values.

Request was successful
*/
type UserGetCommissionOK struct {
	Payload []*models.UserCommission
}

func (o *UserGetCommissionOK) Error() string {
	return fmt.Sprintf("[GET /user/commission][%d] userGetCommissionOK  %+v", 200, o.Payload)
}

func (o *UserGetCommissionOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
