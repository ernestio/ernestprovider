@real @network_interface

Feature: Network Interface

  Scenario: Creting a network_interface
    Given I have no messages on the buffer
    When I request "azure_network_interface.create.azure" with "ni_create.json"
    Then I should get a "azure_network_interface.create.azure.done" response with "name" as "ni_test"
    And I should get a "azure_network_interface.create.azure.done" response with "resource_group_name" as "rg_test"
    And I should get a "azure_network_interface.create.azure.done" response with "location" as "westus"
    And I should get a "azure_network_interface.create.azure.done" response with "ip_configurations.0.name" as "ip_conf_test"
    And I should get a "azure_network_interface.create.azure.done" response with "ip_configurations.0.subnet" as "test_subnet"
    And I should get a "azure_network_interface.create.azure.done" response with "ip_configurations.0.private_ip_address_allocation" as "10.1.1.0"
    And I should get a "azure_network_interface.create.azure.done" response with "ip_configurations.0.name" as "ip_conf_test"
    And I should get a "azure_network_interface.create.azure.done" response with "tags.t1" as "one"
    And I should get a "azure_network_interface.create.azure.done" response with "tags.t2" as "two"
    When I request "azure_network_interface.update.azure" with "ni_update.json"
    Then I should get a "azure_network_interface.update.azure.done" response with "name" as "rg_test"
    And I should get a "azure_network_interface.create.azure.done" response with "tags.t2" as "three"
    When I request "azure_network_interface.get.azure" with "ni_update.json"
    Then I should get a "azure_network_interface.get.azure.done" response with "name" as "ni_test"
    And I should get a "azure_network_interface.get.azure.done" response with "resource_group_name" as "rg_test"
    And I should get a "azure_network_interface.get.azure.done" response with "location" as "westus"
    And I should get a "azure_network_interface.get.azure.done" response with "ip_configurations.0.name" as "ip_conf_test"
    And I should get a "azure_network_interface.get.azure.done" response with "ip_configurations.0.subnet" as "test_subnet"
    And I should get a "azure_network_interface.get.azure.done" response with "ip_configurations.0.private_ip_address_allocation" as "10.1.1.0"
    And I should get a "azure_network_interface.get.azure.done" response with "ip_configurations.0.name" as "ip_conf_test"
    And I should get a "azure_network_interface.get.azure.done" response with "tags.t1" as "one"
    And I should get a "azure_network_interface.get.azure.done" response with "tags.t2" as "three"
    When I request "azure_network_interface.delete.azure" with "ni_update.json"
    And I request "azure_network_interface.get.azure" with "ni_update.json"
    Then I should get a "azure_network_interface.get.azure.error" response with "error" as "Error reading network_interface  - removing from state"
