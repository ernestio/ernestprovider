/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package event

// Collection : collection of events
type Collection struct {
	*Base
	Service        string  `json:"service"`
	ClientID       string  `json:"azure_client_id"`
	ClientSecret   string  `json:"azure_client_secret"`
	TenantID       string  `json:"azure_tenant_id"`
	SubscriptionID string  `json:"azure_subscription_id"`
	Environment    string  `json:"environment"`
	ResourceGroup  string  `json:"resource_group"`
	Resources      []Event `json:"components"`
}
