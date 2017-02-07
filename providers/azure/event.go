package azure

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/ernestio/ernestprovider/event"
	"github.com/fatih/color"
	"github.com/hashicorp/terraform/builtin/providers/azurerm"
	"github.com/hashicorp/terraform/helper/schema"
)

// Resource : ...
type Resource interface {
	SetID(id string)
	GetID() string
	ResourceDataToEvent(d *schema.ResourceData) error
	EventToResourceData(d *schema.ResourceData) error
}

// Event : ...
type Event struct {
	event.Base

	Resource Resource

	Provider     *schema.Provider
	Component    *schema.Resource
	ResourceData *schema.ResourceData
	Schema       map[string]*schema.Schema
	Validator    *event.Validator

	Subject string `json:"-"`
	Body    []byte `json:"-"`
}

// New : Constructor
func New(subject, resourceType string, body []byte, val *event.Validator, res Resource) (event.Event, error) {
	n := Event{Subject: subject, Body: body, Validator: val}
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

	return &n, nil
}

// Validate checks if all criteria are met
func (ev *Event) Validate() error {
	return ev.Validator.Validate(ev.Resource)
}

// Find : Find an object on azure
func (ev *Event) Find() error {
	return errors.New(ev.Subject + " not supported")
}

// Create : Creates a Resource Group on Azure using terraform
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
	ev.ResourceData.SetId(ev.Resource.GetID())
	if err := ev.Component.Read(ev.ResourceData, c); err != nil {
		err := fmt.Errorf("Error getting resource group : %s", err)
		ev.Log("error", err.Error())
		ev.Log("debug", "Original message: "+string(ev.Body))
		return err
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
		panic(err)
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
