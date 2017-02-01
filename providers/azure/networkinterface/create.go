/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package networkinterface

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/arm/network"
)

// Create : Creates a network interface on azure
func (ev *Event) Create() error {
	client := ev.client()
	ev.Log("info", "Preparing arguments for Azure Network interface creation")

	properties := network.InterfacePropertiesFormat{
		EnableIPForwarding: &ev.EnableIPForwarding,
	}

	if ev.NetworkSecurityGroup != "" {
		properties.NetworkSecurityGroup = &network.SecurityGroup{
			ID: &ev.NetworkSecurityGroup,
		}
	}

	if len(ev.DNSServers) > 0 || ev.InternalDNSNameLabel != "" {
		ifaceDNSSettings := network.InterfaceDNSSettings{}
		if len(ev.DNSServers) > 0 {
			ifaceDNSSettings.DNSServers = &ev.DNSServers
		}

		if ev.InternalDNSNameLabel != "" {
			ifaceDNSSettings.InternalDNSNameLabel = &ev.InternalDNSNameLabel
		}

		properties.DNSSettings = &ifaceDNSSettings
	}

	ipConfigs, sgErr := ev.getExpandedAzureRmNetworkInterfaceIPConfigurations()
	if sgErr != nil {
		msg := "Error Building list of network interface IP Configurations: " + sgErr.Error()
		ev.Log("error", msg)
		return errors.New(msg)
	}
	if len(ipConfigs) > 0 {
		properties.IPConfigurations = &ipConfigs
	}

	iface := network.Interface{
		Name:                      &ev.Name,
		Location:                  &ev.Location,
		InterfacePropertiesFormat: &properties,
		Tags: &ev.Tags,
	}

	_, err := client.CreateOrUpdate(ev.ResourceGroupName, ev.Name, iface, make(chan struct{}))
	if err != nil {
		ev.Log("error", err.Error())
		return err
	}

	read, err := client.Get(ev.ResourceGroupName, ev.Name, "")
	if err != nil {
		ev.Log("error", err.Error())
		return err
	}
	if read.ID == nil {
		msg := fmt.Sprintf("Cannot read NIC %s (resource group %s) ID", ev.Name, ev.ResourceGroupName)
		ev.Log("error", msg)
		return errors.New(msg)
	}

	ev.ID = *read.ID

	return ev.Get()
}

func (ev *Event) getExpandedAzureRmNetworkInterfaceIPConfigurations() ([]network.InterfaceIPConfiguration, error) {
	ipConfigs := make([]network.InterfaceIPConfiguration, 0, len(ev.IPConfigurations))

	for _, ip := range ev.IPConfigurations {
		var allocationMethod network.IPAllocationMethod
		switch strings.ToLower(ip.PrivateIPAddressAllocation) {
		case "dynamic":
			allocationMethod = network.Dynamic
		case "static":
			allocationMethod = network.Static
		default:
			return []network.InterfaceIPConfiguration{}, fmt.Errorf(
				"valid values for private_ip_allocation_method are 'dynamic' and 'static' - got '%s'",
				ip.PrivateIPAddressAllocation)
		}

		properties := network.InterfaceIPConfigurationPropertiesFormat{
			Subnet: &network.Subnet{
				ID: &ip.Subnet,
			},
			PrivateIPAllocationMethod: allocationMethod,
		}

		if ip.PrivateIPAddress != "" {
			properties.PrivateIPAddress = &ip.PrivateIPAddress
		}

		if ip.PublicIPAddress != "" {
			properties.PublicIPAddress = &network.PublicIPAddress{
				ID: &ip.PublicIPAddress,
			}
		}

		if len(ip.LoadBalancerBackendAddressPools) > 0 {
			var ids []network.BackendAddressPool
			for _, p := range ip.LoadBalancerBackendAddressPools {
				id := network.BackendAddressPool{
					ID: &p,
				}

				ids = append(ids, id)
			}

			properties.LoadBalancerBackendAddressPools = &ids
		}

		if len(ip.LoadBalancerInboundNatRules) > 0 {
			var natRules []network.InboundNatRule
			for _, r := range ip.LoadBalancerInboundNatRules {
				rule := network.InboundNatRule{
					ID: &r,
				}

				natRules = append(natRules, rule)
			}

			properties.LoadBalancerInboundNatRules = &natRules
		}

		ipConfig := network.InterfaceIPConfiguration{
			Name: &ip.Name,
			InterfaceIPConfigurationPropertiesFormat: &properties,
		}

		ipConfigs = append(ipConfigs, ipConfig)
	}

	return ipConfigs, nil
}
