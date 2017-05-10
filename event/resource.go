/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package event

import (
	"github.com/r3labs/terraform/helper/schema"
)

// Resource : ...
type Resource interface {
	ValidateID(id string) bool
	SetID(id string)
	GetID() string
	SetState(state string)
	ResourceDataToEvent(d *schema.ResourceData) error
	EventToResourceData(d *schema.ResourceData) error
	SetComponents([]Event)
	Clone() (Event, error)
	Error(err error)
}
