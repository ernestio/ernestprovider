#Azure Resource group

A resource group a container that holds related resources for an Azure solution. 

The resource group can include all the resources for the solution, or only those resources that you want to manage as a group. 

You decide how you want to allocate resources to resource groups based on what makes the most sense for your organization.

See azure doumentation for [resource groups](https://docs.microsoft.com/en-us/azure/azure-resource-manager/resource-group-overview#resource-groups)

## Argument reference

The following arguments are supported:

- **name** :  (Required) The name of the resource group. Must be unique on your Azure subscription.

- **location** :  (Required) The location where the resource group should be created. For a list of all Azure locations, please consult [this link](https://azure.microsoft.com/en-us/regions/).

- **tags** :  (Optional) A mapping of tags to assign to the resource.

## Attributes reference

The following attributes are exported:

**id** : The resource group ID.

## Dependencies

Rsource group does not have any dependency

## Example

You'll find a json example [here](../../../internal/definitions/rg_create.json)

## Running "real" tests

This library is provided with a suite of "real" tests to be ran against Azure. In order to run virtual network specific tests, you'll need to setup your test suite as [described here](../../../internal/)

And then run virtual network specific tests from the root of the project with:

```
$ gucumber --tags=@resource_group
