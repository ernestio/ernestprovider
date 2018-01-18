/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package event

// Collection : collection of events
type Collection struct {
	*Base
	Service        string  `json:"service" diff:"-"`
	ClientID       string  `json:"azure_client_id" diff:"-"`
	ClientSecret   string  `json:"azure_client_secret" diff:"-"`
	TenantID       string  `json:"azure_tenant_id" diff:"-"`
	SubscriptionID string  `json:"azure_subscription_id" diff:"-"`
	Environment    string  `json:"environment" diff:"-"`
	ResourceGroup  string  `json:"resource_group" diff:"-"`
	Resources      []Event `json:"components" diff:"-"`
}
