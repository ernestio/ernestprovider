#Azure Storage container

Create an Azure Storage Container.


## Argument reference


The following arguments are supported:

- **name** : (Required) The name of the storage container. Must be unique within the storage service the container is located.

- **resource_group_name** : (Required) The name of the resource group in which to create the storage container. Changing this forces a new resource to be created.

- **storage_account_name** : (Required) Specifies the storage account in which to create the storage container. Changing this forces a new resource to be created.

- **container_access_type** : (Required) The 'interface' for access the container provides. Can be either blob, container or private.


## Attributes reference

The following attributes are exported:

- **id** : The storage contariner Resource ID.

- **properties** : Key-value definition of additional properties associated to the storage container


## Dependencies

This resource has required dependencies on:

- [Resource groups](../resourcegroup/) through resource_group_name field
- [Storage account](../storageaccount/) through storage_account_name field

## Example

You'll find a json example [here](../../../internal/definitions/sc_create.json)

## Running "real" tests

This library is provided with a suite of "real" tests to be ran against Azure. In order to run virtual network specific tests, you'll need to setup your test suite as [described here](../../../internal/)

And then run virtual network specific tests from the root of the project with:

```
$ gucumber --tags=@storage_container
