package virtualmachine

import (
	"strings"
	"testing"

	"github.com/ernestio/ernestprovider/event"
)

func validEvent() Event {
	tags := make(map[string]*string)
	one := "one"
	tags["t1"] = &one

	return Event{
		Name:              "supu",
		ResourceGroupName: "resource_group",
		Location:          "westus",
		Tags:              tags,
		Validator:         event.NewValidator(),
	}
}

func TestRequiredName(t *testing.T) {
	ev := validEvent()
	ev.Name = ""
	err := ev.Validate()

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
	err := ev.Validate()

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
	err := ev.Validate()

	if err == nil {
		t.Error("No error has been received!")
	}

	if !strings.Contains(err.Error(), "ResourceGroupName is a required field") {
		t.Error("Output message does not contain name or required strings")
	}
}

func TestHappyPath(t *testing.T) {
	ev := validEvent()

	err := ev.Validate()
	if err != nil {
		println(err.Error())
		t.Error("I'm in a bad mood.")
	}
}
