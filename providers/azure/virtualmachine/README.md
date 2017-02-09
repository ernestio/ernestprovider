#Azure Virtual Machine

Deploy windows or linux based Azure Virtual Machines.

## Argument reference

The following arguments are supported:

- **name** : (Required) Specifies the name of the virtual machine resource. Changing this forces a new resource to be created.
- **resource_group_name** : (Required) The name of the resource group in which to create the virtual machine.
- **location** : (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.
- **plan** : (Optional) A plan block.
  - **name** : (Required) Specifies the name of the image from the marketplace.
  - **publisher** : (Optional) Specifies the publisher of the image.
  - **product** : (Optional) Specifies the product of the image from the marketplace.
- **availability_set_id** : (Optional) The Id of the Availability Set in which to create the virtual machine
- **boot_diagnostics** : (Optional) A boot diagnostics profile block.
  - **enabled** : (Required) Whether to enable boot diagnostics for the virtual machine.
  - **storage_uri** : (Required) Blob endpoint for the storage account to hold the virtual machine's diagnostic files. This must be the root of a storage account, and not a storage container.
- **vm_size** : (Required) Specifies the size of the virtual machine.
- **storage_image_reference** : (Optional) A Storage Image Reference block.
  - **publisher** : (Required) Specifies the publisher of the image used to create the virtual machine. Changing this forces a new resource to be created.
  - **offer** : (Required) Specifies the offer of the image used to create the virtual machine. Changing this forces a new resource to be created.
  - **sku** : (Required) Specifies the SKU of the image used to create the virtual machine. Changing this forces a new resource to be created.
  - **version** : (Optional) Specifies the version of the image used to create the virtual machine. Changing this forces a new resource to be created.
- **storage_os_disk** : (Required) A Storage OS Disk block.
  - **name** : (Required) Specifies the disk name.
  - **vhd_uri** : (Required) Specifies the vhd uri. Changing this forces a new resource to be created.
  - **create_option** : (Required) Specifies how the virtual machine should be created. Possible values are attach and FromImage.
  - **caching** : (Optional) Specifies the caching requirements.
  - **image_uri** : (Optional) Specifies the image_uri in the form publisherName:offer:skus:version. image_uri can also specify the VHD uri of a custom VM image to clone. When cloning a custom disk image the os_type documented below becomes required.
  - **os_type** : (Optional) Specifies the operating system Type, valid values are windows, linux.
  - **disk_size_gb** : (Optional) Specifies the size of the data disk in gigabytes.
- **delete_os_disk_on_termination** : (Optional) Flag to enable deletion of the OS Disk VHD blob when the VM is deleted, defaults to false
- **storage_data_disk** : (Optional) A list of Storage Data disk blocks.
  - **name** : (Required) Specifies the name of the data disk.
  - **vhd_uri** : (Required) Specifies the uri of the location in storage where the vhd for the virtual machine should be placed.
  - **create_option** : (Required) Specifies how the data disk should be created.
  - **disk_size_gb** : (Required) Specifies the size of the data disk in gigabytes.
  - **caching** : (Optional) Specifies the caching requirements.
  - **lun** : (Required) Specifies the logical unit number of the data disk.
- **delete_data_disks_on_termination** : (Optional) Flag to enable deletion of Storage Disk VHD blobs when the VM is deleted, defaults to false
- **os_profile** : (Required) An OS Profile block.
  - **computer_name** : (Required) Specifies the name of the virtual machine.
  - **admin_username** : (Required) Specifies the name of the administrator account.
  - **admin_password** : (Required) Specifies the password of the administrator account.
  - **custom_data** : (Optional) Specifies a base-64 encoded string of custom data. The base-64 encoded string is decoded to a binary array that is saved as a file on the Virtual Machine. The maximum length of the binary array is 65535 bytes.
- **license_type** : (Optional, when a windows machine) Specifies the Windows OS license type. The only allowable value, if supplied, is Windows_Server.
- **os_profile_windows_config** : (Required, when a windows machine) A Windows config block.
  - **provision_vm_agent** : (Optional)
  - **enable_automatic_upgrades** : (Optional)
  - **winrm** : (Optional) A collection of WinRM configuration blocks.
    - **protocol** : (Required) Specifies the protocol of listener
    - **certificate_url** : (Optional) Specifies URL of the certificate with which new Virtual Machines is provisioned.
  - **additional_unattend_config** : (Optional) An Additional Unattended Config block.
    - **pass** : (Required) Specifies the name of the pass that the content applies to. The only allowable value is oobeSystem.
    - **component** : (Required) Specifies the name of the component to configure with the added content. The only allowable value is Microsoft-Windows-Shell-Setup.
    - **setting_name** : (Required) Specifies the name of the setting to which the content applies. Possible values are: FirstLogonCommands and AutoLogon.
    - **content** : (Optional) Specifies the base-64 encoded XML formatted content that is added to the unattend.xml file for the specified path and component.
- **os_profile_linux_config** : (Required, when a linux machine) A Linux config block.
    - **disable_password_authentication** : (Required) Specifies whether password authentication should be disabled.
    - **ssh_keys** : (Optional) Specifies a collection of path and key_data to be placed on the virtual machine.
- **os_profile_secrets** : (Optional) A collection of Secret blocks.
    - **source_vault_id** : (Required) Specifies the key vault to use.
    - **vault_certificates** : (Required, on windows machines) A collection of Vault Certificates.
      - **certificate_url** : (Required) It is the Base64 encoding of a JSON Object that which is encoded in UTF-8 of which the contents need to be data, dataType and password.
      - **certificate_store** : (Required, on windows machines) Specifies the certificate store on the Virtual Machine where the certificate should be added to.
- **network_interface_ids** : (Required) Specifies the list of resource IDs for the network interfaces associated with the virtual machine.
- **tags** : (Optional) A mapping of tags to assign to the resource.


## Attributes reference

The following attributes are exported:

- **id** : The virtual machine ID.

## Dependencies

This resource has required dependencies on:

- [Resource groups](../resourcegroup/) through resource_group_name field

## Example

You'll find a json example [here](../../../internal/definitions/vm_create.json)

## Running "real" tests

This library is provided with a suite of "real" tests to be ran against Azure. In order to run virtual machine specific tests, you'll need to setup your test suite as [described here](../../../internal/)

And then run virtual machine specific tests from the root of the project with:

```
$ gucumber --tags=@virtual_machine
```
