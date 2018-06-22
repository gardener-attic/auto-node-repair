package resources

// Copyright (c) Microsoft and contributors.  All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

import (
	"context"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"net/http"
)

// PolicyDefinitionsClient is the client for the PolicyDefinitions methods of the Resources service.
type PolicyDefinitionsClient struct {
	BaseClient
}

// NewPolicyDefinitionsClient creates an instance of the PolicyDefinitionsClient client.
func NewPolicyDefinitionsClient(subscriptionID string) PolicyDefinitionsClient {
	return NewPolicyDefinitionsClientWithBaseURI(DefaultBaseURI, subscriptionID)
}

// NewPolicyDefinitionsClientWithBaseURI creates an instance of the PolicyDefinitionsClient client.
func NewPolicyDefinitionsClientWithBaseURI(baseURI string, subscriptionID string) PolicyDefinitionsClient {
	return PolicyDefinitionsClient{NewWithBaseURI(baseURI, subscriptionID)}
}

// CreateOrUpdate create or update policy definition.
// Parameters:
// policyDefinitionName - the policy definition name.
// parameters - the policy definition properties
func (client PolicyDefinitionsClient) CreateOrUpdate(ctx context.Context, policyDefinitionName string, parameters PolicyDefinition) (result PolicyDefinition, err error) {
	req, err := client.CreateOrUpdatePreparer(ctx, policyDefinitionName, parameters)
	if err != nil {
		err = autorest.NewErrorWithError(err, "resources.PolicyDefinitionsClient", "CreateOrUpdate", nil, "Failure preparing request")
		return
	}

	resp, err := client.CreateOrUpdateSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "resources.PolicyDefinitionsClient", "CreateOrUpdate", resp, "Failure sending request")
		return
	}

	result, err = client.CreateOrUpdateResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "resources.PolicyDefinitionsClient", "CreateOrUpdate", resp, "Failure responding to request")
	}

	return
}

// CreateOrUpdatePreparer prepares the CreateOrUpdate request.
func (client PolicyDefinitionsClient) CreateOrUpdatePreparer(ctx context.Context, policyDefinitionName string, parameters PolicyDefinition) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"policyDefinitionName": autorest.Encode("path", policyDefinitionName),
		"subscriptionId":       autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2015-11-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPut(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/providers/Microsoft.Authorization/policydefinitions/{policyDefinitionName}", pathParameters),
		autorest.WithJSON(parameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// CreateOrUpdateSender sends the CreateOrUpdate request. The method will close the
// http.Response Body if it receives an error.
func (client PolicyDefinitionsClient) CreateOrUpdateSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req,
		azure.DoRetryWithRegistration(client.Client))
}

// CreateOrUpdateResponder handles the response to the CreateOrUpdate request. The method always
// closes the http.Response Body.
func (client PolicyDefinitionsClient) CreateOrUpdateResponder(resp *http.Response) (result PolicyDefinition, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusCreated),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// Delete deletes policy definition.
// Parameters:
// policyDefinitionName - the policy definition name.
func (client PolicyDefinitionsClient) Delete(ctx context.Context, policyDefinitionName string) (result autorest.Response, err error) {
	req, err := client.DeletePreparer(ctx, policyDefinitionName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "resources.PolicyDefinitionsClient", "Delete", nil, "Failure preparing request")
		return
	}

	resp, err := client.DeleteSender(req)
	if err != nil {
		result.Response = resp
		err = autorest.NewErrorWithError(err, "resources.PolicyDefinitionsClient", "Delete", resp, "Failure sending request")
		return
	}

	result, err = client.DeleteResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "resources.PolicyDefinitionsClient", "Delete", resp, "Failure responding to request")
	}

	return
}

// DeletePreparer prepares the Delete request.
func (client PolicyDefinitionsClient) DeletePreparer(ctx context.Context, policyDefinitionName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"policyDefinitionName": autorest.Encode("path", policyDefinitionName),
		"subscriptionId":       autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2015-11-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsDelete(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/providers/Microsoft.Authorization/policydefinitions/{policyDefinitionName}", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// DeleteSender sends the Delete request. The method will close the
// http.Response Body if it receives an error.
func (client PolicyDefinitionsClient) DeleteSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req,
		azure.DoRetryWithRegistration(client.Client))
}

// DeleteResponder handles the response to the Delete request. The method always
// closes the http.Response Body.
func (client PolicyDefinitionsClient) DeleteResponder(resp *http.Response) (result autorest.Response, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByClosing())
	result.Response = resp
	return
}

// Get gets policy definition.
// Parameters:
// policyDefinitionName - the policy definition name.
func (client PolicyDefinitionsClient) Get(ctx context.Context, policyDefinitionName string) (result PolicyDefinition, err error) {
	req, err := client.GetPreparer(ctx, policyDefinitionName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "resources.PolicyDefinitionsClient", "Get", nil, "Failure preparing request")
		return
	}

	resp, err := client.GetSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "resources.PolicyDefinitionsClient", "Get", resp, "Failure sending request")
		return
	}

	result, err = client.GetResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "resources.PolicyDefinitionsClient", "Get", resp, "Failure responding to request")
	}

	return
}

// GetPreparer prepares the Get request.
func (client PolicyDefinitionsClient) GetPreparer(ctx context.Context, policyDefinitionName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"policyDefinitionName": autorest.Encode("path", policyDefinitionName),
		"subscriptionId":       autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2015-11-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/providers/Microsoft.Authorization/policydefinitions/{policyDefinitionName}", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// GetSender sends the Get request. The method will close the
// http.Response Body if it receives an error.
func (client PolicyDefinitionsClient) GetSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req,
		azure.DoRetryWithRegistration(client.Client))
}

// GetResponder handles the response to the Get request. The method always
// closes the http.Response Body.
func (client PolicyDefinitionsClient) GetResponder(resp *http.Response) (result PolicyDefinition, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}
