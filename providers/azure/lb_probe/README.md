#Azure Load balancer

Create a LoadBalancer Resource.


## Argument reference

The following arguments are supported:

- **name** : (Required) Specifies the name of the LoadBalancer.

- **resource_group_name** : (Required) The name of the resource group in which to create the LoadBalancer.

- **location** : (Required) Specifies the supported Azure location where the resource exists.

- **frontend_ip_configuration** : (Optional) A frontend ip configuration block.

  - **name** : (Required) Specifies the name of the frontend ip configuration.

  - **subnet_id** : (Optional) Reference to subnet associated with the IP Configuration.

  - **private_ip_address** : (Optional) Private IP Address to assign to the Load Balancer. The last one and first four IPs in any range are reserved and cannot be manually assigned.

  - **private_ip_address_allocation** : (Optional) Defines how a private IP address is assigned. Options are Static or Dynamic.

  - **public_ip_address_id** : (Optional) Reference to Public IP address to be associated with the Load Balancer.

- **tags** : (Optional) A mapping of tags to assign to the resource.


## Attributes reference

The following attributes are exported:

- **id** : The LoadBalancer ID.


## Dependencies

This resource has required dependencies on:

- [Resource groups](../resourcegroup/) through resource_group_name field

And optinal ones through:

- [Public IP](../publicip/) through frontend_ip_configuration::public_ip_address_id field

## Example

You'll find a json example [here](../../../internal/definitions/lb_create.json)

## Running "real" tests

This library is provided with a suite of "real" tests to be ran against Azure. In order to run load balancer specific tests, you'll need to setup your test suite as [described here](../../../internal/)

And then run load balancers specific tests from the root of the project with:

```
$ gucumber --tags=@lb
