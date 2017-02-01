/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package networkinterface

// Delete : Deletes a network interface on azure
func (ev *Event) Delete() error {
	_, err := ev.client().Delete(ev.ResourceGroupName, ev.Name, make(chan struct{}))

	return err
}
