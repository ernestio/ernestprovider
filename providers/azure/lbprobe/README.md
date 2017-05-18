#Azure Load balancer

Create a LoadBalancer Resource.


## Argument reference

The following arguments are supported:

- **name** : (Required) Specifies the name of the Probe.

- **resource_group_name** : (Required) The name of the resource group in which to create the resource.

- **loadbalancer_id** : (Required) The ID of the LoadBalancer in which to create the NAT Rule.

- **protocol** : (Optional) Specifies the protocol of the end point. Possible values are Http or Tcp. If Tcp is specified, a received ACK is required for the probe to be successful. If Http is specified, a 200 OK response from the specified URI is required for the probe to be successful.

- **port** : (Required) Port on which the Probe queries the backend endpoint. Possible values range from 1 to 65535, inclusive.

- **request_path** : (Optional) The URI used for requesting health status from the backend endpoint. Required if protocol is set to Http. Otherwise, it is not allowed.

- **interval_in_seconds** : (Optional) The interval, in seconds between probes to the backend endpoint for health status. The default value is 15, the minimum value is 5.

- **number_of_probes** : (Optional) The number of failed probe attempts after which the backend endpoint is removed from rotation. The default value is 2. NumberOfProbes multiplied by intervalInSeconds value must be greater or equal to 10.Endpoints are returned to rotation when at least one probe is successful.


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
$ gucumber --tags=@lb_probe
