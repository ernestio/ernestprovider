/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package virtualmachine

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/hashicorp/terraform/builtin/providers/azurerm"
	"github.com/hashicorp/terraform/helper/schema"

	aes "github.com/ernestio/crypto/aes"
	"github.com/ernestio/ernestprovider/event"
)

// Event : This is the Ernest representation of an azure networkinterface
type Event struct {
	event.Base
	ID                           string   `json:"id"`
	Name                         string   `json:"name" validate:"required"`
	ResourceGroupName            string   `json:"resource_group_name" validate:"required"`
	Location                     string   `json:"location" validate:"required"`
	AvailabilitySetID            string   `json:"availability_set_id"`
	LicenseType                  string   `json:"license_type"`
	VMSize                       string   `json:"vm_size"`
	DeleteOSDiskOnTermination    bool     `json:"delete_os_disk_on_termination"`
	DeleteDataDisksOnTermination bool     `json:"delete_data_disks_on_termination"`
	NetworkInterfaceIDs          []string `json:"network_interface_ids"`
	Plans                        []struct {
		Name      string `json:"name" validate:"required"`
		Publisher string `json:"publisher" validate:"required"`
		Product   string `json:"product" validate:"required"`
	} `json:"plan" validate:"dive"`
	StorageImageReferences []struct {
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
	StorageDataDisks []struct {
		Name         string `json:"name" validate:"required"`
		VhdURI       string `json:"vhd_uri" validate:"required"`
		CreateOption string `json:"create_option" validate:"required"`
		Size         int32  `json:"disk_size_gb" validate:"required"`
		Lun          int32  `json:"lun"`
	} `json:"storage_data_disk"`
	DiagnosticsProfiles []struct {
		BootDiagnostics []struct {
			Enabled bool   `json:"enabled"`
			URI     string `json:"storage_uri"`
		} `json:"boot_diagnostics"`
	} `json:"diagnostics_profile"`
	OSProfiles []struct {
		ComputerName  string `json:"computer_name" validation:"required"`
		AdminUsername string `json:"admin_username" validation:"required"`
		AdminPassword string `json:"admin_password" validation:"required"`
		CustomData    string `json:"custom_data"`
	} `json:"os_profile"`
	OSProfileWindowsConfigs []struct {
		ProvisionVMAgent        bool `json:"provision_vm_agent"`
		EnableAutomaticUpgrades bool `json:"enable_automatic_upgrades"`
		WinRm                   []struct {
			Protocol       string `json:"protocol" validate:"required"`
			CertificateURL string `json:"certificate_url" validate:"required"`
		} `json:"winrm"`
		AdditionalUnattendConfig []struct {
			Pass        string `json:"pass" validate:"required"`
			Component   string `json:"component" validate:"required"`
			SettingName string `json:"setting_name" validate:"required"`
			Content     string `json:"content" validate:"required"`
		} `json:"additional_unattend_config"`
	} `json:"os_profile_windows_config"`
	OSProfileLinuxConfigs []struct {
		DisablePasswordAuthentication bool `json:"disable_password_authentication"`
		SSHKeys                       []struct {
			Path    string `json:"path" validate:"required"`
			KeyData string `json:"key_data"`
		} `json:"ssh_keys"`
	} `json:"os_profile_linux_config"`
	OSProfileSecrets []struct {
		SourceVaultID           string `json:"source_vault_id"`
		SourceVaultCertificates []struct {
			CertificateURL   string `json:"certificate_url"`
			CertificateStore string `json:"certificate_store"`
		} `json:"vault_certificates"`
	} `json:"os_profile_secrets"`
	Tags map[string]string `json:"tags"`

	ClientID       string `json:"azure_client_id"`
	ClientSecret   string `json:"azure_client_secret"`
	TenantID       string `json:"azure_tenant_id"`
	SubscriptionID string `json:"azure_subscription_id"`
	Environment    string `json:"environment"`

	Provider     *schema.Provider
	Component    *schema.Resource
	ResourceData *schema.ResourceData
	Schema       map[string]*schema.Schema
	ErrorMessage string           `json:"error,omitempty"`
	Subject      string           `json:"-"`
	Body         []byte           `json:"-"`
	CryptoKey    string           `json:"-"`
	Validator    *event.Validator `json:"-"`
}

// New : Constructor
func New(subject string, body []byte, cryptoKey string, val *event.Validator) (event.Event, error) {
	var err error
	n := Event{Subject: subject, Body: body, CryptoKey: cryptoKey, Validator: val}
	n.Provider = azurerm.Provider().(*schema.Provider)
	n.Component = n.Provider.ResourcesMap["azurerm_virtual_machine"]
	n.Schema = n.schema()
	n.Body = body
	n.Subject = subject
	n.CryptoKey = cryptoKey
	n.Validator = val
	if n.ResourceData, err = n.toResourceData(body); err != nil {
		n.Log("error", err.Error())
		return &n, err
	}

	return &n, nil
}

// Validate checks if all criteria are met
func (ev *Event) Validate() error {
	return ev.Validator.Validate(ev)
}

// Find : Find an object on azure
func (ev *Event) Find() error {
	return errors.New(ev.Subject + " not supported")
}

// Create : Creates a Virtual Network on Azure using terraform
// providers
func (ev *Event) Create() error {
	c, err := ev.client()
	if err != nil {
		return err
	}
	if err := ev.Component.Create(ev.ResourceData, c); err != nil {
		err := fmt.Errorf("Error creating the requestd resource : %s", err)
		ev.Log("error", err.Error())
		return err
	}

	return nil
}

// Update : Updates an existing Virtual Network on Azure
// by using azurerm terraform provider resource
func (ev *Event) Update() error {
	c, err := ev.client()
	if err != nil {
		return err
	}
	if err := ev.Component.Update(ev.ResourceData, c); err != nil {
		err := fmt.Errorf("Error creating the requestd resource : %s", err)
		ev.Log("error", err.Error())
		return err
	}

	return nil
}

// Get : Requests and loads the resource to Azure through azurerm
// terraform provider
func (ev *Event) Get() error {
	c, err := ev.client()
	if err != nil {
		return err
	}
	if err := ev.Component.Read(ev.ResourceData, c); err != nil {
		err := fmt.Errorf("Error getting virtual network : %s", err)
		ev.Log("error", err.Error())
		return err
	}

	ev.toEvent()
	return nil
}

// Delete : Deletes the received resource from azure through
// azurerm terraform provider
func (ev *Event) Delete() error {
	c, err := ev.client()
	if err != nil {
		return err
	}
	if err := ev.Component.Delete(ev.ResourceData, c); err != nil {
		err := fmt.Errorf("Error deleting the requested resource : %s", err)
		ev.Log("error", err.Error())
		return err
	}

	return nil
}

// GetBody : Gets the body for this event
func (ev *Event) GetBody() []byte {
	var err error
	if ev.Body, err = json.Marshal(ev); err != nil {
		log.Println(err.Error())
	}
	return ev.Body
}

// GetSubject : Gets the subject for this event
func (ev *Event) GetSubject() string {
	return ev.Subject
}

// Process : starts processing the current message
func (ev *Event) Process() (err error) {
	if err := json.Unmarshal(ev.Body, &ev); err != nil {
		ev.Error(err)
		return err
	}

	return nil
}

// Error : Will respond the current event with an error
func (ev *Event) Error(err error) {
	log.Printf("Error: %s", err.Error())
	ev.ErrorMessage = err.Error()

	ev.Body, err = json.Marshal(ev)
}

// Translates a ResourceData on a valid Ernest Event
func (ev *Event) toEvent() {
	ev.Name = ev.ResourceData.Get("name").(string)
	ev.ResourceGroupName = ev.ResourceData.Get("resource_group_name").(string)
	ev.Location = ev.ResourceData.Get("location").(string)
	ev.AvailabilitySetID = ev.ResourceData.Get("availability_set_id").(string)
	ev.LicenseType = ev.ResourceData.Get("license_type").(string)
	ev.VMSize = ev.ResourceData.Get("vm_size").(string)
	ev.DeleteOSDiskOnTermination = ev.ResourceData.Get("delete_os_disk_on_termination").(bool)
	ev.Tags = ev.ResourceData.Get("tags").(map[string]string)
}

// Translates the current event on a valid ResourceData
func (ev *Event) toResourceData(body []byte) (*schema.ResourceData, error) {
	var d schema.ResourceData
	d.SetSchema(ev.Schema)
	if err := json.Unmarshal(body, &ev); err != nil {
		err := fmt.Errorf("Error on input message : %s", err)
		ev.Log("error", err.Error())
		return nil, err
	}

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
			return nil, err
		}
		if err := d.Set(k, dec); err != nil {
			err := fmt.Errorf("Field '%s' not valid : %s", k, err)
			ev.Log("error", err.Error())
			return nil, err
		}
	}

	fields := make(map[string]interface{})
	fields["name"] = ev.Name
	fields["location"] = ev.Location
	fields["tags"] = ev.Tags
	for k, v := range fields {
		if err := d.Set(k, v); err != nil {
			err := fmt.Errorf("Field '%s' not valid : %s", k, err)
			ev.Log("error", err.Error())
			return nil, err
		}
	}

	return &d, nil
}

// Based on the Provider and Component schemas it calculates
// the necessary schema to be create a new ResourceData
func (ev *Event) schema() (sch map[string]*schema.Schema) {
	if ev.Schema != nil {
		return ev.Schema
	}
	a := ev.Provider.Schema
	b := ev.Component.Schema
	sch = a
	for k, v := range b {
		sch[k] = v
	}
	return sch
}

// Azure virtual network client
func (ev *Event) client() (*azurerm.ArmClient, error) {
	client, err := ev.Provider.ConfigureFunc(ev.ResourceData)
	if err != nil {
		err := fmt.Errorf("Can't connect to provider : %s", err)
		ev.Log("error", err.Error())
		return nil, err
	}
	c := client.(*azurerm.ArmClient)
	return c, nil
}
