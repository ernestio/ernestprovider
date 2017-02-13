@real @network_security_group

Feature: network_security_group

  Scenario: Creting a network_security_group
    Given I have no messages on the buffer
    And I request "azure_resource_group.create.azure" with "rg_create.json"
    When I request "azure_network_security_group.create.azure" with "sg_create.json"
    Then I should get a "azure_network_security_group.create.azure.done" response with "name" as "mytestingIP"
    And I should get a "azure_network_security_group.create.azure.done" response with "location" as "westus"
    And I should get a "azure_network_security_group.create.azure.done" response with "resource_group_name" as "rg_test"
    And I should get a "azure_network_security_group.create.azure.done" response with "tags.t1" as "one"
    And I should get a "azure_network_security_group.create.azure.done" response with "tags.t2" as "two"
    When I request "azure_network_security_group.update.azure" with "sg_update.json"
    Then I should get a "azure_network_security_group.update.azure.done" response with "name" as "mytestingIP"
    And I should get a "azure_network_security_group.update.azure.done" response with "tags.t1" as "one"
    And I should get a "azure_network_security_group.update.azure.done" response with "tags.t2" as "two"
    And I should get a "azure_network_security_group.update.azure.done" response with "tags.t3" as "three"
    When I request "azure_network_security_group.get.azure" with "sg_update.json"
    Then I should get a "azure_network_security_group.get.azure.done" response with "name" as "mytestingIP"
    When I request "azure_network_security_group.delete.azure" with "sg_update.json"
    And I request "azure_network_security_group.get.azure" with "sg_update.json"
    Then I should get a "azure_network_security_group.get.azure.error" response with "error" as "Resource not found"
