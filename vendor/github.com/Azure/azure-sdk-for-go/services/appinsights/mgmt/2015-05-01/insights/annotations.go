package insights

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

// AnnotationsClient is the composite Swagger for Application Insights Management Client
type AnnotationsClient struct {
	BaseClient
}

// NewAnnotationsClient creates an instance of the AnnotationsClient client.
func NewAnnotationsClient(subscriptionID string) AnnotationsClient {
	return NewAnnotationsClientWithBaseURI(DefaultBaseURI, subscriptionID)
}

// NewAnnotationsClientWithBaseURI creates an instance of the AnnotationsClient client.
func NewAnnotationsClientWithBaseURI(baseURI string, subscriptionID string) AnnotationsClient {
	return AnnotationsClient{NewWithBaseURI(baseURI, subscriptionID)}
}

// Create create an Annotation of an Application Insights component.
// Parameters:
// resourceGroupName - the name of the resource group.
// resourceName - the name of the Application Insights component resource.
// annotationProperties - properties that need to be specified to create an annotation of a Application
// Insights component.
func (client AnnotationsClient) Create(ctx context.Context, resourceGroupName string, resourceName string, annotationProperties Annotation) (result ListAnnotation, err error) {
	req, err := client.CreatePreparer(ctx, resourceGroupName, resourceName, annotationProperties)
	if err != nil {
		err = autorest.NewErrorWithError(err, "insights.AnnotationsClient", "Create", nil, "Failure preparing request")
		return
	}

	resp, err := client.CreateSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "insights.AnnotationsClient", "Create", resp, "Failure sending request")
		return
	}

	result, err = client.CreateResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "insights.AnnotationsClient", "Create", resp, "Failure responding to request")
	}

	return
}

// CreatePreparer prepares the Create request.
func (client AnnotationsClient) CreatePreparer(ctx context.Context, resourceGroupName string, resourceName string, annotationProperties Annotation) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"resourceName":      autorest.Encode("path", resourceName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2015-05-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPut(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Insights/components/{resourceName}/Annotations", pathParameters),
		autorest.WithJSON(annotationProperties),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// CreateSender sends the Create request. The method will close the
// http.Response Body if it receives an error.
func (client AnnotationsClient) CreateSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req,
		azure.DoRetryWithRegistration(client.Client))
}

// CreateResponder handles the response to the Create request. The method always
// closes the http.Response Body.
func (client AnnotationsClient) CreateResponder(resp *http.Response) (result ListAnnotation, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Value),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// Delete delete an Annotation of an Application Insights component.
// Parameters:
// resourceGroupName - the name of the resource group.
// resourceName - the name of the Application Insights component resource.
// annotationID - the unique annotation ID. This is unique within a Application Insights component.
func (client AnnotationsClient) Delete(ctx context.Context, resourceGroupName string, resourceName string, annotationID string) (result SetObject, err error) {
	req, err := client.DeletePreparer(ctx, resourceGroupName, resourceName, annotationID)
	if err != nil {
		err = autorest.NewErrorWithError(err, "insights.AnnotationsClient", "Delete", nil, "Failure preparing request")
		return
	}

	resp, err := client.DeleteSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "insights.AnnotationsClient", "Delete", resp, "Failure sending request")
		return
	}

	result, err = client.DeleteResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "insights.AnnotationsClient", "Delete", resp, "Failure responding to request")
	}

	return
}

// DeletePreparer prepares the Delete request.
func (client AnnotationsClient) DeletePreparer(ctx context.Context, resourceGroupName string, resourceName string, annotationID string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"annotationId":      autorest.Encode("path", annotationID),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"resourceName":      autorest.Encode("path", resourceName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2015-05-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsDelete(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Insights/components/{resourceName}/Annotations/{annotationId}", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// DeleteSender sends the Delete request. The method will close the
// http.Response Body if it receives an error.
func (client AnnotationsClient) DeleteSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req,
		azure.DoRetryWithRegistration(client.Client))
}

// DeleteResponder handles the response to the Delete request. The method always
// closes the http.Response Body.
func (client AnnotationsClient) DeleteResponder(resp *http.Response) (result SetObject, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Value),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// Get get the annotation for given id.
// Parameters:
// resourceGroupName - the name of the resource group.
// resourceName - the name of the Application Insights component resource.
// annotationID - the unique annotation ID. This is unique within a Application Insights component.
func (client AnnotationsClient) Get(ctx context.Context, resourceGroupName string, resourceName string, annotationID string) (result ListAnnotation, err error) {
	req, err := client.GetPreparer(ctx, resourceGroupName, resourceName, annotationID)
	if err != nil {
		err = autorest.NewErrorWithError(err, "insights.AnnotationsClient", "Get", nil, "Failure preparing request")
		return
	}

	resp, err := client.GetSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "insights.AnnotationsClient", "Get", resp, "Failure sending request")
		return
	}

	result, err = client.GetResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "insights.AnnotationsClient", "Get", resp, "Failure responding to request")
	}

	return
}

// GetPreparer prepares the Get request.
func (client AnnotationsClient) GetPreparer(ctx context.Context, resourceGroupName string, resourceName string, annotationID string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"annotationId":      autorest.Encode("path", annotationID),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"resourceName":      autorest.Encode("path", resourceName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2015-05-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Insights/components/{resourceName}/Annotations/{annotationId}", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// GetSender sends the Get request. The method will close the
// http.Response Body if it receives an error.
func (client AnnotationsClient) GetSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req,
		azure.DoRetryWithRegistration(client.Client))
}

// GetResponder handles the response to the Get request. The method always
// closes the http.Response Body.
func (client AnnotationsClient) GetResponder(resp *http.Response) (result ListAnnotation, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Value),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// List gets the list of annotations for a component for given time range
// Parameters:
// resourceGroupName - the name of the resource group.
// resourceName - the name of the Application Insights component resource.
// start - the start time to query from for annotations, cannot be older than 90 days from current date.
// end - the end time to query for annotations.
func (client AnnotationsClient) List(ctx context.Context, resourceGroupName string, resourceName string, start string, end string) (result ListAnnotation, err error) {
	req, err := client.ListPreparer(ctx, resourceGroupName, resourceName, start, end)
	if err != nil {
		err = autorest.NewErrorWithError(err, "insights.AnnotationsClient", "List", nil, "Failure preparing request")
		return
	}

	resp, err := client.ListSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "insights.AnnotationsClient", "List", resp, "Failure sending request")
		return
	}

	result, err = client.ListResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "insights.AnnotationsClient", "List", resp, "Failure responding to request")
	}

	return
}

// ListPreparer prepares the List request.
func (client AnnotationsClient) ListPreparer(ctx context.Context, resourceGroupName string, resourceName string, start string, end string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"resourceName":      autorest.Encode("path", resourceName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2015-05-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
		"end":         autorest.Encode("query", end),
		"start":       autorest.Encode("query", start),
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Insights/components/{resourceName}/Annotations", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// ListSender sends the List request. The method will close the
// http.Response Body if it receives an error.
func (client AnnotationsClient) ListSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req,
		azure.DoRetryWithRegistration(client.Client))
}

// ListResponder handles the response to the List request. The method always
// closes the http.Response Body.
func (client AnnotationsClient) ListResponder(resp *http.Response) (result ListAnnotation, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Value),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}
