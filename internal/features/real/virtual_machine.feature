@real @azure @virtual_machine

Feature: Virtual machine

  Scenario: Creting a resource group
    Given I have no messages on the buffer
    And I request "azure_resource_group_name.create.azure" with "rg_create.json"
    And I request "azure_virtual_network.create.azure" with "rg_create.json"
    And I request "azure_subnet.create.azure" with "sub_create.json"
    And I request "azure_network_interface.create.azure" with "ni_create.json"
    And I request "azure_storage_account.create.azure" with "sa_create.json"
    And I request "azure_storage_container.delete.azure" with "sc_create.json"
    And I request "azure_storage_container.create.azure" with "sc_create.json"
    When I request "azure_virtual_machine.create.azure" with "vm_create.json"
    Then I should get a "azure_virtual_machine.create.azure.done" response with "name" as "vm_test"
    And I should get a "azure_virtual_machine.create.azure.done" response with "location" as "westus"
    And I should get a "azure_virtual_machine.create.azure.done" response with "network_interface_ids.0" containing "ni_test"
    And I should get a "azure_virtual_machine.create.azure.done" response with "os_profile.admin_password" as "Password123"
    And I should get a "azure_virtual_machine.create.azure.done" response with "os_profile.admin_username" as "myadmin"
    And I should get a "azure_virtual_machine.create.azure.done" response with "os_profile.computer_name" as "hostname"
    And I should get a "azure_virtual_machine.create.azure.done" response with "resource_group_name" as "rg_test"
    And I should get a "azure_virtual_machine.create.azure.done" response with "tags.t1" as "one"
    And I should get a "azure_virtual_machine.create.azure.done" response with "tags.t2" as "two"
    When I request "azure_virtual_machine.update.azure" with "vm_update.json"
    Then I should get a "azure_virtual_machine.update.azure.done" response with "name" as "vm_test"
    And I should get a "azure_virtual_machine.update.azure.done" response with "location" as "westus"
    And I should get a "azure_virtual_machine.update.azure.done" response with "network_interface_ids.0" containing "ni_test"
    And I should get a "azure_virtual_machine.update.azure.done" response with "os_profile.admin_password" as "Password123"
    And I should get a "azure_virtual_machine.update.azure.done" response with "os_profile.admin_username" as "myadmin"
    And I should get a "azure_virtual_machine.update.azure.done" response with "os_profile.computer_name" as "hostname"
    And I should get a "azure_virtual_machine.update.azure.done" response with "resource_group_name" as "rg_test"
    And I should get a "azure_virtual_machine.update.azure.done" response with "tags.t1" as "one"
    And I should get a "azure_virtual_machine.update.azure.done" response with "tags.t2" as "two"
    And I should get a "azure_virtual_machine.update.azure.done" response with "tags.t3" as "three"
    When I request "azure_virtual_machine.get.azure" with "vm_update.json"
    Then I should get a "azure_virtual_machine.get.azure.done" response with "name" as "vm_test"
    And I should get a "azure_virtual_machine.get.azure.done" response with "location" as "westus"
    And I should get a "azure_virtual_machine.get.azure.done" response with "network_interface_ids.0" containing "ni_test"
    And I should get a "azure_virtual_machine.get.azure.done" response with "os_profile.admin_password" as "Password123"
    And I should get a "azure_virtual_machine.get.azure.done" response with "os_profile.admin_username" as "myadmin"
    And I should get a "azure_virtual_machine.get.azure.done" response with "os_profile.computer_name" as "hostname"
    And I should get a "azure_virtual_machine.get.azure.done" response with "resource_group_name" as "rg_test"
    And I should get a "azure_virtual_machine.get.azure.done" response with "tags.t1" as "one"
    And I should get a "azure_virtual_machine.get.azure.done" response with "tags.t2" as "two"
    And I should get a "azure_virtual_machine.get.azure.done" response with "tags.t3" as "three"
    When I request "azure_virtual_machine.delete.azure" with "vm_update.json"
    Then I should get a "azure_virtual_machine.delete.azure.done" response with "name" as "vm_test"
    When I request "azure_virtual_machine.get.azure" with "vm_update.json"
    Then I should get a "azure_virtual_machine.get.azure.error" response with "error" containing "Resource group 'rg_test' could not be found."