/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package virtualmachine

import (
	"fmt"
	"net/url"
	"strings"
)

// Delete : Deletes a network interface on azure
func (ev *Event) Delete() error {
	vmClient := ev.client()

	resGroup := ev.ResourceGroupName
	name := ev.Name

	if _, err := vmClient.Delete(resGroup, name, make(chan struct{})); err != nil {
		return err
	}

	// delete OS Disk if opted in
	if ev.DeleteOSDiskOnTermination == true {
		ev.Log("info", "delete_os_disk_on_termination is enabled, deleting")

		osDisk, err := ev.expandAzureRmVirtualMachineOsDisk()
		if err != nil {
			err := fmt.Errorf("Error expanding OS Disk: %s", err)
			ev.Log("error", err.Error())
			return err
		}

		if err = ev.resourceArmVirtualMachineDeleteVhd(*osDisk.Vhd.URI, ev.ResourceGroupName); err != nil {
			err := fmt.Errorf("Error deleting OS Disk VHD: %s", err)
			ev.Log("error", err.Error())
			return err
		}
	}

	// delete Data disks if opted in
	if ev.DeleteDataDisksOnTermination {
		ev.Log("info", "delete_data_disks_on_termination is enabled, deleting each data disk")

		disks, err := ev.expandAzureRmVirtualMachineDataDisk()
		if err != nil {
			err := fmt.Errorf("Error expanding Data Disks: %s", err)
			ev.Log("error", err.Error())
			return err
		}

		for _, disk := range disks {
			if err = ev.resourceArmVirtualMachineDeleteVhd(*disk.Vhd.URI, resGroup); err != nil {
				err := fmt.Errorf("Error deleting Data Disk VHD: %s", err)
				ev.Log("error", err.Error())
				return err
			}
		}
	}
	return nil
}

func (ev *Event) resourceArmVirtualMachineDeleteVhd(uri, resGroup string) error {
	vhdURL, err := url.Parse(uri)
	if err != nil {
		return fmt.Errorf("Cannot parse Disk VHD URI: %s", err)
	}

	// VHD URI is in the form: https://storageAccountName.blob.core.windows.net/containerName/blobName
	storageAccountName := strings.Split(vhdURL.Host, ".")[0]
	path := strings.Split(strings.TrimPrefix(vhdURL.Path, "/"), "/")
	containerName := path[0]
	blobName := path[1]

	blobClient, saExists, err := ev.ArmClient.GetBlobStorageClientForStorageAccount(resGroup, storageAccountName)
	if err != nil {
		return fmt.Errorf("Error creating blob store client for VHD deletion: %s", err)
	}

	if !saExists {
		ev.Log("info", fmt.Sprintf("Storage Account %q doesn't exist so the VHD blob won't exist", storageAccountName))
		return nil
	}

	ev.Log("info", fmt.Sprintf("[INFO] Deleting VHD blob %s", blobName))
	_, err = blobClient.DeleteBlobIfExists(containerName, blobName, nil)
	if err != nil {
		return fmt.Errorf("Error deleting VHD blob: %s", err)
	}

	return nil
}
