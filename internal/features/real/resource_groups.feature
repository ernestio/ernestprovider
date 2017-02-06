@real @resource_group

Feature: Resource group

  Scenario: Creting a resource group
    Given I have no messages on the buffer
    When I request "azure_resource_group.create.azure" with "rg_create.json"
    Then I should get a "azure_resource_group.create.azure.done" response with "name" as "rg_test"
    And I should get a "azure_resource_group.create.azure.done" response with "location" as "westus"
    And I should get a "azure_resource_group.create.azure.done" response with "tags.t1" as "one"
    And I should get a "azure_resource_group.create.azure.done" response with "tags.t2" as "two"
    When I request "azure_resource_group.update.azure" with "rg_update.json"
    Then I should get a "azure_resource_group.update.azure.done" response with "name" as "rg_test"
    And I should get a "azure_resource_group.update.azure.done" response with "location" as "westus"
    And I should get a "azure_resource_group.update.azure.done" response with "tags.t1" as "one"
    And I should get a "azure_resource_group.update.azure.done" response with "tags.t2" as "two"
    And I should get a "azure_resource_group.update.azure.done" response with "tags.t3" as "three"
    When I request "azure_resource_group.get.azure" with "rg_update.json"
    Then I should get a "azure_resource_group.get.azure.done" response with "name" as "rg_test"
    And I should get a "azure_resource_group.get.azure.done" response with "location" as "westus"
    And I should get a "azure_resource_group.get.azure.done" response with "tags.t1" as "one"
    And I should get a "azure_resource_group.get.azure.done" response with "tags.t2" as "two"
    And I should get a "azure_resource_group.get.azure.done" response with "tags.t3" as "three"
    When I request "azure_resource_group.delete.azure" with "rg_update.json"
    And I request "azure_resource_group.get.azure" with "rg_update.json"
    Then I should get a "azure_resource_group.get.azure.error" response with "error" containing "Resource group 'rg_test' could not be found."
