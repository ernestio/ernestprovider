@real @sql_database

Feature: Sql database

  Scenario: Creting a sql_database
    Given I have no messages on the buffer
    And I request "azure_resource_group.create.azure" with "rg_create.json"
    And I request "azure_sql_server.create.azure" with "sql_create.json"
    When I request "azure_sql_database.create.azure" with "db_create.json"
    Then I should get a "azure_sql_database.create.azure.done" response with "name" as "mydb"
    And I should get a "azure_sql_database.create.azure.done" response with "resource_group_name" as "rg_test"
    And I should get a "azure_sql_database.create.azure.done" response with "server_name" as "mysqlserver098765"
    And I should get a "azure_sql_database.create.azure.done" response with "tags.t1" as "one"
    And I should get a "azure_sql_database.create.azure.done" response with "tags.t2" as "two"
    When I request "azure_sql_database.update.azure" with "db_update.json"
    Then I should get a "azure_sql_database.update.azure.done" response with "name" as "mydb"
    And I should get a "azure_sql_database.update.azure.done" response with "tags.t1" as "one"
    And I should get a "azure_sql_database.update.azure.done" response with "tags.t2" as "two"
    And I should get a "azure_sql_database.update.azure.done" response with "tags.t3" as "three"
    When I request "azure_sql_database.get.azure" with "db_update.json"
    Then I should get a "azure_sql_database.get.azure.done" response with "name" as "mydb"
    When I request "azure_sql_database.delete.azure" with "db_update.json"
    And I request "azure_sql_database.get.azure" with "db_update.json"
    Then I should get a "azure_sql_database.get.azure.error" response with "error" containing "Resource not found"
