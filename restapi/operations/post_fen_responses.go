// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/xopoww/chess2pic/models"
)

// PostFenOKCode is the HTTP code returned for type PostFenOK
const PostFenOKCode int = 200

/*
PostFenOK API call result

swagger:response postFenOK
*/
type PostFenOK struct {

	/*
	  In: Body
	*/
	Payload *models.APIResult `json:"body,omitempty"`
}

// NewPostFenOK creates PostFenOK with default headers values
func NewPostFenOK() *PostFenOK {

	return &PostFenOK{}
}

// WithPayload adds the payload to the post fen o k response
func (o *PostFenOK) WithPayload(payload *models.APIResult) *PostFenOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the post fen o k response
func (o *PostFenOK) SetPayload(payload *models.APIResult) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PostFenOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}