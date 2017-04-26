@real @lb

Feature: Load Balancer

  Scenario: Creting a load balancer
    Given I have no messages on the buffer
    And I request "resource_group.create.azure" with "rg_create.json"
    And I request "public_ip.create.azure" with "ip_create.json"
    When I request "lb.create.azure" with "lb_create.json"
    Then I should get a "lb.create.azure.done" response with "name" as "mytestingLB"
    And I should get a "lb.create.azure.done" response with "location" as "westus"
    And I should get a "lb.create.azure.done" response with "resource_group_name" as "rg_test"
    And I should get a "lb.create.azure.done" response with "frontend_ip_configuration.name" as "PublicIPAddress"
    And I should get a "lb.create.azure.done" response with "tags.t1" as "one"
    And I should get a "lb.create.azure.done" response with "tags.t2" as "two"
    When I request "lb.update.azure" with "lb_update.json"
    Then I should get a "lb.update.azure.done" response with "name" as "mytestingLB"
    And I should get a "lb.update.azure.done" response with "tags.t1" as "one"
    And I should get a "lb.update.azure.done" response with "tags.t2" as "two"
    And I should get a "lb.update.azure.done" response with "tags.t3" as "three"
    When I request "lb.get.azure" with "lb_update.json"
    Then I should get a "lb.get.azure.done" response with "name" as "mytestingLB"
    When I request "lb.delete.azure" with "lb_update.json"
    And I request "lb.get.azure" with "lb_update.json"
    Then I should get a "lb.get.azure.error" response with "error" as "Resource not found"
