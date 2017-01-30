/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package template

import (
	"errors"

	"github.com/ernestio/ernestprovider"
)

// Event stores the template data
type Event struct {
	ernestprovider.Base
	ErrorMessage string `json:"error,omitempty"`
	Subject      string `json:"-"`
	Body         []byte `json:"-"`
	CryptoKey    string `json:"-"`
}

// New : Constructor
func New(subject string, body []byte, cryptoKey string) ernestprovider.Event {
	n := Event{Subject: subject, Body: body, CryptoKey: cryptoKey}

	return &n
}

// Validate checks if all criteria are met
func (ev *Event) Validate() error {
	return nil
}

// Find : Find an object on azure
func (ev *Event) Find() error {
	return errors.New(ev.Subject + " not supported")
}

// Create : Creates a nat object on azure
func (ev *Event) Create() error {
	return errors.New(ev.Subject + " not supported")
}

// Update : Updates a nat object on azure
func (ev *Event) Update() error {
	return errors.New(ev.Subject + " not supported")
}

// Delete : Deletes a nat object on azure
func (ev *Event) Delete() error {
	return errors.New(ev.Subject + " not supported")
}

// Get : Gets a nat object on azure
func (ev *Event) Get() error {
	return errors.New(ev.Subject + " not supported")
}
