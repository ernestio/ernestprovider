package azure

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/ernestio/ernestprovider/event"
	"github.com/fatih/color"
	"github.com/r3labs/terraform/builtin/providers/azurerm"
	"github.com/r3labs/terraform/helper/schema"
)

// Event : ...
type Event struct {
	event.Base

	Resource event.Resource

	Provider     *schema.Provider
	Component    *schema.Resource
	ResourceData *schema.ResourceData
	ResourceType string
	Schema       map[string]*schema.Schema
	Validator    *event.Validator

	Subject string `json:"-"`
	Body    []byte `json:"-"`
}

// New : Constructor
func New(subject, resourceType string, body []byte, val *event.Validator, res event.Resource) (event.Event, error) {
	var n Event
	n = Event{Subject: subject, Body: body, Validator: val}
	n.Provider = azurerm.Provider().(*schema.Provider)
	n.Component = n.Provider.ResourcesMap[resourceType]
	n.Schema = n.schema()

	var d schema.ResourceData
	d.SetSchema(n.Schema)
	n.ResourceData = &d

	if err := res.EventToResourceData(n.ResourceData); err != nil {
		return nil, err
	}

	n.Body = body
	n.Subject = subject
	n.Validator = val
	n.Resource = res
	n.ResourceType = resourceType

	return &n, nil
}

// Validate checks if all criteria are met
func (ev *Event) Validate() error {
	return ev.Validator.Validate(ev.Resource)
}

// SetID ....
func (ev *Event) SetID(id string) {
	ev.Resource.SetID(id)
}

// GetResource ....
func (ev *Event) GetResource() event.Resource {
	return ev.Resource
}

// SetState ...
func (ev *Event) SetState(state string) {
	ev.Resource.SetState(state)
}

// SetComponents ...
func (ev *Event) SetComponents(components []event.Event) {
	ev.Resource.SetComponents(components)
}

func (ev *Event) getResourceGroup() (string, error) {
	var rg struct {
		ResourceGroup string `json:"resource_group_name"`
	}
	ev.Log("debug", "Subject : "+string(ev.GetSubject()))
	body := ev.GetBody()
	ev.Log("debug", "Body : "+string(body))
	if err := json.Unmarshal(body, &rg); err != nil {
		ev.Log("error", err.Error())
		return "", err
	}
	return rg.ResourceGroup, nil
}

// Find : Find an object on azure
func (ev *Event) Find() error {
	cli, _ := ev.client()
	res, err := ev.getResourceGroup()
	if err != nil {
		return err
	}
	results, err := cli.ListResourcesByGroup(res, "", "")
	if err != nil {
		ev.Log("error", err.Error())
	}
	components := make([]event.Event, 0)
	for _, ids := range results {
		for _, id := range ids {
			if ev.Resource.ValidateID(id) {
				e, err := ev.GetResource().Clone()
				if err != nil {
					log.Println(err)
				}
				e.SetID(id)
				e.SetState("completed")
				if err := e.Get(); err != nil {
					ev.Log("error", err.Error())
				}

				components = append(components, e)
			}
		}
	}
	ev.SetState("completed")
	ev.SetComponents(components)

	return nil
}

// Create : Creates a Resource Group on Azure using terraform
// providers
func (ev *Event) Create() error {
	c, err := ev.client()
	if err != nil {
		return err
	}
	if err := ev.Component.Create(ev.ResourceData, c); err != nil {
		err := fmt.Errorf("Error creating the requested resource : %s", err)
		ev.Log("error", err.Error())
		return err
	}
	ev.Resource.SetID(ev.ResourceData.Id())
	ev.Log("debug", "Created resource group : "+ev.Resource.GetID())

	return nil
}

// Update : Updates an existing Resource Group on Azure
// by using azurerm terraform provider resource
func (ev *Event) Update() error {
	c, err := ev.client()
	if err != nil {
		return err
	}
	ev.ResourceData.SetId(ev.Resource.GetID())
	st := reflect.TypeOf(ev.Component)
	if _, ok := st.MethodByName("Update"); !ok {
		println("Not supported")
		// return errors.New("Not supported")
	}
	if err := ev.Component.Update(ev.ResourceData, c); err != nil {
		err := fmt.Errorf("Error creating the requested resource : %s", err)
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
	id := ev.Resource.GetID()
	ev.ResourceData.SetId(id)
	parts := strings.Split(id, "/")
	if parts[1] == "subscriptions" {
		if err := ev.Component.Read(ev.ResourceData, c); err != nil {
			err := fmt.Errorf("Resource not found : %s", err)
			ev.Log("error", err.Error())
			ev.Log("debug", "Original message: "+string(ev.Body))
			return err
		}
		if ev.ResourceData.Id() == "" {
			err := fmt.Errorf("Resource not found")
			ev.Log("warn", err.Error())
			return err
		}
	}

	if err = ev.Resource.ResourceDataToEvent(ev.ResourceData); err != nil {
		ev.Log("error", err.Error())
		return err
	}
	return nil
}

// Delete : Deletes the received resource from azure through
// azurerm terraform provider
func (ev *Event) Delete() error {
	c, err := ev.client()
	if err != nil {
		return err
	}
	ev.ResourceData.SetId(ev.Resource.GetID())
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
	if ev.Body, err = json.Marshal(ev.Resource); err != nil {
		ev.Log("error", err.Error())
	}
	return ev.Body
}

func (ev *Event) getStatedBody(state string) []byte {
	ev.Resource.SetState(state)
	return ev.GetBody()
}

// GetErroredBody : Gets the body overloaded with errored state
func (ev *Event) GetErroredBody() []byte {
	return ev.getStatedBody("errored")
}

// GetCompletedBody : Gets the body overloaded with errored state
func (ev *Event) GetCompletedBody() []byte {
	return ev.getStatedBody("completed")
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

// Clone :
func (ev *Event) Clone() (event.Event, error) {
	return nil, errors.New("Not supported")
}

// Error : Will respond the current event with an error
func (ev *Event) Error(err error) {
	log.Printf("Error: %s", err.Error())
	ev.Resource.Error(err)
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

// Log : ...
func (ev *Event) Log(level, message string) {
	red := color.New(color.FgRed).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	blue := color.New(color.FgBlue).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()

	switch level {
	case "error":
		level = red("[ERROR]")
	case "warn":
		level = yellow("[WARNING]")
	case "info":
		level = blue("[INFO]")
	case "debug":
		level = green("[DEBUG]")
	}

	log.Println(level, message)

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
