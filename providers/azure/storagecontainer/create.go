/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package storagecontainer

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/storage"
)

// Create : Creates a storage account on azure
func (ev *Event) Create() error {
	blobClient, accountExists, err := ev.client().GetBlobStorageClientForStorageAccount(ev.ResourceGroupName, ev.StorageAccountName)
	if err != nil {
		return err
	}
	if !accountExists {
		return fmt.Errorf("Storage Account %q Not Found", ev.StorageAccountName)
	}

	var accessType storage.ContainerAccessType
	if ev.StorageType == "private" {
		accessType = storage.ContainerAccessType("")
	} else {
		accessType = storage.ContainerAccessType(ev.StorageType)
	}

	ev.Log("info", "Creating container "+ev.Name+" in storage account "+ev.StorageAccountName)
	_, err = blobClient.CreateContainerIfNotExists(ev.Name, accessType)
	if err != nil {
		err := fmt.Errorf("Error creating container %q in storage account %q: %s", ev.Name, ev.StorageAccountName, err)
		ev.Log("error", err.Error())
		return err
	}

	ev.ID = ev.Name

	return ev.Get()
}
