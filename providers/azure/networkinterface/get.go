/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package networkinterface

import (
	"errors"
	"fmt"
	"net/http"
)

// Get : Gets a network interface on azure
func (ev *Event) Get() error {
	resp, err := ev.client().Get(ev.ResourceGroupName, ev.Name, "")
	if err != nil {
		return fmt.Errorf("Error making Read request on Azure Network Interface %s: %s", ev.Name, err)
	}
	if resp.StatusCode == http.StatusNotFound {
		ev.ID = ""
		msg := "Network interface not found"
		ev.Log("warning", msg)
		return errors.New(msg)
	}

	iface := *resp.InterfacePropertiesFormat

	if iface.MacAddress != nil {
		if *iface.MacAddress != "" {
			ev.MacAddress = *iface.MacAddress
		}
	}

	if iface.IPConfigurations != nil && len(*iface.IPConfigurations) > 0 {
		var privateIPAddress *string
		///TODO: Change this to a loop when https://github.com/Azure/azure-sdk-for-go/issues/259 is fixed
		if (*iface.IPConfigurations)[0].InterfaceIPConfigurationPropertiesFormat != nil {
			privateIPAddress = (*iface.IPConfigurations)[0].InterfaceIPConfigurationPropertiesFormat.PrivateIPAddress
		}

		if *privateIPAddress != "" {
			ev.PrivateIPAddress = *privateIPAddress
		}
	}

	if iface.VirtualMachine != nil {
		if *iface.VirtualMachine.ID != "" {
			ev.VirtualMachineID = *iface.VirtualMachine.ID
		}
	}

	if iface.DNSSettings != nil {
		if iface.DNSSettings.AppliedDNSServers != nil && len(*iface.DNSSettings.AppliedDNSServers) > 0 {
			dnsServers := make([]string, 0, len(*iface.DNSSettings.AppliedDNSServers))
			for _, dns := range *iface.DNSSettings.AppliedDNSServers {
				dnsServers = append(dnsServers, dns)
			}

			ev.AppliedDNSServers = dnsServers
		}

		if iface.DNSSettings.InternalFqdn != nil && *iface.DNSSettings.InternalFqdn != "" {

			ev.InternalFQDN = *iface.DNSSettings.InternalFqdn
		}
	}

	ev.Tags = *resp.Tags

	return nil
}
