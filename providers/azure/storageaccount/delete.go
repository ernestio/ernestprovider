/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package storageaccount

import (
	"fmt"
)

// Delete : Deletes a network interface on azure
func (ev *Event) Delete() error {
	client := ev.client()

	_, err := client.Delete(ev.ResourceGroupName, ev.Name)
	if err != nil {
		err := fmt.Errorf("Error issuing AzureRM delete request for storage account %q: %s", ev.Name, err)
		ev.Log("error", err.Error())
		return err
	}

	return nil
}
