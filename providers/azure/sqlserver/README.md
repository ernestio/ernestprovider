#Azure SQL Server


## Argument reference


The following arguments are supported:

- **name** : (Required) The name of the SQL Server.

- **resource_group_name** : (Required) The name of the resource group in which to create the sql server.

- **location** : (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

- **version** : (Required) The version for the new server. Valid values are: 2.0 (for v11 server) and 12.0 (for v12 server).

- **administrator_login** : (Required) The administrator login name for the new server.

- **administrator_login_password**  (Required) The password for the new AdministratorLogin. Please following Azures Password Policy

- **tags** : (Optional) A mapping of tags to assign to the resource.


## Attributes reference

The following attributes are exported:

- **id** : The SQL Server ID.

- **fully_qualified_domain_name** : The fully qualified domain name of the Azure SQL Server (e.g. myServerName.database.windows.net)


## Dependencies

This resource has required dependencies on:

- [Resource groups](../resourcegroup/) through resource_group_name field

## Example

You'll find a json example [here](../../../internal/definitions/sql_create.json)

## Running "real" tests

This library is provided with a suite of "real" tests to be ran against Azure. In order to run virtual network specific tests, you'll need to setup your test suite as [described here](../../../internal/)

And then run sql server specific tests from the root of the project with:

```
$ gucumber --tags=@sql_server
