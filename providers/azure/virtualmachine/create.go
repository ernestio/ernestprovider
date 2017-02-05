/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package virtualmachine

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/arm/compute"
)

// Create : Creates a storage account on azure
func (ev *Event) Create() error {
	ev.Log("info", "preparing arguments for Azure ARM Virtual Machine creation.")

	osDisk, err := ev.expandAzureRmVirtualMachineOsDisk()
	if err != nil {
		ev.Log("error", err.Error())
		return err
	}
	storageProfile := compute.StorageProfile{
		OsDisk: osDisk,
	}

	if len(ev.StorageImageReferences) > 0 {
		imageRef, err := ev.expandAzureRmVirtualMachineImageReference()
		if err != nil {
			return err
		}
		storageProfile.ImageReference = imageRef
	}

	if len(ev.StorageDataDisks) < 0 {
		dataDisks, err := ev.expandAzureRmVirtualMachineDataDisk()
		if err != nil {
			return err
		}
		storageProfile.DataDisks = &dataDisks
	}

	networkProfile := ev.expandAzureRmVirtualMachineNetworkProfile()
	properties := compute.VirtualMachineProperties{
		NetworkProfile: &networkProfile,
		HardwareProfile: &compute.HardwareProfile{
			VMSize: compute.VirtualMachineSizeTypes(ev.VMSize),
		},
		StorageProfile: &storageProfile,
	}

	if len(ev.DiagnosticsProfiles) > 0 {
		diagnosticsProfile := ev.expandAzureRmVirtualMachineDiagnosticsProfile()
		properties.DiagnosticsProfile = &diagnosticsProfile
	}

	osProfile, err := ev.expandAzureRmVirtualMachineOsProfile()
	if err != nil {
		return err
	}
	properties.OsProfile = osProfile

	if ev.AvailabilitySetID != "" {
		availSet := compute.SubResource{
			ID: &ev.AvailabilitySetID,
		}

		properties.AvailabilitySet = &availSet
	}

	vm := compute.VirtualMachine{
		Name:                     &ev.Name,
		Location:                 &ev.Location,
		VirtualMachineProperties: &properties,
		Tags: &ev.Tags,
	}

	if len(ev.Plans) > 0 {
		plan, err := ev.expandAzureRmVirtualMachinePlan()
		if err != nil {
			return err
		}

		vm.Plan = plan
	}

	client := ev.client()
	_, vmErr := client.CreateOrUpdate(ev.ResourceGroupName, ev.Name, vm, make(chan struct{}))
	if vmErr != nil {
		return vmErr
	}

	read, err := client.Get(ev.ResourceGroupName, ev.Name, "")
	if err != nil {
		return err
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read Virtual Machine %s (resource group %s) ID", ev.Name, ev.ResourceGroupName)
	}

	ev.ID = *read.ID

	return ev.Get()
}

func (ev *Event) expandAzureRmVirtualMachineOsDisk() (*compute.OSDisk, error) {
	disk := ev.StorageOSDisk

	osDisk := &compute.OSDisk{
		Name: &disk.Name,
		Vhd: &compute.VirtualHardDisk{
			URI: &disk.VhdURI,
		},
		CreateOption: compute.DiskCreateOptionTypes(disk.CreateOption),
	}

	if disk.ImageURI != "" {
		osDisk.Image = &compute.VirtualHardDisk{
			URI: &disk.ImageURI,
		}
	}

	if disk.OSType != "" {
		if disk.OSType == "linux" {
			osDisk.OsType = compute.Linux
		} else if disk.OSType == "windows" {
			osDisk.OsType = compute.Windows
		} else {
			return nil, fmt.Errorf("os_type must be 'linux' or 'windows'")
		}
	}

	if disk.Caching != "" {
		osDisk.Caching = compute.CachingTypes(disk.Caching)
	}

	return osDisk, nil
}

func (ev *Event) expandAzureRmVirtualMachineNetworkProfile() compute.NetworkProfile {
	nicIDs := ev.NetworkInterfaceIDs
	networkInterfaces := make([]compute.NetworkInterfaceReference, 0, len(nicIDs))

	networkProfile := compute.NetworkProfile{}

	for _, id := range nicIDs {
		networkInterface := compute.NetworkInterfaceReference{
			ID: &id,
		}
		networkInterfaces = append(networkInterfaces, networkInterface)
	}

	networkProfile.NetworkInterfaces = &networkInterfaces

	return networkProfile
}

func (ev *Event) expandAzureRmVirtualMachineImageReference() (*compute.ImageReference, error) {
	si := ev.StorageImageReferences[0]

	return &compute.ImageReference{
		Publisher: &si.Publisher,
		Offer:     &si.Offer,
		Sku:       &si.Sku,
		Version:   &si.Version,
	}, nil
}

func (ev *Event) expandAzureRmVirtualMachineDataDisk() ([]compute.DataDisk, error) {
	disks := ev.StorageDataDisks
	dataDisks := make([]compute.DataDisk, 0, len(disks))
	for _, config := range disks {
		dataDisk := compute.DataDisk{
			Name: &config.Name,
			Vhd: &compute.VirtualHardDisk{
				URI: &config.VhdURI,
			},
			Lun:          &config.Lun,
			DiskSizeGB:   &config.Size,
			CreateOption: compute.DiskCreateOptionTypes(config.CreateOption),
		}

		dataDisks = append(dataDisks, dataDisk)
	}

	return dataDisks, nil
}

func (ev *Event) expandAzureRmVirtualMachineDiagnosticsProfile() compute.DiagnosticsProfile {
	diagnosticsProfile := ev.DiagnosticsProfiles[0]
	bootDiagnostics := diagnosticsProfile.BootDiagnostics[0]
	enabled := bootDiagnostics.Enabled
	storageURI := bootDiagnostics.URI

	return compute.DiagnosticsProfile{
		BootDiagnostics: &compute.BootDiagnostics{
			Enabled:    &enabled,
			StorageURI: &storageURI,
		},
	}
}

func (ev *Event) expandAzureRmVirtualMachineOsProfile() (*compute.OSProfile, error) {
	osProfile := ev.OSProfiles[0]

	adminUsername := osProfile.AdminUsername
	adminPassword := osProfile.AdminPassword
	computerName := osProfile.ComputerName

	profile := &compute.OSProfile{
		AdminUsername: &adminUsername,
		ComputerName:  &computerName,
	}

	if adminPassword != "" {
		profile.AdminPassword = &adminPassword
	}

	if len(ev.OSProfileWindowsConfigs) > 0 {
		winConfig, err := ev.expandAzureRmVirtualMachineOsProfileWindowsConfig()
		if err != nil {
			return nil, err
		}
		if winConfig != nil {
			profile.WindowsConfiguration = winConfig
		}
	}

	if len(ev.OSProfileLinuxConfigs) > 0 {
		linuxConfig, err := ev.expandAzureRmVirtualMachineOsProfileLinuxConfig()
		if err != nil {
			return nil, err
		}
		if linuxConfig != nil {
			profile.LinuxConfiguration = linuxConfig
		}
	}

	if len(ev.OSProfileSecrets) > 0 {
		secrets := ev.expandAzureRmVirtualMachineOsProfileSecrets()
		if secrets != nil {
			profile.Secrets = secrets
		}
	}

	profile.CustomData = &osProfile.CustomData

	return profile, nil
}

func (ev *Event) expandAzureRmVirtualMachineOsProfileWindowsConfig() (*compute.WindowsConfiguration, error) {

	osProfileConfig := ev.OSProfileWindowsConfigs[0]
	config := &compute.WindowsConfiguration{}
	config.ProvisionVMAgent = &osProfileConfig.ProvisionVMAgent
	config.EnableAutomaticUpdates = &osProfileConfig.EnableAutomaticUpgrades

	winRm := osProfileConfig.WinRm
	if len(winRm) > 0 {
		winRmListners := make([]compute.WinRMListener, 0, len(winRm))
		for _, config := range winRm {
			protocol := config.Protocol
			winRmListner := compute.WinRMListener{
				Protocol: compute.ProtocolTypes(protocol),
			}
			winRmListner.CertificateURL = &config.CertificateURL

			winRmListners = append(winRmListners, winRmListner)
		}
		config.WinRM = &compute.WinRMConfiguration{
			Listeners: &winRmListners,
		}
	}

	if len(osProfileConfig.AdditionalUnattendConfig) > 0 {
		additionalConfigContent := make([]compute.AdditionalUnattendContent, 0, len(osProfileConfig.AdditionalUnattendConfig))
		for _, config := range osProfileConfig.AdditionalUnattendConfig {
			addContent := compute.AdditionalUnattendContent{
				PassName:      compute.PassNames(config.Pass),
				ComponentName: compute.ComponentNames(config.Component),
				SettingName:   compute.SettingNames(config.SettingName),
				Content:       &config.Content,
			}

			additionalConfigContent = append(additionalConfigContent, addContent)
		}
		config.AdditionalUnattendContent = &additionalConfigContent
	}
	return config, nil
}

func (ev *Event) expandAzureRmVirtualMachineOsProfileLinuxConfig() (*compute.LinuxConfiguration, error) {
	linuxConfig := ev.OSProfileLinuxConfigs[0]
	disablePasswordAuth := linuxConfig.DisablePasswordAuthentication

	config := &compute.LinuxConfiguration{
		DisablePasswordAuthentication: &disablePasswordAuth,
	}

	linuxKeys := linuxConfig.SSHKeys
	sshPublicKeys := []compute.SSHPublicKey{}
	for _, sshKey := range linuxKeys {
		sshPublicKey := compute.SSHPublicKey{
			Path:    &sshKey.Path,
			KeyData: &sshKey.KeyData,
		}

		sshPublicKeys = append(sshPublicKeys, sshPublicKey)
	}

	config.SSH = &compute.SSHConfiguration{
		PublicKeys: &sshPublicKeys,
	}

	return config, nil
}

func (ev *Event) expandAzureRmVirtualMachineOsProfileSecrets() *[]compute.VaultSecretGroup {
	secretsConfig := ev.OSProfileSecrets
	secrets := make([]compute.VaultSecretGroup, 0, len(secretsConfig))

	for _, config := range secretsConfig {
		sourceVaultID := config.SourceVaultID

		vaultSecretGroup := compute.VaultSecretGroup{
			SourceVault: &compute.SubResource{
				ID: &sourceVaultID,
			},
		}

		certsConfig := config.SourceVaultCertificates
		certs := make([]compute.VaultCertificate, 0, len(certsConfig))
		for _, cc := range certsConfig {
			cert := compute.VaultCertificate{
				CertificateURL: &cc.CertificateURL,
			}
			cert.CertificateStore = &cc.CertificateStore
			certs = append(certs, cert)
		}
		vaultSecretGroup.VaultCertificates = &certs

		secrets = append(secrets, vaultSecretGroup)
	}

	return &secrets
}

func (ev *Event) expandAzureRmVirtualMachinePlan() (*compute.Plan, error) {
	plan := ev.Plans[0]

	return &compute.Plan{
		Publisher: &plan.Publisher,
		Name:      &plan.Name,
		Product:   &plan.Product,
	}, nil
}
