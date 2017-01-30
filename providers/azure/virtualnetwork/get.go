package virtualnetwork

import (
	"errors"
	"net/http"
)

// Get : Gets a nat object on azure
func (ev *Event) Get() error {
	resGroup := ev.ResourceGroupName
	name := ev.Name

	resp, err := ev.client().Get(resGroup, name, "")
	if err != nil {
		return errors.New("Error making Read request on Azure virtual network " + name + ": " + err.Error())
	}
	if resp.StatusCode == http.StatusNotFound {
		ev.ID = ""
		return nil
	}

	// update appropriate values
	ev.Name = *resp.Name
	ev.Location = *resp.Location
	ev.AddressSpace = *resp.AddressSpace.AddressPrefixes
	// ev.Subnets = *resp.Subnets

	dnses := []string{}
	for _, dns := range *resp.DhcpOptions.DNSServers {
		dnses = append(dnses, dns)
	}
	ev.DNSServerNames = dnses
	ev.Tags = *resp.Tags

	vnet := *resp.VirtualNetworkPropertiesFormat
	subnets := []subnet{}
	for _, sub := range *vnet.Subnets {
		s := subnet{}

		s.Name = *sub.Name
		s.AddressPrefix = *sub.SubnetPropertiesFormat.AddressPrefix
		if sub.SubnetPropertiesFormat.NetworkSecurityGroup != nil {
			s.SecurityGroup = *sub.SubnetPropertiesFormat.NetworkSecurityGroup.ID
		}

		subnets = append(subnets, s)
	}
	ev.Subnets = subnets

	return nil
}
