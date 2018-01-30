package sqldatabase

import (
	"encoding/json"

	"github.com/ernestio/ernestprovider/types"
)

type Event struct {
	types.Base
	ID                            string            `json:"id" diff:"-"`
	Name                          string            `json:"name" validate:"required" diff:"-"`
	Location                      string            `json:"location" diff:"-"`
	ResourceGroupName             string            `json:"resource_group_name" validate:"required" diff:"-"`
	ServerName                    string            `json:"server_name" validate:"required" diff:"server_name,immutable"`
	CreateMode                    string            `json:"create_mode" diff:"create_mode"`
	SourceDatabaseID              string            `json:"source_database_id" diff:"source_database_id"`
	RestorePointInTime            string            `json:"restore_point_in_time" diff:"restore_point_in_time"`
	Edition                       string            `json:"edition" diff:"edition"`
	Collation                     string            `json:"collation" diff:"collation"`
	MaxSizeBytes                  string            `json:"max_size_bytes" diff:"max_size_bytes"`
	RequestedServiceObjectiveID   string            `json:"requested_service_objective_id" diff:"requested_service_objective_id"`
	RequestedServiceObjectiveName string            `json:"requested_service_objective_name" diff:"requested_service_objective_name"`
	SourceDatabaseDeletionData    string            `json:"source_database_deletion_date" diff:"source_database_deletion_date"`
	ElasticPoolName               string            `json:"elastic_pool_name" diff:"elastic_pool_name,immutable"`
	Encryption                    string            `json:"encryption" diff:"encryption,immutable"`
	CreationDate                  string            `json:"creation_date" diff:"-"`
	DefaultSecondaryLocation      string            `json:"default_secondary_location" diff:"default_secondary_location,immutable"`
	Tags                          map[string]string `json:"tags" diff:"tags"`
	ClientID                      string            `json:"azure_client_id" diff:"-"`
	ClientSecret                  string            `json:"azure_client_secret" diff:"-"`
	TenantID                      string            `json:"azure_tenant_id" diff:"-"`
	SubscriptionID                string            `json:"azure_subscription_id" diff:"-"`
	Environment                   string            `json:"environment" diff:"-"`
	Components                    []json.RawMessage `json:"components" diff:"-"`
}
