/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package networkinterface

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/r3labs/terraform/helper/hashcode"
	"github.com/r3labs/terraform/helper/schema"

	aes "github.com/ernestio/crypto/aes"
	"github.com/ernestio/ernestprovider/event"
	"github.com/ernestio/ernestprovider/providers/azure"
	types "github.com/ernestio/ernestprovider/types/azure/networkinterface"
	"github.com/ernestio/ernestprovider/validator"
	"github.com/r3labs/terraform/builtin/providers/azurerm"
)

// Event : This is the Ernest representation of an azure networkinterface
type Event struct {
	types.Event
	ErrorMessage string               `json:"error,omitempty" diff:"-"`
	CryptoKey    string               `json:"-" diff:"-"`
	Validator    *validator.Validator `json:"-" diff:"-"`
	GenericEvent event.Event          `json:"-" validate:"-" diff:"-"`
}

// New : Constructor
func New(subject, cryptoKey string, body []byte, val *validator.Validator) (event.Event, error) {
	// var ev event.Resource
	ev := &Event{CryptoKey: cryptoKey, Validator: val}
	body = []byte(strings.Replace(string(body), `"_component":"network_interfaces"`, `"_component":"network_interface"`, 1))
	if err := json.Unmarshal(body, &ev); err != nil {
		err := fmt.Errorf("Error on input message : %s", err)
		return nil, err
	}

	ev.GenericEvent, _ = azure.New(subject, "azurerm_network_interface", body, val, ev)
	return ev.GenericEvent, nil
}

// SetComponents : ....
func (ev *Event) SetComponents(components []event.Event) {
	for _, v := range components {
		ev.Components = append(ev.Components, v.GetBody())
	}
}

// ValidateID : determines if the given id is valid for this resource type
func (ev *Event) ValidateID(id string) bool {
	parts := strings.Split(strings.ToLower(id), "/")
	if len(parts) != 9 {
		return false
	}
	if parts[6] != "microsoft.network" {
		return false
	}
	if parts[7] != "networkinterfaces" {
		return false
	}
	return true
}

// SetID : id setter
func (ev *Event) SetID(id string) {
	ev.ID = id
}

// GetID : id getter
func (ev *Event) GetID() string {
	return ev.ID
}

// SetState : state setter
func (ev *Event) SetState(state string) {
	ev.State = state
}

// ResourceDataToEvent : Translates a ResourceData on a valid Ernest Event
func (ev *Event) ResourceDataToEvent(d *schema.ResourceData) error {

	ev.ID = d.Id()
	if ev.ID == "" {
		ev.Name = d.Get("name").(string)
	} else {
		parts := strings.Split(ev.ID, "/")
		ev.Name = parts[8]
	}
	ev.ComponentID = "network_interface::" + ev.Name
	ev.ResourceGroupName = d.Get("resource_group_name").(string)
	ev.Location = d.Get("location").(string)
	ev.NetworkSecurityGroupID = d.Content["network_security_group_id"].(string)
	ev.MacAddress = d.Get("mac_address").(string)
	ev.PrivateIPAddress = d.Get("private_ip_address").(string)
	ev.VirtualMachineID = d.Get("virtual_machine_id").(string)

	configs := []types.IPConfiguration{}

	cli, _ := ev.GenericEvent.Client()
	list := cli.ListNetworkInterfaceConfigurations(ev.ResourceGroupName, ev.Name)

	for _, mo := range list {
		if mo["interface"] == ev.Name {
			c := types.IPConfiguration{
				Name:                       mo["name"],
				SubnetID:                   mo["subnet_id"],
				PrivateIPAddress:           mo["private_ip_address"],
				PrivateIPAddressAllocation: mo["private_ip_address_allocation"],
				PublicIPAddressID:          mo["public_ip_address_id"],
			}
			c.LoadBalancerBackendAddressPoolIDs = d.Content["lb_pools"].([]string)
			for _, v := range c.LoadBalancerBackendAddressPoolIDs {
				parts := strings.Split(v, "/")
				c.LoadBalancerBackendAddressPools = append(c.LoadBalancerBackendAddressPools, parts[len(parts)-1])
			}
			if mo["load_balancer_inbound_nat_rules_ids"] != "" {
				c.LoadBalancerInboundNatRules = strings.Split(mo["load_balancer_inbound_nat_rules_ids"], ",")
			}
			configs = append(configs, c)
		}
	}
	ev.IPConfigurations = configs
	ev.DNSServers = make([]string, 0)
	for _, v := range d.Get("dns_servers").(*schema.Set).List() {
		ev.DNSServers = append(ev.DNSServers, v.(string))
	}

	ev.InternalDNSNameLabel = d.Get("internal_dns_name_label").(string)
	ev.AppliedDNSServers = make([]string, 0)
	for _, v := range d.Get("applied_dns_servers").(*schema.Set).List() {
		ev.AppliedDNSServers = append(ev.AppliedDNSServers, v.(string))
	}

	ev.InternalFQDN = d.Get("internal_fqdn").(string)
	ev.EnableIPForwarding = *(d.Content["enable_ip_forwarding"].(*bool))

	tags := d.Get("tags").(map[string]interface{})
	ev.Tags = make(map[string]string, 0)
	for k, v := range tags {
		ev.Tags[k] = v.(string)
	}

	return nil
}

// EventToResourceData : Translates the current event on a valid ResourceData
func (ev *Event) EventToResourceData(d *schema.ResourceData) error {
	crypto := aes.New()

	encFields := make(map[string]string)
	encFields["subscription_id"] = ev.SubscriptionID
	encFields["client_id"] = ev.ClientID
	encFields["client_secret"] = ev.ClientSecret
	encFields["tenant_id"] = ev.TenantID
	encFields["environment"] = ev.Environment
	for k, v := range encFields {
		dec, err := crypto.Decrypt(v, ev.CryptoKey)
		if err != nil {
			err := fmt.Errorf("Field '%s' not valid : %s", k, err)
			ev.Log("error", err.Error())
			return err
		}
		if err := d.Set(k, dec); err != nil {
			err := fmt.Errorf("Field '%s' not valid : %s", k, err)
			ev.Log("error", err.Error())
			return err
		}
	}

	fields := make(map[string]interface{})
	fields["name"] = ev.Name
	fields["resource_group_name"] = ev.ResourceGroupName
	fields["location"] = ev.Location
	fields["network_security_group_id"] = ev.NetworkSecurityGroupID
	fields["mac_address"] = ev.MacAddress
	fields["private_ip_address"] = ev.PrivateIPAddress
	fields["virtual_machine_id"] = ev.VirtualMachineID
	fields["ip_configuration"] = ev.mapIPConfigurations()
	fields["dns_servers"] = ev.DNSServers
	fields["internal_dns_name_label"] = ev.InternalDNSNameLabel
	fields["applied_dns_servers"] = ev.AppliedDNSServers
	fields["internal_fqdn"] = ev.InternalFQDN
	fields["enable_ip_forwarding"] = ev.EnableIPForwarding
	fields["tags"] = ev.Tags
	for k, v := range fields {
		if err := d.Set(k, v); err != nil {
			err := fmt.Errorf("Field '%s' not valid : %s", k, err)
			ev.Log("error", err.Error())
			return err
		}
	}

	return nil
}

func (ev *Event) mapIPConfigurations() *schema.Set {
	list := &schema.Set{
		F: resourceArmNetworkInterfaceIPConfigurationHash,
	}
	for _, c := range ev.IPConfigurations {
		conf := map[string]interface{}{}
		conf["name"] = c.Name
		conf["subnet_id"] = c.SubnetID
		conf["private_ip_address"] = c.PrivateIPAddress
		conf["private_ip_address_allocation"] = c.PrivateIPAddressAllocation
		conf["public_ip_address_id"] = c.PublicIPAddressID
		l1 := schema.Set{
			F: resourceHashArnString,
		}
		for _, v := range c.LoadBalancerBackendAddressPoolIDs {
			l1.Add(v)
		}
		conf["load_balancer_backend_address_pools_ids"] = &l1
		l2 := schema.Set{
			F: resourceHashArnString,
		}
		for _, v := range c.LoadBalancerInboundNatRules {
			l2.Add(v)
		}
		conf["load_balancer_inbound_nat_rules_ids"] = &l2
		list.Add(conf)
	}
	return list
}

func resourceArmNetworkInterfaceIPConfigurationHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	buf.WriteString(fmt.Sprintf("%s-", m["name"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["subnet_id"].(string)))
	if m["private_ip_address"] != nil {
		buf.WriteString(fmt.Sprintf("%s-", m["private_ip_address"].(string)))
	}
	buf.WriteString(fmt.Sprintf("%s-", m["private_ip_address_allocation"].(string)))
	if m["public_ip_address_id"] != nil {
		buf.WriteString(fmt.Sprintf("%s-", m["public_ip_address_id"].(string)))
	}
	if m["load_balancer_backend_address_pools_ids"] != nil {
		str := fmt.Sprintf("*Set(%s)", m["load_balancer_backend_address_pools_ids"].(*schema.Set))
		buf.WriteString(fmt.Sprintf("%s-", str))
	}
	if m["load_balancer_inbound_nat_rules_ids"] != nil {
		str := fmt.Sprintf("*Set(%s)", m["load_balancer_inbound_nat_rules_ids"].(*schema.Set))
		buf.WriteString(fmt.Sprintf("%s-", str))
	}

	return hashcode.String(buf.String())
}

func resourceHashArnString(v interface{}) int {
	var buf bytes.Buffer

	buf.WriteString(fmt.Sprintf("%d-", schema.HashString(v.(string))))

	return hashcode.String(buf.String())
}

// Clone : will mark the event as errored
func (ev *Event) Clone() (event.Event, error) {
	body, _ := json.Marshal(ev)
	return New(ev.Subject, ev.CryptoKey, body, ev.Validator)
}

// Error : will mark the event as errored
func (ev *Event) Error(err error) {
	ev.ErrorMessage = err.Error()
	ev.Body, err = json.Marshal(ev)
}

// Client : not implemented
func (ev *Event) Client() (*azurerm.ArmClient, error) {
	return nil, errors.New("Not implemented")
}
