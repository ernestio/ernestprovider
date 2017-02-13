@real @public_ip

Feature: public_ip

  Scenario: Creting a public_ip
    Given I have no messages on the buffer
    And I request "azure_resource_group.create.azure" with "rg_create.json"
    When I request "azure_public_ip.create.azure" with "ip_create.json"
    Then I should get a "azure_public_ip.create.azure.done" response with "name" as "mytestingIP"
    And I should get a "azure_public_ip.create.azure.done" response with "location" as "westus"
    And I should get a "azure_public_ip.create.azure.done" response with "resource_group_name" as "rg_test"
    And I should get a "azure_public_ip.create.azure.done" response with "public_ip_address_allocation" as "static"
    And I should get a "azure_public_ip.create.azure.done" response with "tags.t1" as "one"
    And I should get a "azure_public_ip.create.azure.done" response with "tags.t2" as "two"
    When I request "azure_public_ip.update.azure" with "ip_update.json"
    Then I should get a "azure_public_ip.update.azure.done" response with "name" as "mytestingIP"
    And I should get a "azure_public_ip.update.azure.done" response with "tags.t1" as "one"
    And I should get a "azure_public_ip.update.azure.done" response with "tags.t2" as "two"
    And I should get a "azure_public_ip.update.azure.done" response with "tags.t3" as "three"
    When I request "azure_public_ip.get.azure" with "ip_update.json"
    Then I should get a "azure_public_ip.get.azure.done" response with "name" as "mytestingIP"
    When I request "azure_public_ip.delete.azure" with "ip_update.json"
    And I request "azure_public_ip.get.azure" with "ip_update.json"
    Then I should get a "azure_public_ip.get.azure.error" response with "error" as "Resource not found"
