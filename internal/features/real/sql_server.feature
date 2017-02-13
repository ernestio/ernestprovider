@real @sql_server

Feature: Sql Server

  Scenario: Creting a sql_server
    Given I have no messages on the buffer
    And I request "azure_resource_group.create.azure" with "rg_create.json"
    When I request "azure_sql_server.create.azure" with "sql_create.json"
    Then I should get a "azure_sql_server.create.azure.done" response with "name" as "mysqlserver"
    And I should get a "azure_sql_server.create.azure.done" response with "resource_group_name" as "rg_test"
    And I should get a "azure_sql_server.create.azure.done" response with "version" as "12.0"
    And I should get a "azure_sql_server.create.azure.done" response with "administrator_login" as "sql_admin"
    And I should get a "azure_sql_server.create.azure.done" response with "tags.0" as "a"
    And I should get a "azure_sql_server.create.azure.done" response with "tags.1" as "b"
    When I request "azure_sql_server.update.azure" with "sql_update.json"
    Then I should get a "azure_sql_server.update.azure.done" response with "name" as "mysqlserver"
    And I should get a "azure_sql_server.update.azure.done" response with "tags.0" as "a"
    And I should get a "azure_sql_server.update.azure.done" response with "tags.1" as "b"
    And I should get a "azure_sql_server.update.azure.done" response with "tags.2" as "c"
    When I request "azure_sql_server.get.azure" with "sql_update.json"
    Then I should get a "azure_sql_server.get.azure.done" response with "name" as "mysqlserver"
    When I request "azure_sql_server.delete.azure" with "sql_update.json"
    And I request "azure_sql_server.get.azure" with "sql_update.json"
    Then I should get a "azure_sql_server.get.azure.error" response with "error" as "Resource not found"
