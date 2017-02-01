/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package subnet

// Delete : Deletes a nat object on azure
func (ev *Event) Delete() error {
	// armMutexKV.Lock(vnetName)
	// defer armMutexKV.Unlock(vnetName)
	_, err := ev.client().Delete(ev.ResourceGroupName, ev.VirtualNetworkName, ev.Name, make(chan struct{}))

	return err
}
