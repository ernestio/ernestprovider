package resourcegroup

import (
	"errors"
	"fmt"

	"github.com/jen20/riviera/azure"
)

// Get : Gets a resource group object on azure
func (ev *Event) Get() error {
	client := ev.client()
	rivieraClient := client.RivieraClient

	readRequest := rivieraClient.NewRequestForURI(ev.ID)
	readRequest.Command = &azure.GetResourceGroup{}

	readResponse, err := readRequest.Execute()
	if err != nil {
		return fmt.Errorf("Error reading resource group: %s", err)
	}
	if !readResponse.IsSuccessful() {
		msg := "Error reading resource group " + ev.ID + " - removing from state"
		ev.Log("info", msg)
		ev.ID = ""
		return errors.New(msg)
	}

	resp := readResponse.Parsed.(*azure.GetResourceGroupResponse)

	ev.Name = *resp.Name
	ev.Location = *resp.Location
	ev.Tags = *resp.Tags

	return nil
}
