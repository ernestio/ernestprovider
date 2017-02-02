@real @storage_container

Feature: Storage container

  Scenario: Creting a storage_container
    Given I have no messages on the buffer
    When I request "azure_storage_container.create.azure" with "sc_create.json"
    Then I should get a "azure_storage_container.create.azure.done" response with "name" as "sc_test"
    And I should get a "azure_storage_container.create.azure.done" response with "resource_group_name" as "rg_test"
    And I should get a "azure_storage_container.create.azure.done" response with "storage_account_name" as "sa_name"
    And I should get a "azure_storage_container.create.azure.done" response with "container_access_type" as "private"
    When I request "azure_storage_container.update.azure" with "sc_update.json"
    Then I should get a "azure_storage_container.update.azure.done" response with "name" as "sc_test"
    And I should get a "azure_storage_container.update.azure.done" response with "resource_group_name" as "rg_test"
    And I should get a "azure_storage_container.update.azure.done" response with "storage_account_name" as "sa_name"
    And I should get a "azure_storage_container.update.azure.done" response with "container_access_type" as "public"
    When I request "azure_storage_container.get.azure" with "sc_update.json"
    Then I should get a "azure_storage_container.get.azure.done" response with "name" as "sc_test"
    And I should get a "azure_storage_container.get.azure.done" response with "resource_group_name" as "rg_test"
    And I should get a "azure_storage_container.get.azure.done" response with "storage_account_name" as "sa_name"
    And I should get a "azure_storage_container.get.azure.done" response with "container_access_type" as "public"
    When I request "azure_storage_container.delete.azure" with "sc_update.json"
    And I request "azure_storage_container.get.azure" with "sc_update.json"
    Then I should get a "azure_storage_container.get.azure.error" response with "error" as "Error reading storage_container  - removing from state"
