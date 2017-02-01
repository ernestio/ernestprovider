package networkinterface

import (
	"strings"
	"testing"

	"github.com/ernestio/ernestprovider/event"
)

func validEvent() Event {
	var dns []string
	var ips []IPConfiguration

	tags := make(map[string]*string)
	one := "one"
	tags["t1"] = &one

	dns = append(dns, "8.8.8.8")
	dns = append(dns, "4.4.4.4")

	ips = append(ips, IPConfiguration{
		Name:                       "ip",
		Subnet:                     "10.0.2.0/14",
		PrivateIPAddress:           "10.0.2.1",
		PrivateIPAddressAllocation: "10.0.2.1",
	})

	return Event{
		Name:              "supu",
		ResourceGroupName: "resource_group",
		Location:          "westus",
		IPConfigurations:  ips,
		DNSServers:        dns,
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

func TestRequiredIPConfigurations(t *testing.T) {
	var ips []IPConfiguration
	ev := validEvent()
	ev.IPConfigurations = ips
	err := ev.Validate()

	if err == nil {
		t.Error("No error has been received!")
	}

	if !strings.Contains(err.Error(), "IPConfigurations is a required field") {
		t.Error("Output message does not contain name or required strings")
	}
}

func TestRequiredIPConfigurationsName(t *testing.T) {
	ev := validEvent()
	ev.IPConfigurations[0].Name = ""
	err := ev.Validate()

	if err == nil {
		t.Error("No error has been received!")
	}

	if !strings.Contains(err.Error(), "Name is a required field") {
		t.Error("Output message does not contain name or required strings")
	}
}

func TestRequiredIPConfigurationsSubnet(t *testing.T) {
	var str string
	ev := validEvent()
	ev.IPConfigurations[0].Subnet = str
	err := ev.Validate()

	if err == nil {
		t.Error("No error has been received!")
	}

	if !strings.Contains(err.Error(), "Subnet is a required field") {
		t.Error("Output message does not contain name or required strings")
	}
}

func TestRequiredIPConfigurationsPrivateIPAddressAllocation(t *testing.T) {
	var str string
	ev := validEvent()
	ev.IPConfigurations[0].PrivateIPAddressAllocation = str
	err := ev.Validate()

	if err == nil {
		t.Error("No error has been received!")
	}

	if !strings.Contains(err.Error(), "PrivateIPAddressAllocation is a required field") {
		t.Error("Output message does not contain name or required strings")
	}
}

func TestInvalidDNS(t *testing.T) {
	ev := validEvent()
	ev.DNSServers[0] = "no ip"
	err := ev.Validate()

	if err == nil {
		t.Error("No error has been received!")
	}

	if !strings.Contains(err.Error(), "DNSServers[0] must be a valid IP address") {
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
