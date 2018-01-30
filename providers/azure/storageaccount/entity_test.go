package storageaccount

import (
	"strings"
	"testing"

	types "github.com/ernestio/ernestprovider/types/azure/storageaccount"
	"github.com/ernestio/ernestprovider/validator"
)

func validEvent() Event {
	tags := make(map[string]string)
	tags["t1"] = "one"

	return Event{
		Event: types.Event{
			Name:              "supu",
			ResourceGroupName: "resource_group",
			Location:          "westus",
			AccountType:       "atype",
			Tags:              tags,
		},
	}
}

func TestRequiredName(t *testing.T) {
	ev := validEvent()
	ev.Name = ""
	val := validator.NewValidator()
	err := val.Validate(ev)

	if err == nil {
		t.Error("No error has been received!")
	}

	if !strings.Contains(err.Error(), "Name is a required field") {
		t.Error("Output message does not contain name or required strings")
	}
}

func TestRequiredLocation(t *testing.T) {
	ev := validEvent()
	ev.Location = ""
	val := validator.NewValidator()
	err := val.Validate(ev)

	if err == nil {
		t.Error("No error has been received!")
	}

	if !strings.Contains(err.Error(), "Location is a required field") {
		t.Error("Output message does not contain name or required strings")
	}
}

func TestRequiredResourceGroupName(t *testing.T) {
	ev := validEvent()
	ev.ResourceGroupName = ""
	val := validator.NewValidator()
	err := val.Validate(ev)

	if err == nil {
		t.Error("No error has been received!")
	}

	if !strings.Contains(err.Error(), "ResourceGroupName is a required field") {
		t.Error("Output message does not contain name or required strings")
	}
}

func TestRequiredAccounttype(t *testing.T) {
	ev := validEvent()
	ev.AccountType = ""
	val := validator.NewValidator()
	err := val.Validate(ev)

	if err == nil {
		t.Error("No error has been received!")
	}

	if !strings.Contains(err.Error(), "AccountType is a required field") {
		t.Error("Output message does not contain name or required strings")
	}
}

func TestHappyPath(t *testing.T) {
	ev := validEvent()

	val := validator.NewValidator()
	err := val.Validate(ev)
	if err != nil {
		println(err.Error())
		t.Error("I'm in a bad mood.")
	}
}
