@real @storage_container

Feature: Storage container

  Scenario: Creting a storage_container
    Given I have no messages on the buffer
    And I request "resource_group.create.azure" with "rg_create.json"
    And I request "storage_account.create.azure" with "sa_create.json"
    When I request "storage_container.create.azure" with "sc_create.json"
    Then I should get a "storage_container.create.azure.done" response with "name" as "sctest1283298731982"
    And I should get a "storage_container.create.azure.done" response with "resource_group_name" as "rg_test"
    And I should get a "storage_container.create.azure.done" response with "storage_account_name" as "satest1082376408127"
    And I should get a "storage_container.create.azure.done" response with "container_access_type" as "private"
    When I request "storage_container.get.azure" with "sc_create.json"
    Then I should get a "storage_container.get.azure.done" response with "name" as "sctest1283298731982"
    And I should get a "storage_container.get.azure.done" response with "resource_group_name" as "rg_test"
    And I should get a "storage_container.get.azure.done" response with "storage_account_name" as "satest1082376408127"
    And I should get a "storage_container.get.azure.done" response with "container_access_type" as "private"
    When I request "storage_container.delete.azure" with "sc_create.json"
    And I request "storage_container.get.azure" with "sc_create.json"
    Then I should get a "storage_container.get.azure.error" response with "error" as "Resource not found"
