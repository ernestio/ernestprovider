@real @network_interface

Feature: Network Interface

  Scenario: Creting a network_interface
    Given I have no messages on the buffer
    And I request "resource_group.create.azure" with "rg_create.json"
    And I request "virtual_network.create.azure" with "vn_create.json"
    And I request "subnet.create.azure" with "sub_create.json"
    When I request "network_interface.create.azure" with "ni_create.json"
    Then I should get a "network_interface.create.azure.done" response with "name" as "ni_test"
    And I should get a "network_interface.create.azure.done" response with "resource_group_name" as "rg_test"
    And I should get a "network_interface.create.azure.done" response with "location" as "westus"
    And I should get a "network_interface.create.azure.done" response with "ip_configuration.0.name" as "ip_conf_test"
    And I should get a "network_interface.create.azure.done" response with "ip_configuration.0.subnet_id" containing "sub_test_II"
    And I should get a "network_interface.create.azure.done" response with "ip_configuration.0.private_ip_address_allocation" as "dynamic"
    And I should get a "network_interface.create.azure.done" response with "tags.t1" as "one"
    And I should get a "network_interface.create.azure.done" response with "tags.t2" as "two"
    When I request "network_interface.get.azure" with "ni_create.json"
    Then I should get a "network_interface.get.azure.done" response with "name" as "ni_test"
    And I should get a "network_interface.get.azure.done" response with "resource_group_name" as "rg_test"
    And I should get a "network_interface.get.azure.done" response with "location" as "westus"
    And I should get a "network_interface.get.azure.done" response with "ip_configuration.0.name" as "ip_conf_test"
    And I should get a "network_interface.get.azure.done" response with "ip_configuration.0.subnet_id" containing "sub_test_II"
    And I should get a "network_interface.get.azure.done" response with "ip_configuration.0.private_ip_address_allocation" as "dynamic"
    And I should get a "network_interface.get.azure.done" response with "ip_configuration.0.name" as "ip_conf_test"
    And I should get a "network_interface.get.azure.done" response with "tags.t1" as "one"
    And I should get a "network_interface.get.azure.done" response with "tags.t2" as "two"
    When I request "network_interface.delete.azure" with "ni_create.json"
    And I request "network_interface.get.azure" with "ni_create.json"
    Then I should get a "network_interface.get.azure.error" response with "error" as "Resource not found"
