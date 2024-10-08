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

import "time"

type ExecutionData struct {
	BlockID            Identifier
	ChunkExecutionData []*ChunkExecutionData
}

type ExecutionDataStreamResponse struct {
	Height         uint64
	ExecutionData  *ExecutionData
	BlockTimestamp time.Time
}

type ChunkExecutionData struct {
	Transactions       []*Transaction
	Events             []*Event
	TrieUpdate         *TrieUpdate
	TransactionResults []*LightTransactionResult
}

type TrieUpdate struct {
	RootHash []byte
	Paths    [][]byte
	Payloads []*Payload
}

type Payload struct {
	KeyPart []*KeyPart
	Value   []byte
}

type KeyPart struct {
	Type  uint16
	Value []byte
}

type LightTransactionResult struct {
	// TransactionID is the ID of the transaction this result was emitted from.
	TransactionID Identifier
	// Failed is true if the transaction's execution failed resulting in an error, false otherwise.
	Failed bool
	// ComputationUsed is amount of computation used while executing the transaction.
	ComputationUsed uint64
}
