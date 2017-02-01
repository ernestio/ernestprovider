/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package subnet

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/arm/network"
)

// Create : Creates a nat object on azure
func (ev *Event) Create() error {
	subnetClient := ev.client()

	ev.Log("info", "Preparing arguments for azure arm subnet creation")

	vnetName := ev.VirtualNetworkName
	resGroup := ev.ResourceGroupName
	addressPrefix := ev.AddressPrefix

	properties := network.SubnetPropertiesFormat{
		AddressPrefix: &addressPrefix,
	}

	if ev.NetworkSecurityGroup != "" {
		properties.NetworkSecurityGroup = &network.SecurityGroup{
			ID: &ev.NetworkSecurityGroup,
		}
	}

	if ev.RouteTable != "" {
		properties.RouteTable = &network.RouteTable{
			ID: &ev.RouteTable,
		}
	}

	subnet := network.Subnet{
		Name: &ev.Name,
		SubnetPropertiesFormat: &properties,
	}

	_, err := subnetClient.CreateOrUpdate(resGroup, vnetName, ev.Name, subnet, make(chan struct{}))
	if err != nil {
		return err
	}

	read, err := subnetClient.Get(resGroup, vnetName, ev.Name, "")
	if err != nil {
		return err
	}
	if read.ID == nil {
		msg := fmt.Sprintf("Cannot read Subnet %s/%s (resource group %s) ID", vnetName, ev.Name, resGroup)
		ev.Log("warning", msg)
	}

	ev.ID = *read.ID

	return ev.Get()
}
