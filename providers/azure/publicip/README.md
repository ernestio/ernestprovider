#Azure Public IP

Create a Public IP Address.

## Argument reference


The following arguments are supported:

- **name** : (Required) Specifies the name of the Public IP resource . Changing this forces a new resource to be created.

- **resource_group_name** : (Required) The name of the resource group in which to create the public ip.

- **location** : (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

- **public_ip_address_allocation** : (Required) Defines whether the IP address is stable or dynamic. Options are Static or Dynamic.

- **idle_timeout_in_minutes** : (Optional) Specifies the timeout for the TCP idle connection. The value can be set between 4 and 30 minutes.

- **domain_name_label** : (Optional) Label for the Domain Name. Will be used to make up the FQDN. If a domain name label is specified, an A DNS record is created for the public IP in the Microsoft Azure DNS system.

- **reverse_fqdn** : (Optional) A fully qualified domain name that resolves to this public IP address. If the reverseFqdn is specified, then a PTR DNS record is created pointing from the IP address in the in-addr.arpa domain to the reverse FQDN.

- **tags** : (Optional) A mapping of tags to assign to the resource.

## Attributes reference

The following attributes are exported:

- **id** : The Public IP ID.

- **ip_address** : The IP address value that was allocated.

- **fqdn** : Fully qualified domain name of the A DNS record associated with the public IP. This is the concatenation of the domainNameLabel and the regionalized DNS zone

## Dependencies

This resource has required dependencies on:

- [Resource groups](../resourcegroup/) through resource_group_name field

## Example

You'll find a json example [here](../../../internal/definitions/ip_create.json)

## Running "real" tests

This library is provided with a suite of "real" tests to be ran against Azure. In order to run public ip specific tests, you'll need to setup your test suite as [described here](../../../internal/)

And then run public ip specific tests from the root of the project with:

```
$ gucumber --tags=@ip
