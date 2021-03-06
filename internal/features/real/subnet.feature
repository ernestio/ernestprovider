@real @subnet

Feature: Subnet

  Scenario: Creting a subnet
    Given I have no messages on the buffer
    And I request "resource_group.create.azure" with "rg_create.json"
    And I request "virtual_network.create.azure" with "vn_create.json"
    When I request "subnet.create.azure" with "sub_create.json"
    Then I should get a "subnet.create.azure.done" response with "name" as "sub_test_II"
    And I should get a "subnet.create.azure.done" response with "resource_group_name" as "rg_test"
    And I should get a "subnet.create.azure.done" response with "virtual_network_name" as "vn_test"
    And I should get a "subnet.create.azure.done" response with "address_prefix" as "10.0.2.0/24"
    And I should get a "subnet.create.azure.done" response with "ip_configurations.0" as "a"
    And I should get a "subnet.create.azure.done" response with "ip_configurations.1" as "b"
    When I request "subnet.update.azure" with "sub_update.json"
    Then I should get a "subnet.update.azure.done" response with "name" as "sub_test_II"
    And I should get a "subnet.update.azure.done" response with "ip_configurations.0" as "a"
    And I should get a "subnet.update.azure.done" response with "ip_configurations.1" as "b"
    And I should get a "subnet.update.azure.done" response with "ip_configurations.2" as "c"
    When I request "subnet.get.azure" with "sub_update.json"
    Then I should get a "subnet.get.azure.done" response with "name" as "sub_test_II"
    When I request "subnet.delete.azure" with "sub_update.json"
    And I request "subnet.get.azure" with "sub_update.json"
    Then I should get a "subnet.get.azure.error" response with "error" as "Resource not found"
