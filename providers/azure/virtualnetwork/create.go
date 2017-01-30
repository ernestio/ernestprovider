package virtualnetwork

import (
	"errors"

	"github.com/Azure/azure-sdk-for-go/arm/network"
)

// Create : Creates a nat object on azure
func (ev *Event) Create() error {
	c := ev.client()

	ev.Log("info", "preparing arguments for Azure ARM virtual network creation.")

	resGroup := ev.ResourceGroupName

	vnet := network.VirtualNetwork{
		Name:                           &ev.Name,
		Location:                       &ev.Location,
		VirtualNetworkPropertiesFormat: ev.getVirtualNetworkProperties(),
		Tags: &ev.Tags,
	}

	_, err := c.CreateOrUpdate(resGroup, ev.Name, vnet, make(chan struct{}))
	if err != nil {
		ev.Log("error", err.Error())
		return err
	}

	read, err := c.Get(resGroup, ev.Name, "")
	if err != nil {
		return err
	}
	if read.ID == nil {
		str := "Cannot read Virtual Network " + ev.Name + " (resource group " + resGroup + ") ID"
		ev.Log("error", str)
		return errors.New(str)
	}

	ev.ID = *read.ID

	return ev.Get()
}
