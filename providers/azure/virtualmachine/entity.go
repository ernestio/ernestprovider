/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package virtualmachine

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/ernestio/ernestprovider/event"
	"github.com/ernestio/ernestprovider/providers/azure"

	"github.com/Azure/azure-sdk-for-go/arm/compute"
)

// Event : This is the Ernest representation of an azure networkinterface
type Event struct {
	event.Base
	ID                string `json:"id"`
	Name              string `json:"name" validate:"required"`
	ResourceGroupName string `json:"resource_group_name" validate:"required"`
	Location          string `json:"location" validate:"required"`
	Plans             []struct {
		Name      string `json:"name" validate:"required"`
		Publisher string `json:"publisher" validate:"required"`
		Product   string `json:"product" validate:"required"`
	} `json:"plan" validate:"dive"`
	AvailabilitySetID      string `json:"availability_set_id"`
	LicenseType            string `json:"license_type"`
	VMSize                 string `json:"vm_size"`
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
	DeleteOSDiskOnTermination bool `json:"delete_os_disk_on_termination"`
	StorageDataDisks          []struct {
		Name         string `json:"name" validate:"required"`
		VhdURI       string `json:"vhd_uri" validate:"required"`
		CreateOption string `json:"create_option" validate:"required"`
		Size         int32  `json:"disk_size_gb" validate:"required"`
		Lun          int32  `json:"lun"`
	} `json:"storage_data_disk"`
	DeleteDataDisksOnTermination bool `json:"delete_data_disks_on_termination"`
	DiagnosticsProfiles          []struct {
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
	NetworkInterfaceIDs []string           `json:"network_interface_ids"`
	Tags                map[string]*string `json:"tags"`

	ClientID       string `json:"azure_client_id"`
	ClientSecret   string `json:"azure_client_secret"`
	TenantID       string `json:"azure_tenant_id"`
	SubscriptionID string `json:"azure_subscription_id"`
	Environment    string `json:"environment"`

	ErrorMessage string           `json:"error,omitempty"`
	Subject      string           `json:"-"`
	Body         []byte           `json:"-"`
	CryptoKey    string           `json:"-"`
	Validator    *event.Validator `json:"-"`
	ArmClient    *azure.ArmClient `json:"-"`
}

// New : Constructor
func New(subject string, body []byte, cryptoKey string, val *event.Validator) event.Event {
	n := Event{Subject: subject, Body: body, CryptoKey: cryptoKey, Validator: val}

	return &n
}

// Azure client
func (ev *Event) client() *compute.VirtualMachinesClient {
	client, err := azure.Provider(ev.SubscriptionID, ev.ClientID, ev.ClientSecret, ev.TenantID, ev.Environment, ev.CryptoKey)
	if err != nil {
		panic(err)
	}

	ev.ArmClient = client
	return &client.VMClient
}

// Validate checks if all criteria are met
func (ev *Event) Validate() error {
	return ev.Validator.Validate(ev)
}

// Find : Find an object on azure
func (ev *Event) Find() error {
	return errors.New(ev.Subject + " not supported")
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
