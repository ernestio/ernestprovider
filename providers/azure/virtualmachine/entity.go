/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package virtualmachine

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"

	aes "github.com/ernestio/crypto/aes"
	"github.com/ernestio/ernestprovider/event"
	"github.com/ernestio/ernestprovider/providers/azure"
)

// Event : This is the Ernest representation of an azure networkinterface
type Event struct {
	event.Base
	ID                string `json:"id"`
	Name              string `json:"name" validate:"required"`
	ResourceGroupName string `json:"resource_group_name" validate:"required"`
	Location          string `json:"location" validate:"required"`
	Plan              struct {
		Name      string `json:"name" validate:"required"`
		Publisher string `json:"publisher" validate:"required"`
		Product   string `json:"product" validate:"required"`
	} `json:"plan" validate:"dive"`
	AvailabilitySetID     string `json:"availability_set_id"`
	LicenseType           string `json:"license_type"`
	VMSize                string `json:"vm_size"`
	StorageImageReference struct {
		Publisher string `json:"publisher" validate:"required"`
		Offer     string `json:"offer" validate:"offer"`
		Sku       string `json:"sku" validate:"required"`
		Version   string `json:"version"`
	} `json:"storage_image_reference" validate:"dive"`
	StorageOSDisk struct {
		Name         string `json:"name" validate:"required"`
		VhdURI       string `json:"vhd_uri" validate:"required"`
		CreateOption string `json:"create_option" validate:"required"`
		OSType       string `json:"os_type"`
		ImageURI     string `json:"image_uri"`
		Caching      string `json:"caching"`
	} `json:"storage_os_disk" validate:"dive"`
	DeleteOSDiskOnTermination bool `json:"delete_os_disk_on_termination"`
	StorageDataDisk           struct {
		Name         string `json:"name" validate:"required"`
		VhdURI       string `json:"vhd_uri" validate:"required"`
		CreateOption string `json:"create_option" validate:"required"`
		Size         int32  `json:"disk_size_gb" validate:"required"`
		Lun          int32  `json:"lun"`
	} `json:"storage_data_disk"`
	DeleteDataDisksOnTermination bool `json:"delete_data_disks_on_termination"`
	BootDiagnostics              struct {
		Enabled bool   `json:"enabled"`
		URI     string `json:"storage_uri"`
	} `json:"boot_diagnostics"`
	OSProfile struct {
		ComputerName  string `json:"computer_name" validation:"required"`
		AdminUsername string `json:"admin_username" validation:"required"`
		AdminPassword string `json:"admin_password" validation:"required"`
		CustomData    string `json:"custom_data"`
	} `json:"os_profile"`
	OSProfileWindowsConfig struct {
		ProvisionVMAgent        bool `json:"provision_vm_agent"`
		EnableAutomaticUpgrades bool `json:"enable_automatic_upgrades"`
		WinRm                   struct {
			Protocol       string `json:"protocol" validate:"required"`
			CertificateURL string `json:"certificate_url" validate:"required"`
		} `json:"winrm"`
		AdditionalUnattendConfig struct {
			Pass        string `json:"pass" validate:"required"`
			Component   string `json:"component" validate:"required"`
			SettingName string `json:"setting_name" validate:"required"`
			Content     string `json:"content" validate:"required"`
		} `json:"additional_unattend_config"`
	} `json:"os_profile_windows_config"`
	OSProfileLinuxConfig struct {
		DisablePasswordAuthentication bool `json:"disable_password_authentication"`
		SSHKeys                       struct {
			Path    string `json:"path" validate:"required"`
			KeyData string `json:"key_data"`
		} `json:"ssh_keys"`
	} `json:"os_profile_linux_config"`
	OSProfileSecrets struct {
		SourceVaultID           string `json:"source_vault_id"`
		SourceVaultCertificates struct {
			CertificateURL   string `json:"certificate_url"`
			CertificateStore string `json:"certificate_store"`
		} `json:"vault_certificates"`
	} `json:"os_profile_secrets"`
	NetworkInterfaceIDs []string          `json:"network_interface_ids"`
	Tags                map[string]string `json:"tags"`
	ClientID            string            `json:"azure_client_id"`
	ClientSecret        string            `json:"azure_client_secret"`
	TenantID            string            `json:"azure_tenant_id"`
	SubscriptionID      string            `json:"azure_subscription_id"`
	Environment         string            `json:"environment"`
	ErrorMessage        string            `json:"error,omitempty"`
	CryptoKey           string            `json:"-"`
}

// New : Constructor
func New(subject, cryptoKey string, body []byte, val *event.Validator) (event.Event, error) {
	var ev azure.Resource
	ev = &Event{CryptoKey: cryptoKey}
	if err := json.Unmarshal(body, &ev); err != nil {
		err := fmt.Errorf("Error on input message : %s", err)
		return nil, err
	}

	return azure.New(subject, "azurerm_virtual_machine", body, val, ev)
}

// SetID : id setter
func (ev *Event) SetID(id string) {
	ev.ID = id
}

// GetID : id getter
func (ev *Event) GetID() string {
	return ev.ID
}

// ResourceDataToEvent : Translates a ResourceData on a valid Ernest Event
func (ev *Event) ResourceDataToEvent(d *schema.ResourceData) error {
	ev.Name = d.Get("name").(string)
	ev.ResourceGroupName = d.Get("resource_group_name").(string)
	ev.Location = d.Get("location").(string)

	plan := d.Get("plan").(map[string]interface{})
	ev.Plan.Name = plan["name"].(string)
	ev.Plan.Publisher = plan["publisher"].(string)
	ev.Plan.Product = plan["product"].(string)

	ev.AvailabilitySetID = d.Get("availability_set_id").(string)
	ev.LicenseType = d.Get("license_type").(string)
	ev.VMSize = d.Get("vm_size").(string)

	storageImageReference := d.Get("storage_image_reference").(map[string]interface{})
	ev.StorageImageReference.Publisher = storageImageReference["publisher"].(string)
	ev.StorageImageReference.Offer = storageImageReference["offer"].(string)
	ev.StorageImageReference.Sku = storageImageReference["sku"].(string)
	ev.StorageImageReference.Version = storageImageReference["version"].(string)

	storageOSDisk := d.Get("storage_os_disk").(map[string]interface{})
	ev.StorageOSDisk.Name = storageOSDisk["name"].(string)
	ev.StorageOSDisk.VhdURI = storageOSDisk["vhd_uri"].(string)
	ev.StorageOSDisk.CreateOption = storageOSDisk["create_option"].(string)
	ev.StorageOSDisk.OSType = storageOSDisk["os_type"].(string)
	ev.StorageOSDisk.ImageURI = storageOSDisk["image_uri"].(string)
	ev.StorageOSDisk.Caching = storageOSDisk["caching"].(string)
	ev.DeleteOSDiskOnTermination = d.Get("delete_os_disk_on_termination").(bool)

	storageDataDisk := d.Get("storage_data_disk").(map[string]interface{})
	ev.StorageDataDisk.Name = storageDataDisk["name"].(string)
	ev.StorageDataDisk.VhdURI = storageDataDisk["vhd_uri"].(string)
	ev.StorageDataDisk.CreateOption = storageDataDisk["create_option"].(string)
	ev.StorageDataDisk.Size = storageDataDisk["size_size_go"].(int32)
	ev.StorageDataDisk.Lun = storageDataDisk["lun"].(int32)
	ev.DeleteDataDisksOnTermination = d.Get("delete_data_disks_on_termination").(bool)

	// TODO diagnostics_profile -> TypeSet

	bootDiagnostics := d.Get("boot_diagnostics").(map[string]interface{})
	ev.BootDiagnostics.Enabled = bootDiagnostics["enabled"].(bool)
	ev.BootDiagnostics.URI = bootDiagnostics["storage_uri"].(string)

	osProfile := d.Get("os_profile").(map[string]interface{})
	ev.OSProfile.ComputerName = osProfile["computer_name"].(string)
	ev.OSProfile.AdminUsername = osProfile["admin_username"].(string)
	ev.OSProfile.AdminPassword = osProfile["admin_password"].(string)
	ev.OSProfile.CustomData = osProfile["custom_data"].(string)

	win := d.Get("os_profile_windows_config").(map[string]interface{})
	ev.OSProfileWindowsConfig.ProvisionVMAgent = win["provision_vm_agent"].(bool)
	ev.OSProfileWindowsConfig.EnableAutomaticUpgrades = win["enable_automatic_upgrades"].(bool)
	winrm := win["win_rm"].(map[string]interface{})
	ev.OSProfileWindowsConfig.WinRm.Protocol = winrm["protocol"].(string)
	ev.OSProfileWindowsConfig.WinRm.CertificateURL = winrm["certificate_url"].(string)
	additional := win["additional_unattend_config"].(map[string]interface{})
	ev.OSProfileWindowsConfig.AdditionalUnattendConfig.Pass = additional["pass"].(string)
	ev.OSProfileWindowsConfig.AdditionalUnattendConfig.Component = additional["component"].(string)
	ev.OSProfileWindowsConfig.AdditionalUnattendConfig.SettingName = additional["setting_name"].(string)
	ev.OSProfileWindowsConfig.AdditionalUnattendConfig.Content = additional["content"].(string)

	lin := d.Get("os_profile_linux_config").(map[string]interface{})
	ev.OSProfileLinuxConfig.DisablePasswordAuthentication = lin["disable_password_authentication"].(bool)
	keys := lin["ssh_keys"].(map[string]interface{})
	ev.OSProfileLinuxConfig.SSHKeys.Path = keys["path"].(string)
	ev.OSProfileLinuxConfig.SSHKeys.KeyData = keys["key_data"].(string)

	sec := d.Get("os_profile_linux_config").(map[string]interface{})
	ev.OSProfileSecrets.SourceVaultID = sec["source_vault_id"].(string)
	certs := sec["vault_certificates"].(map[string]interface{})
	ev.OSProfileSecrets.SourceVaultCertificates.CertificateURL = certs["certificate_url"].(string)
	ev.OSProfileSecrets.SourceVaultCertificates.CertificateStore = certs["certificate_store"].(string)

	ev.NetworkInterfaceIDs = make([]string, 0)
	for _, id := range d.Get("network_interface_ids").(*schema.Set).List() {
		ev.NetworkInterfaceIDs = append(ev.NetworkInterfaceIDs, id.(string))
	}

	tags := make(map[string]string, 0)
	for k, v := range d.Get("tags").(map[string]interface{}) {
		tags[k] = v.(string)
	}
	ev.Tags = tags

	return nil
}

// EventToResourceData : Translates the current event on a valid ResourceData
func (ev *Event) EventToResourceData(d *schema.ResourceData) error {
	crypto := aes.New()

	encFields := make(map[string]string)
	encFields["subscription_id"] = ev.SubscriptionID
	encFields["client_id"] = ev.ClientID
	encFields["client_secret"] = ev.ClientSecret
	encFields["tenant_id"] = ev.TenantID
	encFields["environment"] = ev.Environment
	for k, v := range encFields {
		dec, err := crypto.Decrypt(v, ev.CryptoKey)
		if err != nil {
			err := fmt.Errorf("Field '%s' not valid : %s", k, err)
			ev.Log("error", err.Error())
			return err
		}
		if err := d.Set(k, dec); err != nil {
			err := fmt.Errorf("Field '%s' not valid : %s", k, err)
			ev.Log("error", err.Error())
			return err
		}
	}

	fields := make(map[string]interface{})
	fields["name"] = ev.Name
	fields["resource_group_name"] = ev.ResourceGroupName
	fields["location"] = ev.Location
	fields["availability_set_id"] = ev.AvailabilitySetID
	fields["license_type"] = ev.LicenseType
	fields["vm_size"] = ev.VMSize
	fields["storage_image_reference"] = ev.StorageImageReference
	fields["storage_os_disk"] = ev.StorageOSDisk
	fields["delete_os_disk_on_termination"] = ev.DeleteOSDiskOnTermination
	fields["storage_data_disk"] = ev.StorageDataDisk
	fields["delete_data_disks_on_termination"] = ev.DeleteDataDisksOnTermination
	fields["boot_diagnostics"] = ev.BootDiagnostics
	fields["os_profile"] = ev.OSProfile
	fields["os_profile_windows_config"] = ev.OSProfileWindowsConfig
	fields["os_profile_linux_config"] = ev.OSProfileLinuxConfig
	fields["os_profile_secrets"] = ev.OSProfileSecrets
	fields["network_interface_ids"] = ev.NetworkInterfaceIDs
	fields["tags"] = ev.Tags
	for k, v := range fields {
		if err := d.Set(k, v); err != nil {
			err := fmt.Errorf("Field '%s' not valid : %s", k, err)
			ev.Log("error", err.Error())
			return err
		}
	}

	return nil
}

// Error : will mark the event as errored
func (ev *Event) Error(err error) {
	ev.ErrorMessage = err.Error()
	ev.Body, err = json.Marshal(ev)
}
