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
	types "github.com/ernestio/ernestprovider/types/azure/virtualmachine"
	"github.com/ernestio/ernestprovider/validator"
	"github.com/fatih/structs"
)

// Event : This is the Ernest representation of an azure networkinterface
type Event struct {
	types.Event
	ErrorMessage string               `json:"error,omitempty" diff:"-"`
	CryptoKey    string               `json:"-" diff:"-"`
	Validator    *validator.Validator `json:"-" diff:"-"`
	GenericEvent event.Event          `json:"-" validate:"-" diff:"-"`
}

// New : Constructor
func New(subject, cryptoKey string, body []byte, val *validator.Validator) (event.Event, error) {
	ev := &Event{CryptoKey: cryptoKey, Validator: val}
	ev.Powered = true
	body = []byte(strings.Replace(string(body), `"_component":"virtual_machines"`, `"_component":"virtual_machine"`, 1))
	if err := json.Unmarshal(body, &ev); err != nil {
		err := fmt.Errorf("Error on input message : %s", err)
		return nil, err
	}

	ev.GenericEvent, _ = azure.New(subject, "azurerm_virtual_machine", body, val, ev)
	return ev.GenericEvent, nil
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

// SetState : state setter
func (ev *Event) SetState(state string) {
	ev.State = state
}

// ResourceDataToEvent : Translates a ResourceData on a valid Ernest Event
func (ev *Event) ResourceDataToEvent(d *schema.ResourceData) error {
	ev.ID = d.Id()
	if ev.ID == "" {
		ev.Name = d.Content["name"].(string)
	} else {
		parts := strings.Split(ev.ID, "/")
		ev.Name = parts[8]
	}
	ev.ComponentID = "virtual_machine::" + ev.Name
	ev.ResourceGroupName = d.Content["resource_group_name"].(string)
	ev.Location = d.Content["location"].(string)

	if d.Content["plan"] != nil {
		plan := d.Content["plan"].(*schema.Set).List()
		if len(plan) > 0 {
			planConfig := plan[0].(map[string]interface{})
			ev.Plan.Name = planConfig["name"].(string)
			ev.Plan.Publisher = planConfig["publisher"].(string)
			ev.Plan.Product = planConfig["product"].(string)
		}
	}

	if d.Content["availability_set_id"] != nil {
		ev.AvailabilitySetID = d.Content["availability_set_id"].(string)
	}
	if d.Content["license_type"] != nil {
		ev.LicenseType = d.Content["license_type"].(string)
	}
	if d.Content["vm_size"] != nil {
		ev.VMSize = fmt.Sprintf("%s", d.Content["vm_size"])
	}

	if d.Content["storage_image_reference"] != nil {
		storageImageReference := d.Content["storage_image_reference"].(*schema.Set).List()
		if len(storageImageReference) > 0 {
			sir := storageImageReference[0].(map[string]interface{})
			ev.StorageImageReference.Publisher = sir["publisher"].(string)
			ev.StorageImageReference.Offer = sir["offer"].(string)
			ev.StorageImageReference.Sku = sir["sku"].(string)
			ev.StorageImageReference.Version = sir["version"].(string)
		}
	}

	storageOSDisk := d.Content["storage_os_disk"].(*schema.Set).List()
	if len(storageOSDisk) > 0 {
		s := storageOSDisk[0].(map[string]interface{})
		ev.StorageOSDisk.Name = s["name"].(string)
		if s["vhd_uri"] != nil {
			ev.StorageOSDisk.VhdURI = s["vhd_uri"].(string)
		}
		ev.StorageOSDisk.CreateOption = fmt.Sprintf("%s", s["create_option"])

		if s["os_type"] != nil {
			ev.StorageOSDisk.OSType = s["os_type"].(string)
		}
		if s["image_uri"] != nil {
			ev.StorageOSDisk.ImageURI = s["image_uri"].(string)
		}
		if s["caching"] != nil {
			ev.StorageOSDisk.Caching = fmt.Sprintf("%s", s["caching"])
		}

		if s["managed_disk_type"] != nil {
			ev.StorageOSDisk.StorageAccountType = fmt.Sprintf("%s", s["managed_disk_type"])
		}
		if ev.StorageOSDisk.VhdURI == "" {
			ev.StorageOSDisk.ManagedDisk = s["name"].(string)
			parts := strings.Split(ev.ID, "/")
			parts[7] = "disks"
			parts[8] = ev.StorageOSDisk.Name
			ev.StorageOSDisk.ManagedDiskID = strings.Join(parts, "/")
		}
	}

	if d.Content["storage_data_disk"] != nil {
		storageDataDisk := d.Content["storage_data_disk"].([]interface{})
		if len(storageDataDisk) > 0 {
			s := storageDataDisk[0].(map[string]interface{})
			ev.StorageDataDisk.Name = s["name"].(string)
			ev.StorageDataDisk.VhdURI = s["vhd_uri"].(string)
			ev.StorageDataDisk.CreateOption = s["create_option"].(string)
			if ev.StorageDataDisk.VhdURI == "" {
				ev.StorageDataDisk.ManagedDisk = s["name"].(string)
				ev.StorageDataDisk.ManagedDiskID = s["managed_disk_id"].(string)
			}
			if s["disk_size_gb"] != nil {
				x := int32(s["disk_size_gb"].(int))
				ev.StorageDataDisk.Size = &x
			}
			if s["lun"] != nil {
				x := int32(s["lun"].(int))
				ev.StorageDataDisk.Lun = &x
			}
		}
	}

	bootDiagnostics := make([]types.BootDiagnostic, 0)
	if d.Content["boot_diagnostics"] != nil {
		for _, v := range d.Content["boot_diagnostics"].([]interface{}) {
			x := v.(map[string]interface{})
			bootDiagnostics = append(bootDiagnostics, types.BootDiagnostic{
				Enabled: x["enabled"].(bool),
				URI:     x["storage_uri"].(string),
			})
		}
	}
	ev.BootDiagnostics = bootDiagnostics

	if d.Content["os_profile"] != nil {
		osProfile := d.Content["os_profile"].(*schema.Set).List()
		if len(osProfile) > 0 {
			s := osProfile[0].(map[string]interface{})
			if s["computer_name"] != nil {
				ev.OSProfile.ComputerName = s["computer_name"].(string)
			}
			if s["admin_username"] != nil {
				ev.OSProfile.AdminUsername = s["admin_username"].(string)
			}
			if s["admin_password"] != nil {
				ev.OSProfile.AdminPassword = s["admin_password"].(string)
			}
			if s["custom_data"] != nil {
				ev.OSProfile.CustomData = s["custom_data"].(string)
			}
		}
	}

	if d.Content["os_profile_windows_config"] != nil {
		winList := d.Content["os_profile_windows_config"].([]interface{})
		if len(winList) > 0 {
			win := winList[0].(map[string]interface{})
			if win["provision_vm_agent"].(bool) != false && win["enable_automatic_upgrades"].(bool) != false {
				ev.OSProfileWindowsConfig = &types.OSProfileWindowsConfig{}

				ev.OSProfileWindowsConfig.ProvisionVMAgent = win["provision_vm_agent"].(bool)
				ev.OSProfileWindowsConfig.EnableAutomaticUpgrades = win["enable_automatic_upgrades"].(bool)

				if win["winrm"] != nil {
					rms := []types.WinRM{}
					for _, v := range win["winrm"].([]map[string]interface{}) {
						winrm := types.WinRM{}
						if val, ok := v["protocol"]; ok {
							winrm.Protocol = fmt.Sprintf("%s", val)
						}
						if val, ok := v["certificate_url"]; ok {
							winrm.CertificateURL = fmt.Sprintf("%s", val)
						}
						rms = append(rms, winrm)
					}
					ev.OSProfileWindowsConfig.WinRm = rms
				}
				if win["additional_unattend_config"] != nil {
					for i, value := range win["additional_unattend_config"].([]interface{}) {
						v := value.(map[string]interface{})
						ev.OSProfileWindowsConfig.AdditionalUnattendConfig[i].Pass = v["pass"].(string)
						ev.OSProfileWindowsConfig.AdditionalUnattendConfig[i].Component = v["component"].(string)
						ev.OSProfileWindowsConfig.AdditionalUnattendConfig[i].SettingName = v["setting_name"].(string)
						ev.OSProfileWindowsConfig.AdditionalUnattendConfig[i].Content = v["content"].(string)
					}
				}
			} else {
				ev.OSProfileWindowsConfig = nil
			}
		}
	}

	if d.Content["os_profile_linux_config"] != nil {
		linList := d.Content["os_profile_linux_config"].([]interface{})
		if len(linList) > 0 {
			lin := linList[0].(map[string]interface{})
			x := lin["disable_password_authentication"].(bool)
			ev.OSProfileLinuxConfig.DisablePasswordAuthentication = &x
			ev.OSProfileLinuxConfig.SSHKeys = make([]types.SSHKey, 0)
			if lin["ssh_keys"] != nil {
				for _, v := range lin["ssh_keys"].([]map[string]interface{}) {
					ev.OSProfileLinuxConfig.SSHKeys = append(ev.OSProfileLinuxConfig.SSHKeys, types.SSHKey{
						Path:    v["path"].(string),
						KeyData: v["key_data"].(string),
					})
				}
			}
		}
	}

	ev.OSProfileSecrets = make([]types.Secret, 0)
	for _, v := range d.Content["os_profile_secrets"].([]map[string]interface{}) {
		certs := []types.VaultCertificate{}
		for _, wal := range v["vault_certificates"].(*schema.Set).List() {
			w := wal.(map[string]interface{})
			certs = append(certs, types.VaultCertificate{
				CertificateURL:   w["certificate_url"].(string),
				CertificateStore: w["certificate_store"].(string),
			})
		}
		ev.OSProfileSecrets = append(ev.OSProfileSecrets, types.Secret{
			SourceVaultID:           v["source_vault_id"].(string),
			SourceVaultCertificates: certs,
		})
	}

	ev.NetworkInterfaceIDs = make([]string, 0)
	if d.Content["network_interface_ids"] != nil {
		for _, id := range d.Content["network_interface_ids"].([]string) {
			ev.NetworkInterfaceIDs = append(ev.NetworkInterfaceIDs, id)
		}
	}

	ev.DeleteOSDiskOnTermination = false
	ev.DeleteDataDisksOnTermination = false
	tags := make(map[string]string, 0)
	tt := d.Content["tags"].(map[string]interface{})
	for k, val := range tt {
		value := fmt.Sprintf("%s", val)
		if k == "delete_os_disk_on_termination" && value == "true" {
			ev.DeleteOSDiskOnTermination = true
		} else if k == "delete_data_disks_on_termination" && value == "true" {
			ev.DeleteDataDisksOnTermination = true
		} else {
			tags[k] = fmt.Sprintf("%s", value)
		}
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
	fields["powered"] = ev.Powered
	if ev.Plan.Name != "" && ev.Plan.Product != "" && ev.Plan.Publisher != "" {
		fields["plan"] = []interface{}{structs.Map(ev.Plan)}
	}
	fields["availability_set_id"] = ev.AvailabilitySetID
	fields["license_type"] = ev.LicenseType
	fields["vm_size"] = ev.VMSize
	fields["storage_image_reference"] = []interface{}{structs.Map(ev.StorageImageReference)}
	fields["storage_os_disk"] = []interface{}{structs.Map(ev.StorageOSDisk)}
	fields["delete_data_disks_on_termination"] = ev.DeleteDataDisksOnTermination
	fields["delete_os_disk_on_termination"] = ev.DeleteOSDiskOnTermination
	fields["os_profile"] = []interface{}{structs.Map(ev.OSProfile)}

	if ev.StorageDataDisk.Size != nil {
		ddisk := make(map[string]interface{})
		ddisk["name"] = ev.StorageDataDisk.Name
		ddisk["vhd_uri"] = ev.StorageDataDisk.VhdURI
		ddisk["create_option"] = ev.StorageDataDisk.CreateOption
		ddisk["disk_size_gb"] = *ev.StorageDataDisk.Size
		if ddisk["vhd_uri"] == "" {
			ddisk["managed_disk_id"] = ev.StorageDataDisk.ManagedDiskID
		}
		if ev.StorageDataDisk.Lun != nil {
			ddisk["lun"] = *ev.StorageDataDisk.Lun
		}
		fields["storage_data_disk"] = []interface{}{ddisk}
	}

	var diagnostics []interface{}
	for _, bd := range ev.BootDiagnostics {
		diagnostics = append(diagnostics, structs.Map(bd))
	}

	fields["boot_diagnostics"] = diagnostics
	if ev.OSProfileWindowsConfig != nil {
		fields["os_profile_windows_config"] = []interface{}{structs.Map(ev.OSProfileWindowsConfig)}
	}

	secrets := make([]interface{}, 0)
	for _, v := range ev.OSProfileSecrets {
		secrets = append(secrets, structs.Map(v))
	}
	fields["os_profile_secrets"] = secrets

	lconfig := make(map[string]interface{})
	if ev.OSProfileLinuxConfig.DisablePasswordAuthentication != nil {
		lconfig["disable_password_authentication"] = *ev.OSProfileLinuxConfig.DisablePasswordAuthentication
	}
	var sshkeys []interface{}
	for i := range ev.OSProfileLinuxConfig.SSHKeys {
		ev.OSProfileLinuxConfig.SSHKeys[i].Path = strings.Replace(ev.OSProfileLinuxConfig.SSHKeys[i].Path, "\\u003c", "<", -1)
		ev.OSProfileLinuxConfig.SSHKeys[i].KeyData = strings.Replace(ev.OSProfileLinuxConfig.SSHKeys[i].KeyData, "\\u003c", "<", -1)
		ev.OSProfileLinuxConfig.SSHKeys[i].Path = strings.Replace(ev.OSProfileLinuxConfig.SSHKeys[i].Path, "\\u003e", ">", -1)
		ev.OSProfileLinuxConfig.SSHKeys[i].KeyData = strings.Replace(ev.OSProfileLinuxConfig.SSHKeys[i].KeyData, "\\u003e", ">", -1)
		sshkeys = append(sshkeys, structs.Map(ev.OSProfileLinuxConfig.SSHKeys[i]))
	}
	lconfig["ssh_keys"] = sshkeys

	fields["network_interface_ids"] = ev.NetworkInterfaceIDs
	if ev.OSProfileLinuxConfig.DisablePasswordAuthentication != nil || len(ev.OSProfileLinuxConfig.SSHKeys) > 0 {
		fields["os_profile_linux_config"] = []interface{}{lconfig}
	}

	if ev.DeleteDataDisksOnTermination {
		ev.Tags["delete_data_disks_on_termination"] = "true"
	}
	if ev.DeleteOSDiskOnTermination {
		ev.Tags["delete_os_disk_on_termination"] = "true"
	}
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

// Clone : will mark the event as errored
func (ev *Event) Clone() (event.Event, error) {
	body, _ := json.Marshal(ev)
	return New(ev.Subject, ev.CryptoKey, body, ev.Validator)
}

// Error : will mark the event as errored
func (ev *Event) Error(err error) {
	ev.ErrorMessage = err.Error()
	ev.Body, err = json.Marshal(ev)
}
