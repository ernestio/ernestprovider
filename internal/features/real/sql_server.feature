@real @sql_server

Feature: Sql Server

  Scenario: Creting a sql_server
    Given I have no messages on the buffer
    And I request "resource_group.create.azure" with "rg_create.json"
    When I request "sql_server.create.azure" with "sql_create.json"
    Then I should get a "sql_server.create.azure.done" response with "name" as "mysqlserver098765"
    And I should get a "sql_server.create.azure.done" response with "resource_group_name" as "rg_test"
    And I should get a "sql_server.create.azure.done" response with "version" as "12.0"
    And I should get a "sql_server.create.azure.done" response with "administrator_login" as "sql_admin"
    And I should get a "sql_server.create.azure.done" response with "tags.t1" as "one"
    And I should get a "sql_server.create.azure.done" response with "tags.t2" as "two"
    When I request "sql_server.update.azure" with "sql_update.json"
    Then I should get a "sql_server.update.azure.done" response with "name" as "mysqlserver098765"
    And I should get a "sql_server.update.azure.done" response with "tags.t1" as "one"
    And I should get a "sql_server.update.azure.done" response with "tags.t2" as "two"
    And I should get a "sql_server.update.azure.done" response with "tags.t3" as "three"
    When I request "sql_server.get.azure" with "sql_update.json"
    Then I should get a "sql_server.get.azure.done" response with "name" as "mysqlserver098765"
    When I request "sql_server.delete.azure" with "sql_update.json"
    And I request "sql_server.get.azure" with "sql_update.json"
    Then I should get a "sql_server.get.azure.error" response with "error" containing "Resource not found"
