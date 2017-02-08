#Azure Subnet

Subnet is a range of IP addresses in the VNet, you can divide a VNet into multiple subnets for organization and security. 

VMs and PaaS role instances deployed to subnets (same or different) within a VNet can communicate with each other without any extra configuration. 

You can also configure route tables and NSGs to a subnet.

## Argument reference


The following arguments are supported:

- **name : ** (Required) The name of the virtual network. Changing this forces a new resource to be created.

- **resource_group_name : ** (Required) The name of the resource group in which to create the subnet.

- **virtual_network_name : ** (Required) The name of the virtual network to which to attach the subnet.

- **address_prefix : ** (Required) The address prefix to use for the subnet.

- **network_security_group_id : ** (Optional) The ID of the Network Security Group to associate with the subnet.

- **route_table_id : ** (Optional) The ID of the Route Table to associate with the subnet.


## Attributes reference

The following attributes are exported:

- **id : ** The subnet ID.
- **ip_configurations : ** The collection of IP Configurations with IPs within this subnet.

## Dependencies

This resource has required dependencies on:

- [Resource groups](../resourcegroup/) through resource_group_name field
- [Virtual Networks](../virtualnetwork/) thorugh virtual_network_name field

And optional dependencies on:

- [Network security groups](../networksecuritygroups/) though network_security_group_id field
- [Route tables](../routetables/) through route_table_id field


