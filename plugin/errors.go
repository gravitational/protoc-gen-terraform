/*
Copyright 2015-2021 Gravitational, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package plugin

import "fmt"

// invalidFieldError is generated when there is something wrong with struct field
type invalidFieldError struct {
	msg    string
	field  string
	reason string
}

// newInvalidFieldError creates unknown type error
func newInvalidFieldError(b *fieldBuilder, reason string) *invalidFieldError {
	return &invalidFieldError{msg: b.descriptor.GetName(), field: b.fieldDescriptor.GetName(), reason: reason}
}

// Error returns error message
func (e *invalidFieldError) Error() string {
	return fmt.Sprintf("%v (%v.%v)", e.reason, e.msg, e.field)
}

// invalidMessageError is generated when message is invalid (oneOf)
type invalidMessageError struct {
	name   string
	reason string
}

// newInvalidMessageError creates unknown type error
func newInvalidMessageError(name string, reason string) *invalidMessageError {
	return &invalidMessageError{name: name, reason: reason}
}

// Error returns error message
func (e *invalidMessageError) Error() string {
	return fmt.Sprintf("%v (%v)", e.reason, e.name)
}
