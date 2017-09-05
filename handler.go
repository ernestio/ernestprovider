/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package ernestprovider

import (
	"errors"
	"log"
	"strings"

	"github.com/ernestio/ernestprovider/event"
	"github.com/ernestio/ernestprovider/providers/azure/availabilityset"
	"github.com/ernestio/ernestprovider/providers/azure/lb"
	"github.com/ernestio/ernestprovider/providers/azure/lbbackendaddresspool"
	"github.com/ernestio/ernestprovider/providers/azure/lbprobe"
	"github.com/ernestio/ernestprovider/providers/azure/lbrule"
	"github.com/ernestio/ernestprovider/providers/azure/localnetworkgateway"
	"github.com/ernestio/ernestprovider/providers/azure/manageddisk"
	"github.com/ernestio/ernestprovider/providers/azure/networkinterface"
	"github.com/ernestio/ernestprovider/providers/azure/publicip"
	"github.com/ernestio/ernestprovider/providers/azure/resourcegroup"
	"github.com/ernestio/ernestprovider/providers/azure/securitygroup"
	"github.com/ernestio/ernestprovider/providers/azure/sqldatabase"
	"github.com/ernestio/ernestprovider/providers/azure/sqlfirewallrule"
	"github.com/ernestio/ernestprovider/providers/azure/sqlserver"
	"github.com/ernestio/ernestprovider/providers/azure/storageaccount"
	"github.com/ernestio/ernestprovider/providers/azure/storagecontainer"
	"github.com/ernestio/ernestprovider/providers/azure/subnet"
	"github.com/ernestio/ernestprovider/providers/azure/virtualmachine"
	"github.com/ernestio/ernestprovider/providers/azure/virtualnetwork"
)

// Handle : Handles the given event
func Handle(ev event.Event) (string, []byte) {
	var err error

	if err := ev.Process(); err != nil {
		ev.Log("error", err.Error())
		return ev.GetSubject() + ".error", ev.GetErroredBody()
	}

	parts := strings.Split(ev.GetSubject(), ".")
	switch parts[1] {
	case "create":
		if err := ev.Validate(); err != nil {
			ev.Log("error", err.Error())
			return ev.GetSubject() + ".error", ev.GetErroredBody()
		}
		err = ev.Create()
	case "update":
		if err := ev.Validate(); err != nil {
			ev.Log("error", err.Error())
			return ev.GetSubject() + ".error", ev.GetErroredBody()
		}
		err = ev.Update()
	case "delete":
		if err := ev.Validate(); err != nil {
			ev.Log("error", err.Error())
			return ev.GetSubject() + ".error", ev.GetErroredBody()
		}
		err = ev.Delete()
	case "get":
		err = ev.Get()
	case "find":
		err = ev.Find()
	case "validate":
		if err := ev.Validate(); err != nil {
			ev.Log("error", err.Error())
			return ev.GetSubject() + ".error", ev.GetErroredBody()
		}
	}

	if err != nil {
		ev.Error(err)
		return ev.GetSubject() + ".error", ev.GetErroredBody()
	}

	ev.Log("debug", "Component successfully processed")
	body := ev.GetCompletedBody()
	ev.Log("debug", string(body))
	return ev.GetSubject() + ".done", body
}

// GetAndHandle : Gets an event and Handles its results
func GetAndHandle(subject string, data []byte, key string) (string, []byte) {
	ev, err := GetEvent(subject, data, key)
	if err != nil {
		log.Println("[ERROR] : Event not found (A) - " + err.Error())
		return subject + ".error", data
	}
	if ev == nil {
		log.Println("[ERROR] : Event not found (" + subject + ") ")
		return subject + ".error", data
	}

	return Handle(ev)
}

// GetEvent : Gets a valid event based on a subject
func GetEvent(subject string, data []byte, key string) (event.Event, error) {
	parts := strings.Split(subject, ".")
	switch parts[2] {
	case "aws":
		return getAWSEvent(subject, data, key)
	case "azure":
		return getAzureEvent(subject, data, key)
	}
	return nil, errors.New("Unkown provider")
}

func getAzureEvent(subject string, data []byte, key string) (event.Event, error) {
	var ev event.Event
	var err error
	parts := strings.Split(subject, ".")
	val := event.NewValidator()
	switch parts[0] {
	case "public_ip", "public_ips":
		ev, err = publicip.New(subject, key, data, val)
	case "virtual_network", "virtual_networks":
		ev, err = virtualnetwork.New(subject, key, data, val)
	case "resource_group", "resource_groups":
		ev, err = resourcegroup.New(subject, key, data, val)
	case "subnet", "subnets":
		ev, err = subnet.New(subject, key, data, val)
	case "network_interface", "network_interfaces":
		ev, err = networkinterface.New(subject, key, data, val)
	case "managed_disk", "managed_disks":
		ev, err = manageddisk.New(subject, key, data, val)
	case "storage_account", "storage_accounts":
		ev, err = storageaccount.New(subject, key, data, val)
	case "storage_container", "storage_containers":
		ev, err = storagecontainer.New(subject, key, data, val)
	case "virtual_machine", "virtual_machines":
		ev, err = virtualmachine.New(subject, key, data, val)
	case "availability_set", "availability_sets":
		ev, err = availabilityset.New(subject, key, data, val)
	case "lb", "lbs":
		ev, err = lb.New(subject, key, data, val)
	case "lb_rule", "lb_rules":
		ev, err = lbrule.New(subject, key, data, val)
	case "lb_probe", "lb_probes":
		ev, err = lbprobe.New(subject, key, data, val)
	case "lb_backend_address_pool", "lb_backend_address_pools":
		ev, err = lbbackendaddresspool.New(subject, key, data, val)
	case "sql_server", "sql_servers":
		ev, err = sqlserver.New(subject, key, data, val)
	case "local_network_gateway", "local_network_gateways":
		ev, err = localnetworkgateway.New(subject, key, data, val)
	case "security_group", "security_groups":
		ev, err = securitygroup.New(subject, key, data, val)
	case "sql_database", "sql_databases":
		ev, err = sqldatabase.New(subject, key, data, val)
	case "sql_firewall_rule", "sql_firewall_rules":
		ev, err = sqlfirewallrule.New(subject, key, data, val)
	}
	return ev, err
}

func getAWSEvent(subject string, data []byte, key string) (event.Event, error) {
	return nil, errors.New("Unconfigured provider")
}
