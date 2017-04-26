@real @azure @storage_account

Feature: Storage account

  Scenario: Creting a storage_account
    Given I have no messages on the buffer
    And I request "resource_group.create.azure" with "rg_create.json"
    When I request "storage_account.create.azure" with "sa_create.json"
    Then I should get a "storage_account.create.azure.done" response with "name" as "satest1082376408127"
    And I should get a "storage_account.create.azure.done" response with "resource_group_name" as "rg_test"
    And I should get a "storage_account.create.azure.done" response with "location" as "westus"
    And I should get a "storage_account.create.azure.done" response with "account_type" as "Standard_LRS"
    And I should get a "storage_account.create.azure.done" response with "tags.t1" as "one"
    And I should get a "storage_account.create.azure.done" response with "tags.t2" as "two"
    When I request "storage_account.update.azure" with "sa_update.json"
    Then I should get a "storage_account.update.azure.done" response with "name" as "satest1082376408127"
    And I should get a "storage_account.update.azure.done" response with "tags.t3" as "three"
    When I request "storage_account.get.azure" with "sa_update.json"
    Then I should get a "storage_account.get.azure.done" response with "name" as "satest1082376408127"
    And I should get a "storage_account.get.azure.done" response with "resource_group_name" as "rg_test"
    And I should get a "storage_account.get.azure.done" response with "location" as "westus"
    And I should get a "storage_account.get.azure.done" response with "account_type" as "Standard_LRS"
    And I should get a "storage_account.get.azure.done" response with "tags.t1" as "one"
    And I should get a "storage_account.get.azure.done" response with "tags.t2" as "two"
    And I should get a "storage_account.get.azure.done" response with "tags.t3" as "three"
    When I request "storage_account.delete.azure" with "sa_update.json"
    And I request "storage_account.get.azure" with "sa_update.json"
    Then I should get a "storage_account.get.azure.error" response with "error" as "Resource not found"
