#Azure Load balancer Rule

Create a LoadBalancer Rule.


## Argument reference

The following arguments are supported:



## Attributes reference

The following attributes are exported:

- **id** : The ID of the LoadBalancer to which the resource is attached.


## Dependencies

This resource has required dependencies on:

- [Resource groups](../resourcegroup/) through resource_group_name field

And optinal ones through:

- [Loadbalancer id](../lb/) through loadbalancer_id field

- [Backend address pool](../lb_backend_address_pool/) throught backend_address_pool_id

## Example

You'll find a json example [here](../../../internal/definitions/lb_create.json)

## Running "real" tests

This library is provided with a suite of "real" tests to be ran against Azure. In order to run load balancer specific tests, you'll need to setup your test suite as [described here](../../../internal/)

And then run load balancers specific tests from the root of the project with:

```
$ gucumber --tags=@lb_rule
