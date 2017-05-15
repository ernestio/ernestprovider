#Azure SQL firewall rule


## Argument reference


The following arguments are supported:

- **name** : (Required) The name of the firewall rule.

- **resource_group_name** : (Required) The name of the resource group in which to create the firewall rule. This must be the same as firewall rule Server resource group currently.

- **location** : (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

- **server_name** : (Required) The name of the SQL Server on which to create the firewall rule.


## Dependencies

This resource has required dependencies on:

- [Resource groups](../resourcegroup/) through resource_group_name field
- [SQL Server name](../sqlserver/) through server_name field

## Example

You'll find a json example [here](../../../internal/definitions/db_create.json)

## Running "real" tests

This library is provided with a suite of "real" tests to be ran against Azure. In order to run virtual network specific tests, you'll need to setup your test suite as [described here](../../../internal/)

And then run sql firewall rule specific tests from the root of the project with:

```
$ gucumber --tags=@sql_firewall rule
