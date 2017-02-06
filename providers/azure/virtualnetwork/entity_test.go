package virtualnetwork

import (
	"fmt"
	"strings"
	"testing"

	"github.com/ernestio/ernestprovider/event"
)

func validEvent() Event {
	var subnets []subnet
	var dns []string
	var address []string

	subnets = append(subnets, subnet{
		Name:          "subnet",
		AddressPrefix: "10.2.0.1/24",
	})

	address = append(address, "10.2.0.1/24")
	dns = append(dns, "10.2.0.1/24")

	return Event{
		Name:           "supu",
		AddressSpace:   address,
		DNSServerNames: dns,
		Subnets:        subnets,
		Validator:      event.NewValidator(),
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

func TestEmptyAddressSpace(t *testing.T) {
	ev := validEvent()
	ev.AddressSpace = []string{}
	err := ev.Validate()

	if err == nil {
		t.Error("No error has been received!")
	}

	if !strings.Contains(err.Error(), "AddressSpace must contain at least 1 item") {
		fmt.Println(err.Error())
		t.Error("Output message does not contain name or required strings")
	}
}

func TestEmptySubnets(t *testing.T) {
	ev := validEvent()
	ev.Subnets = []subnet{}
	err := ev.Validate()

	if err == nil {
		t.Error("No error has been received!")
	}

	if !strings.Contains(err.Error(), "Subnet") {
		fmt.Println(err.Error())
		t.Error("Output message does not contain name or required strings")
	}
}

func TestSubnetsEmptyName(t *testing.T) {
	ev := validEvent()
	ev.Subnets[0].Name = ""
	err := ev.Validate()

	if err == nil {
		t.Error("No error has been received!")
	}

	if !strings.Contains(err.Error(), "Name is a required field") {
		fmt.Println(err.Error())
		t.Error("Output message does not contain name or required strings")
	}
}

func TestSubnetsInvalidPrefix(t *testing.T) {
	ev := validEvent()
	ev.Subnets[0].AddressPrefix = "supu"
	err := ev.Validate()

	if err == nil {
		t.Error("No error has been received!")
	}

	if !strings.Contains(err.Error(), "AddressPrefix must contain a valid CIDR notation") {
		fmt.Println(err.Error())
		t.Error("Output message does not contain name or required strings")
	}
}

func TestInvalidDNSServers(t *testing.T) {
	ev := validEvent()
	ev.DNSServerNames[0] = "supu"
	err := ev.Validate()

	if err == nil {
		t.Error("No error has been received!")
	}

	if !strings.Contains(err.Error(), "DNSServerNames[0] must be a valid IP address") {
		fmt.Println(err.Error())
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
