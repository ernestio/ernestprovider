/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package storageaccount

import (
	"errors"
	"fmt"
	"net/http"
)

// Get : Gets a storage account on azure
func (ev *Event) Get() error {
	client := ev.client()

	resp, err := client.GetProperties(ev.ResourceGroupName, ev.Name)
	if err != nil {
		return fmt.Errorf("Error reading the state of AzureRM Storage Account %q: %s", ev.Name, err)
	}
	if resp.StatusCode == http.StatusNotFound {
		ev.Log("error", "Not found "+ev.Name)
		return errors.New("Storage account " + ev.Name + " not found")
	}

	keys, err := client.ListKeys(ev.ResourceGroupName, ev.Name)
	if err != nil {
		return err
	}

	accessKeys := *keys.Keys
	ev.PrimaryAccessKey = *accessKeys[0].Value
	ev.SecondaryAccessKey = *accessKeys[1].Value
	ev.Location = *resp.Location
	ev.AccountType = string(resp.Sku.Name)
	ev.PrimaryLocation = *resp.AccountProperties.PrimaryLocation
	ev.SecondaryLocation = *resp.AccountProperties.SecondaryLocation

	if resp.AccountProperties.PrimaryEndpoints != nil {
		ev.PrimaryBlobEndpoint = *resp.AccountProperties.PrimaryEndpoints.Blob
		ev.PrimaryQueueEndpoint = *resp.AccountProperties.PrimaryEndpoints.Queue
		ev.PrimaryTableEndpoint = *resp.AccountProperties.PrimaryEndpoints.Table
		ev.PrimaryFileEndpoint = *resp.AccountProperties.PrimaryEndpoints.File
	}

	if resp.AccountProperties.SecondaryEndpoints != nil {
		if resp.AccountProperties.SecondaryEndpoints.Blob != nil {
			ev.SecondaryBlobEndpoint = *resp.AccountProperties.SecondaryEndpoints.Blob
		} else {
			ev.SecondaryBlobEndpoint = ""
		}
		if resp.AccountProperties.SecondaryEndpoints.Queue != nil {
			ev.SecondaryQueueEndpoint = *resp.AccountProperties.SecondaryEndpoints.Queue
		} else {
			ev.SecondaryQueueEndpoint = ""
		}
		if resp.AccountProperties.SecondaryEndpoints.Table != nil {
			ev.SecondaryTableEndpoint = *resp.AccountProperties.SecondaryEndpoints.Table
		} else {
			ev.SecondaryTableEndpoint = ""
		}
	}

	ev.Name = *resp.Name
	ev.Tags = *resp.Tags

	return nil
}
