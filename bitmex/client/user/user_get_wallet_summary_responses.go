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

// UserGetWalletSummaryReader is a Reader for the UserGetWalletSummary structure.
type UserGetWalletSummaryReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *UserGetWalletSummaryReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {

	case 200:
		result := NewUserGetWalletSummaryOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewUserGetWalletSummaryOK creates a UserGetWalletSummaryOK with default headers values
func NewUserGetWalletSummaryOK() *UserGetWalletSummaryOK {
	return &UserGetWalletSummaryOK{}
}

/*UserGetWalletSummaryOK handles this case with default header values.

Request was successful
*/
type UserGetWalletSummaryOK struct {
	Payload []*models.Transaction
}

func (o *UserGetWalletSummaryOK) Error() string {
	return fmt.Sprintf("[GET /user/walletSummary][%d] userGetWalletSummaryOK  %+v", 200, o.Payload)
}

func (o *UserGetWalletSummaryOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
