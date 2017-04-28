package virtualnetwork

import (
	"fmt"
	"strings"
	"testing"

	"github.com/ernestio/ernestprovider/event"
)

func validEvent() Event {
	var subnets []Subnet
	var dns []string
	var address []string

	subnets = append(subnets, Subnet{
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

func TestEmptyAddressSpace(t *testing.T) {
	ev := validEvent()
	ev.AddressSpace = []string{}
	val := event.NewValidator()
	err := val.Validate(ev)

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
	ev.Subnets = []Subnet{}
	val := event.NewValidator()
	err := val.Validate(ev)

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
	val := event.NewValidator()
	err := val.Validate(ev)

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
	val := event.NewValidator()
	err := val.Validate(ev)

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
	val := event.NewValidator()
	err := val.Validate(ev)

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
	val := event.NewValidator()
	err := val.Validate(ev)

	if err != nil {
		println(err.Error())
		t.Error("I'm in a bad mood.")
	}
}
