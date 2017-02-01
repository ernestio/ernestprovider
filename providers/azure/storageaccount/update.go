/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package storageaccount

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/arm/storage"
)

// Update : Update a network interface on azure
func (ev *Event) Update() error {
	client := ev.client()

	storageAccountName := ev.Name
	resourceGroupName := ev.ResourceGroupName
	accountType := ev.AccountType

	sku := storage.Sku{
		Name: storage.SkuName(accountType),
	}

	opts := storage.AccountUpdateParameters{
		Sku:  &sku,
		Tags: &ev.Tags,
	}
	_, err := client.Update(resourceGroupName, storageAccountName, opts)

	if err != nil {
		err := fmt.Errorf("Error updating Azure Storage Account type %q: %s", storageAccountName, err)
		ev.Log("error", err.Error())
		return err
	}

	return nil
}
