package resourcegroup

import (
	"errors"
	"fmt"

	"github.com/jen20/riviera/azure"
)

// Delete : Deletes a resource group object on azure
func (ev *Event) Delete() error {
	client := ev.client()
	rivieraClient := client.RivieraClient

	deleteRequest := rivieraClient.NewRequestForURI(ev.ID)
	deleteRequest.Command = &azure.DeleteResourceGroup{}

	deleteResponse, err := deleteRequest.Execute()
	if err != nil {
		msg := fmt.Sprintf("Error deleting resource group: %s", err)
		ev.Log("error", msg)
		return errors.New(msg)
	}
	if !deleteResponse.IsSuccessful() {
		msg := fmt.Sprintf("Error deleting resource group: %s", deleteResponse.Error)
		ev.Log("error", msg)
		return errors.New(msg)
	}

	return nil
}
