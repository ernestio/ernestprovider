package subnet

import (
	"strings"
	"testing"

	"github.com/ernestio/ernestprovider/event"
)

func validEvent() Event {
	var configs []string
	configs = append(configs, "A")
	configs = append(configs, "B")

	return Event{
		Name:                 "supu",
		ResourceGroupName:    "group",
		VirtualNetworkName:   "vn_name",
		AddressPrefix:        "prefix",
		NetworkSecurityGroup: "sec",
		RouteTable:           "route",
		IPConfigurations:     configs,
		Validator:            event.NewValidator(),
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

func TestRequiredVirtualNetworkName(t *testing.T) {
	ev := validEvent()
	ev.VirtualNetworkName = ""
	err := ev.Validate()

	if err == nil {
		t.Error("No error has been received!")
	}

	if !strings.Contains(err.Error(), "VirtualNetworkName is a required field") {
		t.Error("Output message does not contain name or required strings")
	}
}

func TestRequiredAddressPrefix(t *testing.T) {
	ev := validEvent()
	ev.AddressPrefix = ""
	err := ev.Validate()

	if err == nil {
		t.Error("No error has been received!")
	}

	if !strings.Contains(err.Error(), "AddressPrefix is a required field") {
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
