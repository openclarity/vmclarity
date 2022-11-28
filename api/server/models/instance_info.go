// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// InstanceInfo instance info
//
// swagger:model InstanceInfo
type InstanceInfo struct {

	// id
	ID string `json:"id,omitempty"`

	// instance ID
	InstanceID string `json:"instanceID,omitempty"`

	// instance name
	InstanceName string `json:"instanceName,omitempty"`

	// instance provider
	InstanceProvider CloudProvider `json:"instanceProvider,omitempty"`

	// instance region
	InstanceRegion string `json:"instanceRegion,omitempty"`
}

// Validate validates this instance info
func (m *InstanceInfo) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateInstanceProvider(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *InstanceInfo) validateInstanceProvider(formats strfmt.Registry) error {
	if swag.IsZero(m.InstanceProvider) { // not required
		return nil
	}

	if err := m.InstanceProvider.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("instanceProvider")
		} else if ce, ok := err.(*errors.CompositeError); ok {
			return ce.ValidateName("instanceProvider")
		}
		return err
	}

	return nil
}

// ContextValidate validate this instance info based on the context it is used
func (m *InstanceInfo) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateInstanceProvider(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *InstanceInfo) contextValidateInstanceProvider(ctx context.Context, formats strfmt.Registry) error {

	if err := m.InstanceProvider.ContextValidate(ctx, formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("instanceProvider")
		} else if ce, ok := err.(*errors.CompositeError); ok {
			return ce.ValidateName("instanceProvider")
		}
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *InstanceInfo) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *InstanceInfo) UnmarshalBinary(b []byte) error {
	var res InstanceInfo
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
