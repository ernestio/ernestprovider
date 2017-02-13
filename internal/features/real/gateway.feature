@real @local_network_gateway

Feature: Local network gateway

  Scenario: Creting a local_network_gateway
    Given I have no messages on the buffer
    And I request "azure_resource_group.create.azure" with "rg_create.json"
    When I request "azure_local_network_gateway.create.azure" with "ng_create.json"
    Then I should get a "azure_local_network_gateway.create.azure.done" response with "name" as "ng_test"
    And I should get a "azure_local_network_gateway.create.azure.done" response with "location" as "westus"
    And I should get a "azure_local_network_gateway.create.azure.done" response with "resource_group_name" as "rg_test"
    And I should get a "azure_local_network_gateway.create.azure.done" response with "gateway_address" as "12.13.14.15"
    And I should get a "azure_local_network_gateway.create.azure.done" response with "address_space.0" as "10.0.0.0/16"
    When I request "azure_local_network_gateway.update.azure" with "ng_update.json"
    Then I should get a "azure_local_network_gateway.update.azure.done" response with "name" as "ng_test"
    And I should get a "azure_local_network_gateway.update.azure.done" response with "gateway_address" as "12.13.14.15"
    And I should get a "azure_local_network_gateway.update.azure.done" response with "address_space.0" as "10.0.0.0/16"
    And I should get a "azure_local_network_gateway.update.azure.done" response with "address_space.1" as "10.0.0.0/24"
    When I request "azure_local_network_gateway.get.azure" with "ng_update.json"
    Then I should get a "azure_local_network_gateway.get.azure.done" response with "name" as "ng_test"
    When I request "azure_local_network_gateway.delete.azure" with "ng_update.json"
    And I request "azure_local_network_gateway.get.azure" with "ng_update.json"
    Then I should get a "azure_local_network_gateway.get.azure.error" response with "error" as "Resource not found"
