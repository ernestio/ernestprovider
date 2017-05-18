#Azure Managed disks

Create a Managed disk Resource.


## Argument reference

The following arguments are supported:

- **name** : (Required) Specifies the name of the LoadBalancer.

- **resource_group_name** : (Required) The name of the resource group in which to create the LoadBalancer.

- **location** : (Required) Specifies the supported Azure location where the resource exists.

- **storage_account_type** : (Required) The type of storage to use for the managed disk. Allowable values are Standard_LRS or Premium_LRS.

- **create_option** : (Required) The method to use when creating the managed disk. [Import, Empty, Copy]

- **source_uri** : (Optional) URI to a valid VHD file to be used when create_option is Import.

- **source_resource_id** : (Optional) ID of an existing managed disk to copy when create_option is Copy.

- **os_type** : (Optional) Specify a value when the source of an Import or Copy operation targets a source that contains an operating system. Valid values are Linux or Windows

- **disk_size_gb** : (Required) Specifies the size of the managed disk to create in gigabytes. If create_option is Copy, then the value must be equal to or greater than the source's size.

- **tags** : (Optional) A mapping of tags to assign to the resource.


## Attributes reference

The following attributes are exported:

- **id** : The managed disk ID.


## Dependencies

This resource has required dependencies on:

- [Resource groups](../resourcegroup/) through resource_group_name field

And optinal ones through:

## Example

You'll find a json example [here](../../../internal/definitions/md_create.json)

## Running "real" tests

This library is provided with a suite of "real" tests to be ran against Azure. In order to run load balancer specific tests, you'll need to setup your test suite as [described here](../../../internal/)

And then run load balancers specific tests from the root of the project with:

```
$ gucumber --tags=@managed_disk
