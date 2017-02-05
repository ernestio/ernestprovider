/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package virtualmachine

import (
	"fmt"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/arm/compute"
)

// Get : Gets a storage account on azure
func (ev *Event) Get() error {
	resGroup := ev.ResourceGroupName
	name := ev.Name

	resp, err := ev.client().Get(resGroup, name, "")

	if err != nil {
		err := fmt.Errorf("Error making Read request on Azure Virtual Machine %s: %s", name, err)
		ev.Log("error", err.Error())
		return err
	}
	if resp.StatusCode == http.StatusNotFound {
		ev.ID = ""
		return nil
	}

	if resp.Plan != nil {
		ev.Plans[0].Name = *resp.Plan.Name
		ev.Plans[0].Publisher = *resp.Plan.Publisher
		ev.Plans[0].Product = *resp.Plan.Product
	}

	if resp.VirtualMachineProperties.AvailabilitySet != nil {
		ev.AvailabilitySetID = *resp.VirtualMachineProperties.AvailabilitySet.ID
	}

	ev.VMSize = string(resp.VirtualMachineProperties.HardwareProfile.VMSize)

	if resp.VirtualMachineProperties.StorageProfile.ImageReference != nil {
		ev.StorageImageReferences = append(ev.StorageImageReferences, resp.VirtualMachineProperties.StorageProfile.ImageReference)
	}

	if err := d.Set("storage_os_disk", schema.NewSet(resourceArmVirtualMachineStorageOsDiskHash, flattenAzureRmVirtualMachineOsDisk(resp.Properties.StorageProfile.OsDisk))); err != nil {
		return fmt.Errorf("[DEBUG] Error setting Virtual Machine Storage OS Disk error: %#v", err)
	}

	if resp.Properties.StorageProfile.DataDisks != nil {
		if err := d.Set("storage_data_disk", flattenAzureRmVirtualMachineDataDisk(resp.Properties.StorageProfile.DataDisks)); err != nil {
			return fmt.Errorf("[DEBUG] Error setting Virtual Machine Storage Data Disks error: %#v", err)
		}
	}

	if err := d.Set("os_profile", schema.NewSet(resourceArmVirtualMachineStorageOsProfileHash, flattenAzureRmVirtualMachineOsProfile(resp.Properties.OsProfile))); err != nil {
		return fmt.Errorf("[DEBUG] Error setting Virtual Machine Storage OS Profile: %#v", err)
	}

	if resp.Properties.OsProfile.WindowsConfiguration != nil {
		if err := d.Set("os_profile_windows_config", flattenAzureRmVirtualMachineOsProfileWindowsConfiguration(resp.Properties.OsProfile.WindowsConfiguration)); err != nil {
			return fmt.Errorf("[DEBUG] Error setting Virtual Machine Storage OS Profile Windows Configuration: %#v", err)
		}
	}

	if resp.Properties.OsProfile.LinuxConfiguration != nil {
		if err := d.Set("os_profile_linux_config", flattenAzureRmVirtualMachineOsProfileLinuxConfiguration(resp.Properties.OsProfile.LinuxConfiguration)); err != nil {
			return fmt.Errorf("[DEBUG] Error setting Virtual Machine Storage OS Profile Linux Configuration: %#v", err)
		}
	}

	if resp.Properties.OsProfile.Secrets != nil {
		if err := d.Set("os_profile_secrets", flattenAzureRmVirtualMachineOsProfileSecrets(resp.Properties.OsProfile.Secrets)); err != nil {
			return fmt.Errorf("[DEBUG] Error setting Virtual Machine Storage OS Profile Secrets: %#v", err)
		}
	}

	if resp.Properties.DiagnosticsProfile != nil {
		if err := d.Set("diagnostics_profile", flattenAzureRmVirtualMachineDiagnosticsProfile(resp.Properties.DiagnosticsProfile)); err != nil {
			return fmt.Errorf("[DEBUG] Error setting Virtual Machine Diagnostics Profile: %#v", err)
		}
	}

	if resp.Properties.NetworkProfile != nil {
		if err := d.Set("network_interface_ids", flattenAzureRmVirtualMachineNetworkInterfaces(resp.Properties.NetworkProfile)); err != nil {
			return fmt.Errorf("[DEBUG] Error setting Virtual Machine Storage Network Interfaces: %#v", err)
		}
	}

	flattenAndSetTags(d, resp.Tags)

	return nil
}

func (ev *Event) flattenAzureRmVirtualMachinePlan(plan *compute.Plan) map[string]interface{} {
	result := make(map[string]interface{})
	result["name"] = *plan.Name
	result["publisher"] = *plan.Publisher
	result["product"] = *plan.Product

	return result
}
