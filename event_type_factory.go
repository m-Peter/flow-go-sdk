/*
 * Flow Go SDK
 *
 * Copyright Flow Foundation
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package flow

import (
	"fmt"
)

type eventTypeFactory struct {
	address      string
	contractName string
	eventName    string
}

func (f eventTypeFactory) WithAddressString(address string) eventTypeFactory {
	f.address = address
	return f
}

func (f eventTypeFactory) WithAddress(address Address) eventTypeFactory {
	f.address = address.Hex()
	return f
}

func (f eventTypeFactory) WithContractName(contract string) eventTypeFactory {
	f.contractName = contract
	return f
}

func (f eventTypeFactory) WithEventName(event string) eventTypeFactory {
	f.eventName = event
	return f
}

func (f eventTypeFactory) String() string {
	return fmt.Sprintf("A.%s.%s.%s", f.address, f.contractName, f.eventName)
}

// NewEventTypeFactory helper function for constructing event names
func NewEventTypeFactory() eventTypeFactory {
	return eventTypeFactory{}
}
