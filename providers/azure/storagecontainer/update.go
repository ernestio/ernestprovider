/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package storagecontainer

import (
	"errors"
)

// Update : Update a network interface on azure
func (ev *Event) Update() error {
	return errors.New(ev.Subject + " not supported")
}
