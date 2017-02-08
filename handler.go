/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package ernestprovider

import (
	"errors"
	"log"
	"strings"

	"github.com/ernestio/ernestprovider/event"
	"github.com/ernestio/ernestprovider/providers/azure/networkinterface"
	"github.com/ernestio/ernestprovider/providers/azure/resourcegroup"
	"github.com/ernestio/ernestprovider/providers/azure/storageaccount"
	"github.com/ernestio/ernestprovider/providers/azure/storagecontainer"
	"github.com/ernestio/ernestprovider/providers/azure/subnet"
	"github.com/ernestio/ernestprovider/providers/azure/virtualnetwork"
)

// Handle : Handles the given event
func Handle(ev *event.Event) (string, []byte) {
	var err error

	n := *ev
	if err := n.Process(); err != nil {
		n.Log("error", err.Error())
		return n.GetSubject() + ".error", n.GetBody()
	}

	parts := strings.Split(n.GetSubject(), ".")
	switch parts[1] {
	case "create":
		if err := n.Validate(); err != nil {
			n.Log("error", err.Error())
			return n.GetSubject() + ".error", n.GetBody()
		}
		err = n.Create()
	case "update":
		if err := n.Validate(); err != nil {
			n.Log("error", err.Error())
			return n.GetSubject() + ".error", n.GetBody()
		}
		err = n.Update()
	case "delete":
		if err := n.Validate(); err != nil {
			n.Log("error", err.Error())
			return n.GetSubject() + ".error", n.GetBody()
		}
		err = n.Delete()
	case "get":
		err = n.Get()
	case "find":
		err = n.Find()
	case "validate":
		if err := n.Validate(); err != nil {
			n.Log("error", err.Error())
			return n.GetSubject() + ".error", n.GetBody()
		}
	}

	if err != nil {
		n.Error(err)
		return n.GetSubject() + ".error", n.GetBody()
	}

	return n.GetSubject() + ".done", n.GetBody()
}

// GetAndHandle : Gets an event and Handles its results
func GetAndHandle(subject string, data []byte, key string) (string, []byte) {
	ev, err := GetEvent(subject, data, key)
	if err != nil {
		log.Println("[ERROR] : Event not found (A) - " + err.Error())
		return subject + ".error", data
	}
	if *ev == nil {
		log.Println("[ERROR] : Event not found (B) ")
		return subject + ".error", data
	}

	return Handle(ev)
}

// GetEvent : Gets a valid event based on a subject
func GetEvent(subject string, data []byte, key string) (*event.Event, error) {
	parts := strings.Split(subject, ".")
	switch parts[2] {
	case "aws":
		return getAWSEvent(subject, data, key)
	case "azure":
		return getAzureEvent(subject, data, key)
	}
	return nil, errors.New("Unkown provider")
}

func getAzureEvent(subject string, data []byte, key string) (*event.Event, error) {
	var ev event.Event
	var err error
	parts := strings.Split(subject, ".")
	val := event.NewValidator()
	switch parts[0] {
	case "azure_virtual_network":
		ev, err = virtualnetwork.New(subject, key, data, val)
	case "azure_resource_group":
		ev, err = resourcegroup.New(subject, key, data, val)
	case "azure_subnet":
		ev, err = subnet.New(subject, key, data, val)
	case "azure_network_interface":
		ev, err = networkinterface.New(subject, key, data, val)
	case "azure_storage_account":
		ev, err = storageaccount.New(subject, key, data, val)
	case "azure_storage_container":
		ev, err = storagecontainer.New(subject, key, data, val)
	}
	return &ev, err
}

func getAWSEvent(subject string, data []byte, key string) (*event.Event, error) {
	return nil, errors.New("Unconfigured provider")
}
