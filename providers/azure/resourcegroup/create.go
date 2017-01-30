package resourcegroup

import (
	"errors"
	"fmt"

	"github.com/jen20/riviera/azure"
)

// Create : Creates a resource group object on azure
func (ev *Event) Create() error {
	rivieraClient := ev.client().RivieraClient

	createRequest := rivieraClient.NewRequest()
	createRequest.Command = &azure.CreateResourceGroup{
		Name:     ev.Name,
		Location: ev.Location,
		Tags:     ev.Tags,
	}

	createResponse, err := createRequest.Execute()
	if err != nil {
		msg := fmt.Sprintf("Error creating resource group: %s", err.Error())
		ev.Log("error", msg)
		return errors.New(msg)
	}
	if !createResponse.IsSuccessful() {
		msg := fmt.Sprintf("Error creating resource group: %s", createResponse.Error)
		ev.Log("error", msg)
		return errors.New(msg)
	}

	resp := createResponse.Parsed.(*azure.CreateResourceGroupResponse)
	ev.ID = *resp.ID

	return ev.Get()
}
