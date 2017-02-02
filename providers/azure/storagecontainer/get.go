/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package storagecontainer

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/storage"
)

// Get : Gets a storage account on azure
func (ev *Event) Get() error {
	blobClient, accountExists, err := ev.client().GetBlobStorageClientForStorageAccount(ev.ResourceGroupName, ev.StorageAccountName)
	if err != nil {
		return err
	}
	if !accountExists {
		ev.Log("debug", "Storage account "+ev.StorageAccountName+" not found, removing container "+ev.ID+" from state")
		ev.ID = ""
		return nil
	}

	containers, err := blobClient.ListContainers(storage.ListContainersParameters{
		Prefix:  ev.Name,
		Timeout: 90,
	})
	if err != nil {
		err := fmt.Errorf("Failed to retrieve storage containers in account %q: %s", ev.Name, err)
		ev.Log("error", err.Error())
		return err
	}

	var found bool
	for _, cont := range containers.Containers {
		if cont.Name == ev.Name {
			found = true

			props := make(map[string]interface{})
			props["last_modified"] = cont.Properties.LastModified
			props["lease_status"] = cont.Properties.LeaseStatus
			props["lease_state"] = cont.Properties.LeaseState
			props["lease_duration"] = cont.Properties.LeaseDuration

			ev.Properties = props
		}
	}

	if !found {
		ev.Log("info", "Storage container "+ev.Name+" does not exist in account "+ev.StorageAccountName+", removing from state...")
		ev.ID = ""
	}
	return nil
}
