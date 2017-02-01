/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package storageaccount

import (
	"errors"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/arm/storage"
)

// Create : Creates a storage account on azure
func (ev *Event) Create() error {
	client := ev.client()

	sku := storage.Sku{
		Name: storage.SkuName(ev.AccountType),
	}

	opts := storage.AccountCreateParameters{
		Location: &ev.Location,
		Sku:      &sku,
		Tags:     &ev.Tags,
	}

	_, err := client.Create(ev.ResourceGroupName, ev.Name, opts, nil)
	if err != nil {
		ev.Log("error", err.Error())
		return err
	}

	// The only way to get the ID back apparently is to read the resource again
	read, err := client.GetProperties(ev.ResourceGroupName, ev.Name)

	// Set the ID right away if we have one
	if err == nil && read.ID != nil {
		msg := fmt.Sprintf("[INFO] storage account %q ID: %q", ev.Name, *read.ID)
		ev.Log("info", msg)
		ev.ID = *read.ID
	}

	// Check the read error now that we know it would exist without a create err
	if err != nil {
		ev.Log("error", err.Error())
		return err
	}

	// If we got no ID then the resource group doesn't yet exist
	if read.ID == nil {
		msg := fmt.Sprintf("Cannot read Storage Account %s (resource group %s) ID",
			ev.Name, ev.ResourceGroupName)
		ev.Log("error", msg)
		return errors.New(msg)
	}

	return ev.Get()
}
