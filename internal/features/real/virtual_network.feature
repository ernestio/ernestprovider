@real @virtual_network

Feature: Virtual network

  Scenario: Creting a virtual network
    Given I have no messages on the buffer
    #And I request "azure_virtualnetwork.delete.azure" with "vn_create.json"
    #And I request "azure_resource_group.create.azure" with "rg_create.json"
    #And I request "azure_subnet.create.azure" with "sub_create.json"
    When I request "azure_virtualnetwork.create.azure" with "vn_create.json"
    Then I should get a "azure_virtualnetwork.create.azure.done" response with "name" as "testing_resource"
    And I should get a "azure_virtualnetwork.create.azure.done" response with "address_space.0" as "10.1.2.0/24"
    And I should get a "azure_virtualnetwork.create.azure.done" response with "dns_server_names.0" as "8.8.8.8"
    And I should get a "azure_virtualnetwork.create.azure.done" response with "dns_server_names.1" as "4.4.4.4"
    And I should get a "azure_virtualnetwork.create.azure.done" response with "subnets.0.name" as "subnet1"
    And I should get a "azure_virtualnetwork.create.azure.done" response with "subnets.0.address_prefix" as "10.1.2.0/24"
    And I should get a "azure_virtualnetwork.create.azure.done" response with "location" as "westus"
    And I should get a "azure_virtualnetwork.create.azure.done" response with "resource_group_name" as "ernest"
    And I should get a "azure_virtualnetwork.create.azure.done" response with "tags.t1" as "one"
    And I should get a "azure_virtualnetwork.create.azure.done" response with "tags.t2" as "two"
