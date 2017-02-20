/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package virtualmachine

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/r3labs/terraform/helper/schema"

	aes "github.com/ernestio/crypto/aes"
	"github.com/ernestio/ernestprovider/event"
	"github.com/ernestio/ernestprovider/providers/azure"
	"github.com/fatih/structs"
)

// Event : This is the Ernest representation of an azure networkinterface
type Event struct {
	event.Base
	ID                string `json:"id"`
	Name              string `json:"name" validate:"required"`
	ResourceGroupName string `json:"resource_group_name" validate:"required"`
	Location          string `json:"location" validate:"required"`
	Plan              struct {
		Name      string `json:"name"`
		Publisher string `json:"publisher"`
		Product   string `json:"product"`
	} `json:"plan"`
	AvailabilitySetID     string `json:"availability_set_id"`
	LicenseType           string `json:"license_type"`
	VMSize                string `json:"vm_size"`
	StorageImageReference struct {
		Publisher string `json:"publisher" validate:"required" structs:"publisher"`
		Offer     string `json:"offer" validate:"required" structs:"offer"`
		Sku       string `json:"sku" validate:"required" structs:"sku"`
		Version   string `json:"version" structs:"version"`
	} `json:"storage_image_reference" validate:"dive"`
	StorageOSDisk struct {
		Name         string `json:"name" validate:"required" structs:"name"`
		VhdURI       string `json:"vhd_uri" validate:"required" structs:"vhd_uri"`
		CreateOption string `json:"create_option" validate:"required" structs:"create_option"`
		OSType       string `json:"os_type" structs:"os_type"`
		ImageURI     string `json:"image_uri" structs:"image_uri"`
		Caching      string `json:"caching" structs:"caching"`
	} `json:"storage_os_disk" validate:"dive"`
	DeleteOSDiskOnTermination bool `json:"delete_os_disk_on_termination"`
	StorageDataDisk           struct {
		Name         string `json:"name" validate:"required" structs:"name"`
		VhdURI       string `json:"vhd_uri" validate:"required" structs:"vhd_uri"`
		CreateOption string `json:"create_option" validate:"required" structs:"create_option"`
		Size         int32  `json:"disk_size_gb" validate:"required" structs:"disk_size_gb"`
		Lun          int32  `json:"lun" structs:"lun"`
	} `json:"storage_data_disk"`
	DeleteDataDisksOnTermination bool             `json:"delete_data_disks_on_termination"`
	BootDiagnostics              []bootDiagnostic `json:"boot_diagnostics"`
	OSProfile                    struct {
		ComputerName  string `json:"computer_name" structs:"computer_name"`
		AdminUsername string `json:"admin_username" structs:"admin_username"`
		AdminPassword string `json:"admin_password" structs:"admin_password"`
		CustomData    string `json:"custom_data" structs:"custom_data"`
	} `json:"os_profile"`
	OSProfileWindowsConfig struct {
		ProvisionVMAgent        bool `json:"provision_vm_agent" structs:"provision_vm_agent"`
		EnableAutomaticUpgrades bool `json:"enable_automatic_upgrades" structs:"enable_automatic_upgrades"`
		WinRm                   []struct {
			Protocol       string `json:"protocol" structs:"protocol"`
			CertificateURL string `json:"certificate_url" structs:"certification_url"`
		} `json:"winrm" structs:"winrm"`
		AdditionalUnattendConfig []struct {
			Pass        string `json:"pass" structs:"pass"`
			Component   string `json:"component" structs:"component"`
			SettingName string `json:"setting_name" structs:"setting_name"`
			Content     string `json:"content" structs:"content"`
		} `json:"additional_unattend_config" structs:"additional_unattend_config"`
	} `json:"os_profile_windows_config"`
	OSProfileLinuxConfig struct {
		DisablePasswordAuthentication bool     `json:"disable_password_authentication" structs:"disable_password_authentication"`
		SSHKeys                       []sshKey `json:"ssh_keys" structs:"ssh_keys"`
	} `json:"os_profile_linux_config"`
	OSProfileSecrets    []secret          `json:"os_profile_secrets"`
	NetworkInterfaceIDs []string          `json:"network_interface_ids"`
	Tags                map[string]string `json:"tags"`
	ClientID            string            `json:"azure_client_id"`
	ClientSecret        string            `json:"azure_client_secret"`
	TenantID            string            `json:"azure_tenant_id"`
	SubscriptionID      string            `json:"azure_subscription_id"`
	Environment         string            `json:"environment"`
	ErrorMessage        string            `json:"error,omitempty"`
	Components          []json.RawMessage `json:"components"`
	CryptoKey           string            `json:"-"`
}

type secret struct {
	SourceVaultID           string             `json:"source_vault_id" structs:"source_vault_id"`
	SourceVaultCertificates []vaultCertificate `json:"vault_certificates" structs:"vault_certificates"`
}

type vaultCertificate struct {
	CertificateURL   string `json:"certificate_url" structs:"certificate_url"`
	CertificateStore string `json:"certificate_store" structs:"certificate_store"`
}
type sshKey struct {
	Path    string `json:"path" validate:"required" structs:"path"`
	KeyData string `json:"key_data" structs:"key_data"`
}

type bootDiagnostic struct {
	Enabled bool   `json:"enabled"`
	URI     string `json:"storage_uri"`
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

// SetComponents : ....
func (ev *Event) SetComponents(components []event.Event) {
	for _, v := range components {
		ev.Components = append(ev.Components, v.GetBody())
	}
}

// ValidateID : determines if the given id is valid for this resource type
func (ev *Event) ValidateID(id string) bool {
	parts := strings.Split(strings.ToLower(id), "/")
	if len(parts) != 9 {
		return false
	}
	if parts[6] != "microsoft.compute" {
		return false
	}
	if parts[7] != "virtualmachines" {
		return false
	}
	return true
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

	plan := d.Get("plan").(*schema.Set).List()
	if len(plan) > 0 {
		planConfig := plan[0].(map[string]interface{})
		ev.Plan.Name = planConfig["name"].(string)
		ev.Plan.Publisher = planConfig["publisher"].(string)
		ev.Plan.Product = planConfig["product"].(string)
	}

	ev.AvailabilitySetID = d.Get("availability_set_id").(string)
	ev.LicenseType = d.Get("license_type").(string)
	ev.VMSize = d.Get("vm_size").(string)

	storageImageReference := d.Get("storage_image_reference").(*schema.Set).List()
	if len(storageImageReference) > 0 {
		sir := storageImageReference[0].(map[string]interface{})
		ev.StorageImageReference.Publisher = sir["publisher"].(string)
		ev.StorageImageReference.Offer = sir["offer"].(string)
		ev.StorageImageReference.Sku = sir["sku"].(string)
		ev.StorageImageReference.Version = sir["version"].(string)
	}

	storageOSDisk := d.Get("storage_os_disk").(*schema.Set).List()
	if len(storageOSDisk) > 0 {
		s := storageOSDisk[0].(map[string]interface{})
		ev.StorageOSDisk.Name = s["name"].(string)
		ev.StorageOSDisk.VhdURI = s["vhd_uri"].(string)
		ev.StorageOSDisk.CreateOption = s["create_option"].(string)
		ev.StorageOSDisk.OSType = s["os_type"].(string)
		ev.StorageOSDisk.ImageURI = s["image_uri"].(string)
		ev.StorageOSDisk.Caching = s["caching"].(string)
	}
	ev.DeleteOSDiskOnTermination = d.Get("delete_os_disk_on_termination").(bool)

	storageDataDisk := d.Get("storage_data_disk").([]interface{})
	if len(storageDataDisk) > 0 {
		s := storageDataDisk[0].(map[string]interface{})
		ev.StorageDataDisk.Name = s["name"].(string)
		ev.StorageDataDisk.VhdURI = s["vhd_uri"].(string)
		ev.StorageDataDisk.CreateOption = s["create_option"].(string)
		if s["disk_size_gb"] != nil {
			ev.StorageDataDisk.Size = int32(s["disk_size_gb"].(int))
		}
		if s["lun"] != nil {
			ev.StorageDataDisk.Lun = int32(s["lun"].(int))
		}
	}
	ev.DeleteDataDisksOnTermination = d.Get("delete_data_disks_on_termination").(bool)

	// TODO diagnostics_profile -> TypeSet

	bootDiagnostics := make([]bootDiagnostic, 0)
	for _, v := range d.Get("boot_diagnostics").([]interface{}) {
		x := v.(map[string]interface{})
		bootDiagnostics = append(bootDiagnostics, bootDiagnostic{
			Enabled: x["enabled"].(bool),
			URI:     x["storage_uri"].(string),
		})
	}
	ev.BootDiagnostics = bootDiagnostics

	osProfile := d.Get("os_profile").(*schema.Set).List()
	if len(osProfile) > 0 {
		s := osProfile[0].(map[string]interface{})
		ev.OSProfile.ComputerName = s["computer_name"].(string)
		ev.OSProfile.AdminUsername = s["admin_username"].(string)
		ev.OSProfile.AdminPassword = s["admin_password"].(string)
		ev.OSProfile.CustomData = s["custom_data"].(string)
	}

	winList := d.Get("os_profile_windows_config").(*schema.Set).List()
	if len(winList) > 0 {
		win := winList[0].(map[string]interface{})
		ev.OSProfileWindowsConfig.ProvisionVMAgent = win["provision_vm_agent"].(bool)
		ev.OSProfileWindowsConfig.EnableAutomaticUpgrades = win["enable_automatic_upgrades"].(bool)

		for i, v := range win["win_rm"].([]map[string]interface{}) {
			ev.OSProfileWindowsConfig.WinRm[i].Protocol = v["protocol"].(string)
			ev.OSProfileWindowsConfig.WinRm[i].CertificateURL = v["certificate_url"].(string)
		}

		for i, v := range win["additional_unattend_config"].([]map[string]interface{}) {
			ev.OSProfileWindowsConfig.AdditionalUnattendConfig[i].Pass = v["pass"].(string)
			ev.OSProfileWindowsConfig.AdditionalUnattendConfig[i].Component = v["component"].(string)
			ev.OSProfileWindowsConfig.AdditionalUnattendConfig[i].SettingName = v["setting_name"].(string)
			ev.OSProfileWindowsConfig.AdditionalUnattendConfig[i].Content = v["content"].(string)
		}
	}

	linList := d.Get("os_profile_linux_config").(*schema.Set).List()
	if len(linList) > 0 {
		lin := linList[0].(map[string]interface{})
		ev.OSProfileLinuxConfig.DisablePasswordAuthentication = lin["disable_password_authentication"].(bool)
		ev.OSProfileLinuxConfig.SSHKeys = make([]sshKey, 0)
		for _, key := range lin["ssh_keys"].([]interface{}) {
			v := key.(map[string]interface{})
			ev.OSProfileLinuxConfig.SSHKeys = append(ev.OSProfileLinuxConfig.SSHKeys, sshKey{
				Path:    v["path"].(string),
				KeyData: v["key_data"].(string),
			})
		}
	}

	ev.OSProfileSecrets = make([]secret, 0)
	for _, val := range d.Get("os_profile_secrets").(*schema.Set).List() {
		v := val.(map[string]interface{})
		certs := []vaultCertificate{}
		for _, wal := range v["vault_certificates"].(*schema.Set).List() {
			w := wal.(map[string]interface{})
			certs = append(certs, vaultCertificate{
				CertificateURL:   w["certificate_url"].(string),
				CertificateStore: w["certificate_store"].(string),
			})
		}
		ev.OSProfileSecrets = append(ev.OSProfileSecrets, secret{
			SourceVaultID:           v["source_vault_id"].(string),
			SourceVaultCertificates: certs,
		})
	}

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
	fields["storage_image_reference"] = []interface{}{structs.Map(ev.StorageImageReference)}
	fields["storage_os_disk"] = []interface{}{structs.Map(ev.StorageOSDisk)}
	fields["delete_os_disk_on_termination"] = ev.DeleteOSDiskOnTermination
	fields["storage_data_disk"] = []interface{}{structs.Map(ev.StorageDataDisk)}
	fields["delete_data_disks_on_termination"] = ev.DeleteDataDisksOnTermination
	fields["boot_diagnostics"] = ev.BootDiagnostics
	fields["os_profile"] = []interface{}{structs.Map(ev.OSProfile)}

	if len(ev.OSProfileWindowsConfig.WinRm) > 0 {
		fields["os_profile_windows_config"] = []interface{}{structs.Map(ev.OSProfileWindowsConfig)}
	}
	fields["os_profile_linux_config"] = []interface{}{structs.Map(ev.OSProfileLinuxConfig)}
	secrets := make([]interface{}, 0)
	for _, v := range ev.OSProfileSecrets {
		secrets = append(secrets, structs.Map(v))
	}

	fields["os_profile_secrets"] = secrets
	fields["network_interface_ids"] = ev.NetworkInterfaceIDs
	fields["tags"] = ev.Tags
	for k, v := range fields {
		if err := d.Set(k, v); err != nil {
			ev.Log("error", err.Error())
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
