#Azure Local network gateway

Creates a new local network gateway connection over which specific connections can be configured.


## Argument reference

The following arguments are supported:

- **name** : (Required) The name of the local network gateway. Changing this forces a new resource to be created.

- **resource_group_name** : (Required) The name of the resource group in which to create the local network gateway.

- **location** : (Required) The location/region where the local network gatway is created. Changing this forces a new resource to be created.

- **gateway_address** : (Required) The IP address of the gatway to which to connect.

- **address_space** : (Required) The list of string CIDRs representing the addredss spaces the gateway exposes.


## Attributes reference

The following attributes are exported:

- **id** : The local network gateway unique ID within Azure.


## Dependencies

This resource has required dependencies on:

- [Resource groups](../resourcegroup/) through resource_group_name field


## Example

You'll find a json example [here](../../../internal/definitions/ng_create.json)

## Running "real" tests

This library is provided with a suite of "real" tests to be ran against Azure. In order to run local network gateway specific tests, you'll need to setup your test suite as [described here](../../../internal/)

And then run local network gateway specific tests from the root of the project with:

```
$ gucumber --tags=@local_network_gateway
