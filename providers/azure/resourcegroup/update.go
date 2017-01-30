package resourcegroup

import (
	"errors"
	"fmt"

	"github.com/jen20/riviera/azure"
)

// Update : Updates a resource group object on azure
func (ev *Event) Update() error {
	client := ev.client()
	rivieraClient := client.RivieraClient

	updateRequest := rivieraClient.NewRequestForURI(ev.ID)
	updateRequest.Command = &azure.UpdateResourceGroup{
		Name: ev.Name,
		Tags: ev.Tags,
	}

	updateResponse, err := updateRequest.Execute()
	if err != nil {
		msg := fmt.Sprintf("Error updating resource group: %s", err)
		ev.Log("error", msg)
		return errors.New(msg)
	}
	if !updateResponse.IsSuccessful() {
		msg := fmt.Sprintf("Error updating resource group: %s", updateResponse.Error)
		ev.Log("error", msg)
		return errors.New(msg)
	}

	return ev.Get()
}
