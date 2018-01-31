package networkinterface

import (
	"strings"
	"testing"

	types "github.com/ernestio/ernestprovider/types/azure/networkinterface"
	"github.com/ernestio/ernestprovider/validator"
)

func validEvent() Event {
	var dns []string
	var ips []types.IPConfiguration

	tags := make(map[string]string)
	tags["t1"] = "one"

	dns = append(dns, "8.8.8.8")
	dns = append(dns, "4.4.4.4")

	ips = append(ips, types.IPConfiguration{
		Name:                       "ip",
		SubnetID:                   "10.0.2.0/14",
		PrivateIPAddress:           "10.0.2.1",
		PrivateIPAddressAllocation: "10.0.2.1",
	})

	return Event{
		Event: types.Event{
			Name:              "supu",
			ResourceGroupName: "resource_group",
			Location:          "westus",
			IPConfigurations:  ips,
			DNSServers:        dns,
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

/*
func TestRequiredLocation(t *testing.T) {
	ev := validEvent()
	ev.Location = ""

	val := event.NewValidator()
	err := val.Validate(ev)

	if err == nil {
		t.Error("No error has been received!")
	}

	if !strings.Contains(err.Error(), "Location is a required field") {
		t.Error("Output message does not contain name or required strings")
	}
}
*/

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

/*
func TestRequiredIPConfigurations(t *testing.T) {
	var ips []IPConfiguration
	ev := validEvent()
	ev.IPConfigurations = ips

	val := event.NewValidator()
	err := val.Validate(ev)

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

	val := event.NewValidator()
	err := val.Validate(ev)

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
	ev.IPConfigurations[0].SubnetID = str

	val := event.NewValidator()
	err := val.Validate(ev)

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

	val := event.NewValidator()
	err := val.Validate(ev)

	if err == nil {
		t.Error("No error has been received!")
	}

	if !strings.Contains(err.Error(), "PrivateIPAddressAllocation is a required field") {
		t.Error("Output message does not contain name or required strings")
	}
}
*/

func TestInvalidDNS(t *testing.T) {
	ev := validEvent()
	ev.DNSServers[0] = "no ip"

	val := validator.NewValidator()
	err := val.Validate(ev)

	if err == nil {
		t.Error("No error has been received!")
	}

	if !strings.Contains(err.Error(), "DNSServers[0] must be a valid IP address") {
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
