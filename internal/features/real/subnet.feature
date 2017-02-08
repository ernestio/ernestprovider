@real @subnet

Feature: Subnet

  Scenario: Creting a subnet
    Given I have no messages on the buffer
    And I request "azure_resource_group.create.azure" with "rg_create.json"
    And I request "azure_virtual_network.create.azure" with "vn_create.json"
    When I request "azure_subnet.create.azure" with "sub_create.json"
    Then I should get a "azure_subnet.create.azure.done" response with "name" as "test_subnet"
    And I should get a "azure_subnet.create.azure.done" response with "resource_group_name" as "rg_test"
    And I should get a "azure_subnet.create.azure.done" response with "virtual_networks_name" as "testing_resource"
    And I should get a "azure_subnet.create.azure.done" response with "address_prefix" as "10.1.2.0/24"
    And I should get a "azure_subnet.create.azure.done" response with "network_security_group_id" as "net_sec"
    And I should get a "azure_subnet.create.azure.done" response with "route_table_id" as "rt_test"
    And I should get a "azure_subnet.create.azure.done" response with "ip_configurations.0" as "a"
    And I should get a "azure_subnet.create.azure.done" response with "ip_configurations.1" as "b"
    #When I request "azure_subnet.update.azure" with "sub_update.json"
    #Then I should get a "azure_subnet.update.azure.done" response with "name" as "rg_test"
    #And I should get a "azure_subnet.update.azure.done" response with "ip_configurations.0" as "a"
    #And I should get a "azure_subnet.update.azure.done" response with "ip_configurations.1" as "b"
    #And I should get a "azure_subnet.update.azure.done" response with "ip_configurations.3" as "c"
    #When I request "azure_subnet.get.azure" with "sub_update.json"
    #Then I should get a "azure_subnet.get.azure.done" response with "name" as "rg_test"
    #When I request "azure_subnet.delete.azure" with "sub_update.json"
    #And I request "azure_subnet.get.azure" with "sub_update.json"
    #Then I should get a "azure_subnet.get.azure.error" response with "error" as "Error reading subnet  - removing from state"
