/*
Copyright 2017 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http:  www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

This file was copied and modified from the kubernetes autoscaler project
https://github.com/kubernetes/autoscaler/blob/cluster-autorepair-1.0.0/cluster-autoscaler/util/errors/error.go
*/

package errors

import (
	"fmt"
)

// AutorepairErrorType describes a high-level category of a given error
type AutorepairErrorType string

// AutorepairrError contains information about Autorepair errors
type AutorepairError interface {
	// Error implements golang error interface
	Error() string

	// Type returns the typ of AutorepairError
	Type() AutorepairErrorType

	// AddPrefix adds a prefix to error message.
	// Returns the error it's called for convienient inline use.
	// Example:
	// if err := DoSomething(myObject); err != nil {
	//	return err.AddPrefix("can't do something with %v: ", myObject)
	// }
	AddPrefix(msg string, args ...interface{}) AutorepairError
}

type autorepairErrorImpl struct {
	errorType AutorepairErrorType
	msg       string
}

const (
	// CloudProviderError is an error related to underlying infrastructure
	CloudProviderError AutorepairErrorType = "cloudProviderError"
	// ApiCallError is an error related to communication with k8s API server
	ApiCallError AutorepairErrorType = "apiCallError"
	// InternalError is an error inside Cluster autorepair
	InternalError AutorepairErrorType = "internalError"
	// TransientError is an error that causes us to skip a single loop, but
	// does not require any additional action.
	TransientError AutorepairErrorType = "transientError"
)

// NewAutorepairError returns new autorepair error with a message constructed from format string
func NewAutorepairError(errorType AutorepairErrorType, msg string, args ...interface{}) AutorepairError {
	return autorepairErrorImpl{
		errorType: errorType,
		msg:       fmt.Sprintf(msg, args...),
	}
}

// ToAutorepairError converts an error to AutorepairError with given type,
// unless it already is an AutorepairError (in which case it's not modified).
func ToAutorepairError(defaultType AutorepairErrorType, err error) AutorepairError {
	if e, ok := err.(AutorepairError); ok {
		return e
	}
	return NewAutorepairError(defaultType, err.Error())
}

// Error implements golang error interface
func (e autorepairErrorImpl) Error() string {
	return e.msg
}

// Type returns the typ of AutorepairError
func (e autorepairErrorImpl) Type() AutorepairErrorType {
	return e.errorType
}

// AddPrefix adds a prefix to error message.
// Returns the error it's called for convienient inline use.
// Example:
// if err := DoSomething(myObject); err != nil {
//	return err.AddPrefix("can't do something with %v: ", myObject)
// }
func (e autorepairErrorImpl) AddPrefix(msg string, args ...interface{}) AutorepairError {
	e.msg = fmt.Sprintf(msg, args...) + e.msg
	return e
}
