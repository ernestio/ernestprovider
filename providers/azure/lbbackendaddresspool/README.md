#Azure Load balancer

Create a LoadBalancer Resource.


## Argument reference

The following arguments are supported:

- **name** : (Required) Specifies the name of the Backend Address Pool.

- **resource_group_name** : (Required) The name of the resource group in which to create the resource.

- **loadbalancer_id** : (Required) The ID of the LoadBalancer in which to create the Backend Address Pool.
»

## Attributes reference

The following attributes are exported:

- **id** : The ID of the LoadBalancer to which the resource is attached.


## Dependencies

This resource has required dependencies on:

- [Resource groups](../resourcegroup/) through resource_group_name field

And optinal ones through:

- [Loadbalancer ID](../lb/) through loadbalancer_id

## Example

You'll find a json example [here](../../../internal/definitions/lb_create.json)

## Running "real" tests

This library is provided with a suite of "real" tests to be ran against Azure. In order to run load balancer specific tests, you'll need to setup your test suite as [described here](../../../internal/)

And then run load balancers specific tests from the root of the project with:

```
$ gucumber --tags=@lb_backend_address_pool
