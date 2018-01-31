package virtualmachine

import (
	"encoding/json"

	"github.com/ernestio/ernestprovider/types"
)

type Event struct {
	types.Base
	ID                string `json:"id" diff:"-"`
	Name              string `json:"name" validate:"required" diff:"-"`
	ResourceGroupName string `json:"resource_group_name" validate:"required" diff:"-"`
	Location          string `json:"location" validate:"required" diff:"-"`
	Powered           bool   `json:"powered" diff:"powered"`
	Plan              struct {
		Name      string `json:"name" structs:"name" diff:"-"`
		Publisher string `json:"publisher" structs:"publisher" diff:"-"`
		Product   string `json:"product" structs:"product" diff:"-"`
	} `json:"plan" structs:"plan" diff:"-"`
	AvailabilitySet       string `json:"availability_set" diff:"-"`
	AvailabilitySetID     string `json:"availability_set_id" diff:"-"`
	LicenseType           string `json:"license_type" diff:"-"`
	VMSize                string `json:"vm_size" diff:"vm_size"`
	StorageImageReference struct {
		Publisher string `json:"publisher" structs:"publisher" diff:"-"`
		Offer     string `json:"offer" structs:"offer" diff:"-"`
		Sku       string `json:"sku" structs:"sku" diff:"-"`
		Version   string `json:"version" structs:"version" diff:"-"`
	} `json:"storage_image_reference" validate:"dive" diff:"-"`
	StorageOSDisk struct {
		Name               string `json:"name" structs:"name" diff:"-"`
		VhdURI             string `json:"vhd_uri" structs:"vhd_uri" diff:"-"`
		StorageAccount     string `json:"storage_account" structs:"-" diff:"-"`
		StorageContainer   string `json:"storage_container" structs:"-" diff:"-"`
		StorageAccountType string `json:"managed_disk_type" structs:"managed_disk_type" diff:"-"`
		ManagedDisk        string `json:"managed_disk" structs:"-" diff:"-"`
		ManagedDiskID      string `json:"managed_disk_id" structs:"managed_disk_id" diff:"-"`
		CreateOption       string `json:"create_option" structs:"create_option" diff:"-"`
		OSType             string `json:"os_type" structs:"os_type" diff:"-"`
		ImageURI           string `json:"image_uri" structs:"image_uri" diff:"-"`
		Caching            string `json:"caching" structs:"caching" diff:"-"`
	} `json:"storage_os_disk" validate:"dive" diff:"-"`
	DeleteOSDiskOnTermination bool `json:"delete_os_disk_on_termination" diff:"-"`
	StorageDataDisk           struct {
		Name               string `json:"name" structs:"name" diff:"-"`
		VhdURI             string `json:"vhd_uri" structs:"vhd_uri" diff:"-"`
		StorageAccount     string `json:"storage_account" structs:"-" diff:"-"`
		StorageAccountType string `json:"managed_disk_type" structs:"managed_disk_type" diff:"-"`
		StorageContainer   string `json:"storage_container" structs:"-" diff:"-"`
		ManagedDisk        string `json:"managed_disk" structs:"-" diff:"-"`
		ManagedDiskID      string `json:"managed_disk_id" structs:"managed_disk_id" diff:"-"`
		CreateOption       string `json:"create_option" structs:"create_option" diff:"-"`
		Size               *int32 `json:"disk_size_gb" structs:"disk_size_gb" diff:"size"`
		Lun                *int32 `json:"lun" structs:"lun" diff:"-"`
	} `json:"storage_data_disk" diff:"storage_data_disk"`
	DeleteDataDisksOnTermination bool             `json:"delete_data_disks_on_termination" diff:"-"`
	BootDiagnostics              []BootDiagnostic `json:"boot_diagnostics,omitempty" diff:"-"`
	OSProfile                    struct {
		ComputerName  string `json:"computer_name" structs:"computer_name" diff:"-"`
		AdminUsername string `json:"admin_username" structs:"admin_username" diff:"-"`
		AdminPassword string `json:"admin_password" structs:"admin_password" diff:"-"`
		CustomData    string `json:"custom_data" structs:"custom_data" diff:"-"`
	} `json:"os_profile" diff:"-"`
	OSProfileWindowsConfig *OSProfileWindowsConfig `json:"os_profile_windows_config,omitempty" diff:"-"`
	OSProfileLinuxConfig   struct {
		DisablePasswordAuthentication *bool    `json:"disable_password_authentication" structs:"disable_password_authentication" diff:"-"`
		SSHKeys                       []SSHKey `json:"ssh_keys" structs:"ssh_keys" diff:"-"`
	} `json:"os_profile_linux_config" structs:"os_profile_linux_config" diff:"-"`
	OSProfileSecrets    []Secret          `json:"os_profile_secrets" diff:"-"`
	NetworkInterfaces   []string          `json:"network_interfaces" diff:"network_interfaces"`
	NetworkInterfaceIDs []string          `json:"network_interface_ids" diff:"-"`
	Tags                map[string]string `json:"tags" diff:"tags"`
	ClientID            string            `json:"azure_client_id" diff:"-"`
	ClientSecret        string            `json:"azure_client_secret" diff:"-"`
	TenantID            string            `json:"azure_tenant_id" diff:"-"`
	SubscriptionID      string            `json:"azure_subscription_id" diff:"-"`
	Environment         string            `json:"environment" diff:"-"`
	Components          []json.RawMessage `json:"components" diff:"-"`
}

// OSProfileWindowsConfig ...
type OSProfileWindowsConfig struct {
	ProvisionVMAgent         bool               `json:"provision_vm_agent" structs:"provision_vm_agent" diff:"-"`
	EnableAutomaticUpgrades  bool               `json:"enable_automatic_upgrades" structs:"enable_automatic_upgrades" diff:"-"`
	WinRm                    []WinRM            `json:"winrm,omitempty" structs:"winrm,omitempty" diff:"-"`
	AdditionalUnattendConfig []UnattendedConfig `json:"additional_unattend_config,omitempty" structs:"additional_unattend_config,omitempty" diff:"-"`
}

type Secret struct {
	SourceVaultID           string             `json:"source_vault_id" structs:"source_vault_id" diff:"-"`
	SourceVaultCertificates []VaultCertificate `json:"vault_certificates" structs:"vault_certificates" diff:"-"`
}

type VaultCertificate struct {
	CertificateURL   string `json:"certificate_url" structs:"certificate_url" diff:"-"`
	CertificateStore string `json:"certificate_store" structs:"certificate_store" diff:"-"`
}

// WinRM ...
type WinRM struct {
	Protocol       string `json:"protocol" structs:"protocol" diff:"-"`
	CertificateURL string `json:"certificate_url" structs:"certification_url,omitempty" diff:"-"`
}

// SSHKey ...
type SSHKey struct {
	Path    string `json:"path" structs:"path" diff:"-"`
	KeyData string `json:"key_data" structs:"key_data" diff:"-"`
}

// BootDiagnostic ...
type BootDiagnostic struct {
	Enabled bool   `json:"enabled" structs:"enabled" diff:"-"`
	URI     string `json:"storage_uri" structs:"storage_uri" diff:"-"`
}

// UnattendedConfig ...
type UnattendedConfig struct {
	Pass        string `json:"pass" structs:"pass,omitempty" diff:"-"`
	Component   string `json:"component" structs:"component,omitempty" diff:"-"`
	SettingName string `json:"setting_name" structs:"setting_name,omitempty" diff:"-"`
	Content     string `json:"content" structs:"content,omitempty" diff:"-"`
}
