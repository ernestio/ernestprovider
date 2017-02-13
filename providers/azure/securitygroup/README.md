#Azure network security group

A network security group (NSG) contains a list of access control list (ACL) rules that allow or deny network traffic to your VM instances in a Virtual Network. NSGs can be associated with either subnets or individual VM instances within that subnet. When a NSG is associated with a subnet, the ACL rules apply to all the VM instances in that subnet. In addition, traffic to an individual VM can be restricted further by associating a NSG directly to that VM.

[Official Docs](https://docs.microsoft.com/en-us/azure/virtual-network/virtual-networks-nsg)

## Argument reference


The following arguments are supported:

- **name** : (Required) Specifies the name of the network security group. Changing this forces a new resource to be created.

- **resource_group_name** : (Required) The name of the resource group in which to create the availability set.

- **location** : (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

- **security_rule** : (Optional) Can be specified multiple times to define multiple security rules. Each security_rule block.
  
  - **name** : (Required) The name of the security rule.

  - **description** : (Optional) A description for this rule. Restricted to 140 characters.

  - **protocol** : (Required) Network protocol this rule applies to. Can be Tcp, Udp or * to match both.

  - **source_port_range** : (Required) Source Port or Range. Integer or range between 0 and 65535 or * to match any.

  - **destination_port_range** : (Required) Destination Port or Range. Integer or range between 0 and 65535 or * to match any.

  - **source_address_prefix** : (Required) CIDR or source IP range or * to match any IP. Tags such as ‘VirtualNetwork’, ‘AzureLoadBalancer’ and ‘Internet’ can also be used.

  - **destination_address_prefix** : (Required) CIDR or destination IP range or * to match any IP. Tags such as ‘VirtualNetwork’, ‘AzureLoadBalancer’ and ‘Internet’ can also be used.

  - **access** : (Required) Specifies whether network traffic is allowed or denied. Possible values are “Allow” and “Deny”.

  - **priority** : (Required) Specifies the priority of the rule. The value can be between 100 and 4096. The priority number must be unique for each rule in the collection. The lower the priority number, the higher the priority of the rule.

  - **direction** : (Required) The direction specifies if rule will be evaluated on incoming or outgoing traffic. Possible values are “Inbound” and “Outbound”.

- **tags** : (Optional) A mapping of tags to assign to the resource.


## Attributes reference

The following attributes are exported:

- **id** : The Network Security Group ID.

## Dependencies

This resource has required dependencies on:

- [Resource groups](../resourcegroup/) through resource_group_name field

## Example

You'll find a json example [here](../../../internal/definitions/sg_create.json)

## Running "real" tests

This library is provided with a suite of "real" tests to be ran against Azure. In order to run network security group specific tests, you'll need to setup your test suite as [described here](../../../internal/)

And then run network security group specific tests from the root of the project with:

```
$ gucumber --tags=@network_security_group
```
