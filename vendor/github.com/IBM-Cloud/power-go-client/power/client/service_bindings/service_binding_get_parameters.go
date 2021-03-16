// Code generated by go-swagger; DO NOT EDIT.

package service_bindings

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"net/http"
	"time"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"

	strfmt "github.com/go-openapi/strfmt"
)

// NewServiceBindingGetParams creates a new ServiceBindingGetParams object
// with the default values initialized.
func NewServiceBindingGetParams() *ServiceBindingGetParams {
	var ()
	return &ServiceBindingGetParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewServiceBindingGetParamsWithTimeout creates a new ServiceBindingGetParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewServiceBindingGetParamsWithTimeout(timeout time.Duration) *ServiceBindingGetParams {
	var ()
	return &ServiceBindingGetParams{

		timeout: timeout,
	}
}

// NewServiceBindingGetParamsWithContext creates a new ServiceBindingGetParams object
// with the default values initialized, and the ability to set a context for a request
func NewServiceBindingGetParamsWithContext(ctx context.Context) *ServiceBindingGetParams {
	var ()
	return &ServiceBindingGetParams{

		Context: ctx,
	}
}

// NewServiceBindingGetParamsWithHTTPClient creates a new ServiceBindingGetParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewServiceBindingGetParamsWithHTTPClient(client *http.Client) *ServiceBindingGetParams {
	var ()
	return &ServiceBindingGetParams{
		HTTPClient: client,
	}
}

/*ServiceBindingGetParams contains all the parameters to send to the API endpoint
for the service binding get operation typically these are written to a http.Request
*/
type ServiceBindingGetParams struct {

	/*XBrokerAPIOriginatingIdentity
	  identity of the user that initiated the request from the Platform

	*/
	XBrokerAPIOriginatingIdentity *string
	/*XBrokerAPIVersion
	  version number of the Service Broker API that the Platform will use

	*/
	XBrokerAPIVersion string
	/*BindingID
	  binding id of binding to create

	*/
	BindingID string
	/*InstanceID
	  instance id of instance to provision

	*/
	InstanceID string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the service binding get params
func (o *ServiceBindingGetParams) WithTimeout(timeout time.Duration) *ServiceBindingGetParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the service binding get params
func (o *ServiceBindingGetParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the service binding get params
func (o *ServiceBindingGetParams) WithContext(ctx context.Context) *ServiceBindingGetParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the service binding get params
func (o *ServiceBindingGetParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the service binding get params
func (o *ServiceBindingGetParams) WithHTTPClient(client *http.Client) *ServiceBindingGetParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the service binding get params
func (o *ServiceBindingGetParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithXBrokerAPIOriginatingIdentity adds the xBrokerAPIOriginatingIdentity to the service binding get params
func (o *ServiceBindingGetParams) WithXBrokerAPIOriginatingIdentity(xBrokerAPIOriginatingIdentity *string) *ServiceBindingGetParams {
	o.SetXBrokerAPIOriginatingIdentity(xBrokerAPIOriginatingIdentity)
	return o
}

// SetXBrokerAPIOriginatingIdentity adds the xBrokerApiOriginatingIdentity to the service binding get params
func (o *ServiceBindingGetParams) SetXBrokerAPIOriginatingIdentity(xBrokerAPIOriginatingIdentity *string) {
	o.XBrokerAPIOriginatingIdentity = xBrokerAPIOriginatingIdentity
}

// WithXBrokerAPIVersion adds the xBrokerAPIVersion to the service binding get params
func (o *ServiceBindingGetParams) WithXBrokerAPIVersion(xBrokerAPIVersion string) *ServiceBindingGetParams {
	o.SetXBrokerAPIVersion(xBrokerAPIVersion)
	return o
}

// SetXBrokerAPIVersion adds the xBrokerApiVersion to the service binding get params
func (o *ServiceBindingGetParams) SetXBrokerAPIVersion(xBrokerAPIVersion string) {
	o.XBrokerAPIVersion = xBrokerAPIVersion
}

// WithBindingID adds the bindingID to the service binding get params
func (o *ServiceBindingGetParams) WithBindingID(bindingID string) *ServiceBindingGetParams {
	o.SetBindingID(bindingID)
	return o
}

// SetBindingID adds the bindingId to the service binding get params
func (o *ServiceBindingGetParams) SetBindingID(bindingID string) {
	o.BindingID = bindingID
}

// WithInstanceID adds the instanceID to the service binding get params
func (o *ServiceBindingGetParams) WithInstanceID(instanceID string) *ServiceBindingGetParams {
	o.SetInstanceID(instanceID)
	return o
}

// SetInstanceID adds the instanceId to the service binding get params
func (o *ServiceBindingGetParams) SetInstanceID(instanceID string) {
	o.InstanceID = instanceID
}

// WriteToRequest writes these params to a swagger request
func (o *ServiceBindingGetParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if o.XBrokerAPIOriginatingIdentity != nil {

		// header param X-Broker-API-Originating-Identity
		if err := r.SetHeaderParam("X-Broker-API-Originating-Identity", *o.XBrokerAPIOriginatingIdentity); err != nil {
			return err
		}

	}

	// header param X-Broker-API-Version
	if err := r.SetHeaderParam("X-Broker-API-Version", o.XBrokerAPIVersion); err != nil {
		return err
	}

	// path param binding_id
	if err := r.SetPathParam("binding_id", o.BindingID); err != nil {
		return err
	}

	// path param instance_id
	if err := r.SetPathParam("instance_id", o.InstanceID); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
