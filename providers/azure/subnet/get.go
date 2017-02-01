/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package subnet

import (
	"errors"
	"fmt"
	"net/http"
)

// Get : Gets a nat object on azure
func (ev *Event) Get() error {
	resp, err := ev.client().Get(ev.ResourceGroupName, ev.VirtualNetworkName, ev.Name, "")
	if err != nil {
		msg := fmt.Sprintf("Error making Read request on Azure Subnet %s: %s", ev.Name, err)
		ev.Log("error", msg)
		return errors.New(msg)
	}
	if resp.StatusCode == http.StatusNotFound {
		ev.ID = ""
		msg := fmt.Sprintf("Subnet %s not found", ev.Name)
		ev.Log("info", msg)
		return errors.New(msg)
	}

	if resp.SubnetPropertiesFormat.IPConfigurations != nil && len(*resp.SubnetPropertiesFormat.IPConfigurations) > 0 {
		ips := make([]string, 0, len(*resp.SubnetPropertiesFormat.IPConfigurations))
		for _, ip := range *resp.SubnetPropertiesFormat.IPConfigurations {
			ips = append(ips, *ip.ID)
		}

		ev.IPConfigurations = ips
	}

	return nil
}
