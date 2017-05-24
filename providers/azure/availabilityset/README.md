#Azure Availability Set

Create an Availability Sen Availability Set
## Argument reference

The following arguments are supported:

- **name** : (Required) Specifies the name of the Probe.

- **resource_group_name** : (Required) The name of the resource group in which to create the resource.

- **location** :  (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

- **platform_update_domain_count** : (Optional) Specifies the number of update domains that are used. Defaults to 5.

- **platform_fault_domain_count** :  (Optional) Specifies the number of fault domains that are used. Defaults to 3.

- **managed** : (Optional) Specifies whether the availability set is managed or not. Possible values are true (to specify aligned) or false (to specify classic). Default is false.

- **tags** : (Optional) A mapping of tags to assign to the resource.



## Attributes reference

The following attributes are exported:

- **id** : The virtual AvailabilitySet ID.


## Dependencies

This resource has required dependencies on:

- [Resource groups](../resourcegroup/) through resource_group_name field

And optinal ones through:

## Example

You'll find a json example [here](../../../internal/definitions/as_create.json)

## Running "real" tests

This library is provided with a suite of "real" tests to be ran against Azure. In order to run load balancer specific tests, you'll need to setup your test suite as [described here](../../../internal/)

And then run load balancers specific tests from the root of the project with:

```
$ gucumber --tags=@availability_set
