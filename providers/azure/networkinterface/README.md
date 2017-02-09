#Azure Network Interface


## Argument reference

The following arguments are supported:

- **name** : (Required) The name of the network interface. Changing this forces a new resource to be created.

- **resource_group_name** : (Required) The name of the resource group in which to create the network interface.

- **location** : (Required) The location/region where the network interface is created. Changing this forces a new resource to be created.

- **network_security_group_id** : (Optional) The ID of the Network Security Group to associate with the network interface.

- **internal_dns_name_label** : (Optional) Relative DNS name for this NIC used for internal communications between VMs in the same VNet

- **enable_ip_forwarding** : (Optional) Enables IP Forwarding on the NIC. Defaults to false.

- **dns_servers** : (Optional) List of DNS servers IP addresses to use for this NIC, overrides the VNet-level server list

- **ip_configuration** : (Optional) Collection of ipConfigurations associated with this NIC.

  - **name** : (Required) User-defined name of the IP.

  - **subnet_id** : (Required) Reference to a subnet in which this NIC has been created.

  - **private_ip_address** : (Optional) Static IP Address.

  - **private_ip_address_allocation** : (Required) Defines how a private IP address is assigned. Options are Static or Dynamic.

  - **public_ip_address_id** : (Optional) Reference to a Public IP Address to associate with this NIC

  - **load_balancer_backend_address_pools_ids** : (Optional) List of Load Balancer Backend Address Pool IDs references to which this NIC belongs

load_balancer_inbound_nat_rules_ids - (Optional) List of Load Balancer Inbound Nat Rules IDs involving this NIC

- **tags** : (Optional) A mapping of tags to assign to the resource.


## Attributes reference

The following attributes are exported:

- **id** : The virtual NetworkConfiguration ID.

- **mac_address** : The media access control (MAC) address of the network interface.

- **private_ip_address** : The private ip address of the network interface.

- **virtual_machine_id** : Reference to a VM with which this NIC has been associated.

- **applied_dns_servers** : If the VM that uses this NIC is part of an Availability Set, then this list will have the union of all DNS servers from all NICs that are part of the Availability Set

- **internal_fqdn** : Fully qualified DNS name supporting internal communications between VMs in the same VNet

## Dependencies

Rsource group does not have any dependency

## Example

You'll find a json example [here](../../../internal/definitions/ni_create.json)

## Running "real" tests

This library is provided with a suite of "real" tests to be ran against Azure. In order to run network interface specific tests, you'll need to setup your test suite as [described here](../../../internal/)

And then run network interface specific tests from the root of the project with:

```
$ gucumber --tags=@resource_group
