#Azure Virtual Network

An Azure virtual network (VNet) is a representation of your own network in the cloud. 

It is a logical isolation of the Azure cloud dedicated to your subscription. 

You can fully control the IP address blocks, DNS settings, security policies, and route tables within this network. 

You can also further segment your VNet into subnets and launch Azure IaaS virtual machines (VMs) and/or Cloud services (PaaS role instances). 

Additionally, you can connect the virtual network to your on-premises network using one of the connectivity options available in Azure. 

In essence, you can expand your network to Azure, with complete control on IP address blocks with the benefit of enterprise scale Azure provides.

## Argument reference


The following arguments are supported:

- **name** : (Required) The name of the virtual network. Changing this forces a new resource to be created.

- **resource_group_name** : (Required) The name of the resource group in which to create the virtual network.

- **address_space** : (Required) The address space that is used the virtual network. You can supply more than one address space. Changing this forces a new resource to be created.

- **location** : (Required) The location/region where the virtual network is created. Changing this forces a new resource to be created.

- **dns_servers** : (Optional) List of IP addresses of DNS servers

- **subnet** : (Optional) Can be specified multiple times to define multiple subnets. Each subnet block supports fields documented below.

- **tags** : (Optional) A mapping of tags to assign to the resource.

The subnet block supports:

- **name** : (Required) The name of the subnet.

- **address_prefix** : (Required) The address prefix to use for the subnet.

- **security_group** : (Optional) The Network Security Group to associate with the subnet. (Referenced by id, ie. azurerm_network_security_group.test.id)


## Attributes reference

The following attributes are exported:

- **id** :  The Vritual network ID.

## Dependencies

This resource has required dependencies on:

- [Resource groups](../resourcegroup/) through resource_group_name field

## Example

You'll find a json example [here](../../../internal/definitions/vn_create.json)

## Running "real" tests

This library is provided with a suite of "real" tests to be ran against Azure. In order to run virtual network specific tests, you'll need to setup your test suite as [described here](../../../internal/)

And then run virtual network specific tests from the root of the project with:

```
$ gucumber --tags=@virtual_network
```
