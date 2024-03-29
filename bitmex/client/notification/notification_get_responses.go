// Code generated by go-swagger; DO NOT EDIT.

package notification

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"

	models "github.com/frankrap/bitmexwrap/bitmex/models"
)

// NotificationGetReader is a Reader for the NotificationGet structure.
type NotificationGetReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *NotificationGetReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {

	case 200:
		result := NewNotificationGetOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	case 400:
		result := NewNotificationGetBadRequest()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	case 401:
		result := NewNotificationGetUnauthorized()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	case 404:
		result := NewNotificationGetNotFound()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewNotificationGetOK creates a NotificationGetOK with default headers values
func NewNotificationGetOK() *NotificationGetOK {
	return &NotificationGetOK{}
}

/*NotificationGetOK handles this case with default header values.

Request was successful
*/
type NotificationGetOK struct {
	Payload []*models.Notification
}

func (o *NotificationGetOK) Error() string {
	return fmt.Sprintf("[GET /notification][%d] notificationGetOK  %+v", 200, o.Payload)
}

func (o *NotificationGetOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewNotificationGetBadRequest creates a NotificationGetBadRequest with default headers values
func NewNotificationGetBadRequest() *NotificationGetBadRequest {
	return &NotificationGetBadRequest{}
}

/*NotificationGetBadRequest handles this case with default header values.

Parameter Error
*/
type NotificationGetBadRequest struct {
	Payload *models.Error
}

func (o *NotificationGetBadRequest) Error() string {
	return fmt.Sprintf("[GET /notification][%d] notificationGetBadRequest  %+v", 400, *o.Payload)
}

func (o *NotificationGetBadRequest) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewNotificationGetUnauthorized creates a NotificationGetUnauthorized with default headers values
func NewNotificationGetUnauthorized() *NotificationGetUnauthorized {
	return &NotificationGetUnauthorized{}
}

/*NotificationGetUnauthorized handles this case with default header values.

Unauthorized
*/
type NotificationGetUnauthorized struct {
	Payload *models.Error
}

func (o *NotificationGetUnauthorized) Error() string {
	return fmt.Sprintf("[GET /notification][%d] notificationGetUnauthorized  %+v", 401, o.Payload)
}

func (o *NotificationGetUnauthorized) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewNotificationGetNotFound creates a NotificationGetNotFound with default headers values
func NewNotificationGetNotFound() *NotificationGetNotFound {
	return &NotificationGetNotFound{}
}

/*NotificationGetNotFound handles this case with default header values.

Not Found
*/
type NotificationGetNotFound struct {
	Payload *models.Error
}

func (o *NotificationGetNotFound) Error() string {
	return fmt.Sprintf("[GET /notification][%d] notificationGetNotFound  %+v", 404, o.Payload)
}

func (o *NotificationGetNotFound) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
