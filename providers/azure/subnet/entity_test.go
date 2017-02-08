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
	}
}

func TestRequiredName(t *testing.T) {
	ev := validEvent()
	ev.Name = ""
	val := event.NewValidator()
	err := val.Validate(ev)

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
	val := event.NewValidator()
	err := val.Validate(ev)

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
	val := event.NewValidator()
	err := val.Validate(ev)

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
	val := event.NewValidator()
	err := val.Validate(ev)

	if err == nil {
		t.Error("No error has been received!")
	}

	if !strings.Contains(err.Error(), "AddressPrefix is a required field") {
		t.Error("Output message does not contain name or required strings")
	}
}

func TestHappyPath(t *testing.T) {
	ev := validEvent()

	val := event.NewValidator()
	err := val.Validate(ev)
	if err != nil {
		println(err.Error())
		t.Error("I'm in a bad mood.")
	}
}
