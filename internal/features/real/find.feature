@find

Feature: Finding resources by resource group

  Scenario: Finding resources
    Given I have no messages on the buffer
    When I request "azure_public_ip.find.azure" with "ng_create.json"
    And I print the last info
