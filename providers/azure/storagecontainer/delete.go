/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package storagecontainer

import (
	"fmt"
)

// Delete : Deletes a network interface on azure
func (ev *Event) Delete() error {
	blobClient, accountExists, err := ev.client().GetBlobStorageClientForStorageAccount(ev.ResourceGroupName, ev.StorageAccountName)
	if err != nil {
		return err
	}
	if !accountExists {
		ev.Log("info", "Storage Account "+ev.StorageAccountName+" doesn't exist so the container won't exist")
		return nil
	}

	ev.Log("info", "Deleting storage container "+ev.Name+" in account "+ev.StorageAccountName)
	if _, err := blobClient.DeleteContainerIfExists(ev.Name); err != nil {
		err := fmt.Errorf("Error deleting storage container %q from storage account %q: %s", ev.Name, ev.StorageAccountName, err)
		ev.Log("error", err.Error())
		return err
	}

	ev.ID = ""
	return nil

}
