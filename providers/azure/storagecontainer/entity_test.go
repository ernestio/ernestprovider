/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package storagecontainer

import (
	"strings"
	"testing"

	"github.com/ernestio/ernestprovider/event"
)

func validEvent() Event {
	properties := make(map[string]interface{})

	return Event{
		Name:               "supu",
		ResourceGroupName:  "rg_test",
		StorageAccountName: "sa_test",
		Properties:         properties,
		Validator:          event.NewValidator(),
	}
}

func TestRequiredName(t *testing.T) {
	ev := validEvent()
	ev.Name = ""
	err := ev.Validate()

	if err == nil {
		t.Error("No error has been received!")
	}

	if !strings.Contains(err.Error(), "Name is a required field") {
		t.Error("Output message does not contain name or required strings")
	}
}

func TestRequiredStorageAccountName(t *testing.T) {
	ev := validEvent()
	ev.StorageAccountName = ""
	err := ev.Validate()

	if err == nil {
		t.Error("No error has been received!")
	}

	if !strings.Contains(err.Error(), "StorageAccountName is a required field") {
		t.Error("Output message does not contain name or required strings")
	}
}

func TestRequiredResourceGroupName(t *testing.T) {
	ev := validEvent()
	ev.ResourceGroupName = ""
	err := ev.Validate()

	if err == nil {
		t.Error("No error has been received!")
	}

	if !strings.Contains(err.Error(), "ResourceGroupName is a required field") {
		t.Error("Output message does not contain name or required strings")
	}
}

func TestHappyPath(t *testing.T) {
	ev := validEvent()

	err := ev.Validate()
	if err != nil {
		println(err.Error())
		t.Error("I'm in a bad mood.")
	}
}
