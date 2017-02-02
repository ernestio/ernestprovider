@real @storage_account

Feature: Storage account

  Scenario: Creting a storage_account
    Given I have no messages on the buffer
    When I request "azure_storage_account.create.azure" with "sa_create.json"
    Then I should get a "azure_storage_account.create.azure.done" response with "name" as "sa_test"
    And I should get a "azure_storage_account.create.azure.done" response with "resource_group_name" as "rg_test"
    And I should get a "azure_storage_account.create.azure.done" response with "location" as "westus"
    And I should get a "azure_storage_account.create.azure.done" response with "account_type" as "test"
    And I should get a "azure_storage_account.create.azure.done" response with "tags.t1" as "one"
    And I should get a "azure_storage_account.create.azure.done" response with "tags.t2" as "two"
    When I request "azure_storage_account.update.azure" with "sa_update.json"
    Then I should get a "azure_storage_account.update.azure.done" response with "name" as "rg_test"
    And I should get a "azure_storage_account.create.azure.done" response with "tags.t2" as "three"
    When I request "azure_storage_account.get.azure" with "sa_update.json"
    Then I should get a "azure_storage_account.get.azure.done" response with "name" as "sa_test"
    And I should get a "azure_storage_account.get.azure.done" response with "resource_group_name" as "rg_test"
    And I should get a "azure_storage_account.get.azure.done" response with "location" as "westus"
    And I should get a "azure_storage_account.create.azure.done" response with "account_type" as "test"
    And I should get a "azure_storage_account.get.azure.done" response with "tags.t1" as "one"
    And I should get a "azure_storage_account.get.azure.done" response with "tags.t2" as "three"
    When I request "azure_storage_account.delete.azure" with "sa_update.json"
    And I request "azure_storage_account.get.azure" with "sa_update.json"
    Then I should get a "azure_storage_account.get.azure.error" response with "error" as "Error reading storage_account  - removing from state"