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

package http

import (
	"context"
	"encoding/base64"
	"fmt"
	"testing"

	jsoncdc "github.com/onflow/cadence/encoding/json"
	"github.com/stretchr/testify/require"

	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/access/http/convert"
	"github.com/onflow/flow-go-sdk/access/http/internal/unittest"
	"github.com/onflow/flow-go-sdk/access/http/models"
	"github.com/onflow/flow-go-sdk/test"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func clientTest(
	f func(ctx context.Context, t *testing.T, handler *mockHandler, client *Client),
) func(t *testing.T) {
	return func(t *testing.T) {
		h := &mockHandler{}
		client := &Client{
			&BaseClient{handler: h},
		}
		f(context.Background(), t, h, client)
		h.AssertExpectations(t)
	}
}

func TestClient_Factories(t *testing.T) {

	client, err := NewClient(EmulatorHost)
	assert.NoError(t, err)
	assert.NotNil(t, client)

	client, err = NewClient(EmulatorHost)
	assert.NoError(t, err)
	assert.NotNil(t, client)

	client, err = NewClient(EmulatorHost)
	assert.NoError(t, err)
	assert.NotNil(t, client)

	client, err = NewClient(EmulatorHost)
	assert.NoError(t, err)
	assert.NotNil(t, client)
}
func TestClient_ClientOptions(t *testing.T) {
	t.Run("WithJSONOptions", func(t *testing.T) {

		expectedJsonOption := []jsoncdc.Option{jsoncdc.WithBackwardsCompatibility()}

		// Confirm that the options are set
		client, err := NewClient(EmulatorHost, WithJSONOptions(expectedJsonOption...))
		assert.NoError(t, err)
		assert.NotNil(t, client)

		// hard to run a contains check on the options due to it comparing functions, so just check the length
		assert.Equal(t, len(client.httpClient.jsonOptions), len(expectedJsonOption)+len(DefaultClientOptions().jsonOptions))
	})
}

func TestBaseClient_GetNodeInfo(t *testing.T) {
	t.Run("Success", clientTest(func(ctx context.Context, t *testing.T, handler *mockHandler, client *Client) {
		httpInfo := &models.NodeVersionInfo{
			Semver:               "1.0",
			Commit:               "123",
			SporkId:              flow.HexToID("0x01").String(),
			ProtocolVersion:      "2",
			SporkRootBlockHeight: "123",
			NodeRootBlockHeight:  "321",
		}
		expectedInfo, err := convert.ToNodeVersionInfo(httpInfo)
		require.NoError(t, err)

		handler.
			On("getNodeVersionInfo", mock.Anything).
			Return(httpInfo, nil)

		info, err := client.GetNodeVersionInfo(ctx)
		require.NoError(t, err)
		assert.Equal(t, expectedInfo, info)
	}))
}

func TestBaseClient_GetNetworkParameters(t *testing.T) {
	const handlerName = "getNetworkParameters"

	t.Run("Success", clientTest(func(ctx context.Context, t *testing.T, handler *mockHandler, client *Client) {
		httpParams := unittest.NetworkParametersFlowFixture()
		expectedParams := convert.ToNetworkParameters(&httpParams)

		handler.
			On(handlerName, mock.Anything).
			Return(&httpParams, nil)

		params, err := client.GetNetworkParameters(ctx)
		assert.NoError(t, err)
		assert.Equal(t, params, expectedParams)
	}))

	t.Run("Failure", clientTest(func(ctx context.Context, t *testing.T, handler *mockHandler, client *Client) {
		handler.
			On(handlerName, mock.Anything).
			Return(nil, HTTPError{
				Url:     "/",
				Code:    400,
				Message: "bad request",
			})

		params, err := client.GetNetworkParameters(ctx)
		assert.EqualError(t, err, "bad request")
		assert.Nil(t, params)
	}))
}

func TestBaseClient_GetBlockByID(t *testing.T) {
	const handlerName = "getBlockByID"
	t.Run("Success", clientTest(func(ctx context.Context, t *testing.T, handler *mockHandler, client *Client) {
		httpBlock := unittest.BlockFlowFixture()
		expectedBlock, err := convert.ToBlock(&httpBlock)
		assert.NoError(t, err)

		handler.
			On(handlerName, mock.Anything, httpBlock.Header.Id).
			Return(&httpBlock, nil)

		block, err := client.GetBlockByID(ctx, flow.HexToID(httpBlock.Header.Id))
		assert.NoError(t, err)
		assert.Equal(t, block, expectedBlock)
	}))

	t.Run("Get Block Header", clientTest(func(ctx context.Context, t *testing.T, handler *mockHandler, client *Client) {
		httpBlock := unittest.BlockFlowFixture()
		expectedBlock, err := convert.ToBlock(&httpBlock)
		assert.NoError(t, err)

		handler.
			On(handlerName, mock.Anything, httpBlock.Header.Id).
			Return(&httpBlock, nil)

		header, err := client.GetBlockHeaderByID(ctx, flow.HexToID(httpBlock.Header.Id))
		assert.NoError(t, err)
		assert.Equal(t, header, &expectedBlock.BlockHeader)
	}))

	t.Run("Not found", clientTest(func(ctx context.Context, t *testing.T, handler *mockHandler, client *Client) {
		handler.
			On(handlerName, mock.Anything, mock.Anything).
			Return(nil, HTTPError{
				Url:     "/",
				Code:    404,
				Message: "block not found",
			})
		block, err := client.GetBlockByID(ctx, flow.HexToID("0x1"))
		assert.EqualError(t, err, "block not found")
		assert.Nil(t, block)
	}))
}

func TestBaseClient_GetBlockByHeight(t *testing.T) {
	const handlerName = "getBlocksByHeights"

	t.Run("Success", clientTest(func(ctx context.Context, t *testing.T, handler *mockHandler, client *Client) {
		httpBlock := unittest.BlockFlowFixture()
		expectedBlock, err := convert.ToBlock(&httpBlock)
		assert.NoError(t, err)

		handler.
			On(handlerName, mock.Anything, httpBlock.Header.Height, "", "").
			Return([]*models.Block{&httpBlock}, nil)

		block, err := client.GetBlockByHeight(ctx, expectedBlock.Height)
		assert.NoError(t, err)
		assert.Equal(t, block, expectedBlock)
	}))

	t.Run("Not found", clientTest(func(ctx context.Context, t *testing.T, handler *mockHandler, client *Client) {
		handler.
			On(handlerName, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(nil, HTTPError{
				Url:     "/",
				Code:    404,
				Message: "block not found",
			})

		block, err := client.GetBlockByHeight(ctx, 10)
		assert.EqualError(t, err, "block not found")
		assert.Nil(t, block)
	}))

	t.Run("Get Block Header", clientTest(func(ctx context.Context, t *testing.T, handler *mockHandler, client *Client) {
		httpBlock := unittest.BlockFlowFixture()
		expectedBlock, err := convert.ToBlock(&httpBlock)
		assert.NoError(t, err)

		handler.
			On(handlerName, mock.Anything, httpBlock.Header.Height, "", "").
			Return([]*models.Block{&httpBlock}, nil)

		block, err := client.GetBlockHeaderByHeight(ctx, expectedBlock.Height)
		assert.NoError(t, err)
		assert.Equal(t, block, &expectedBlock.BlockHeader)
	}))
}

func TestBaseClient_GetLatestBlock(t *testing.T) {
	const handlerName = "getBlocksByHeights"

	t.Run("Block Sealed", clientTest(func(ctx context.Context, t *testing.T, handler *mockHandler, client *Client) {
		httpBlock := unittest.BlockFlowFixture()
		expectedBlock, err := convert.ToBlock(&httpBlock)
		assert.NoError(t, err)

		handler.
			On(handlerName, mock.Anything, "sealed", "", "").
			Return([]*models.Block{&httpBlock}, nil)

		block, err := client.GetLatestBlock(ctx, true)
		assert.NoError(t, err)
		assert.Equal(t, block, expectedBlock)
	}))

	t.Run("Block Not Sealed", clientTest(func(ctx context.Context, t *testing.T, handler *mockHandler, client *Client) {
		httpBlock := unittest.BlockFlowFixture()
		expectedBlock, err := convert.ToBlock(&httpBlock)
		assert.NoError(t, err)

		handler.
			On(handlerName, mock.Anything, "final", "", "").
			Return([]*models.Block{&httpBlock}, nil)

		block, err := client.GetLatestBlock(ctx, false)
		assert.NoError(t, err)
		assert.Equal(t, block, expectedBlock)
	}))

	t.Run("Final Header", clientTest(func(ctx context.Context, t *testing.T, handler *mockHandler, client *Client) {
		httpBlock := unittest.BlockFlowFixture()
		expectedBlock, err := convert.ToBlock(&httpBlock)
		assert.NoError(t, err)

		handler.
			On(handlerName, mock.Anything, "final", "", "").
			Return([]*models.Block{&httpBlock}, nil)

		block, err := client.GetLatestBlockHeader(ctx, false)
		assert.NoError(t, err)
		assert.Equal(t, block, &expectedBlock.BlockHeader)
	}))

	t.Run("Sealed Header", clientTest(func(ctx context.Context, t *testing.T, handler *mockHandler, client *Client) {
		httpBlock := unittest.BlockFlowFixture()
		expectedBlock, err := convert.ToBlock(&httpBlock)
		assert.NoError(t, err)

		handler.
			On(handlerName, mock.Anything, "sealed", "", "").
			Return([]*models.Block{&httpBlock}, nil)

		block, err := client.GetLatestBlockHeader(ctx, true)
		assert.NoError(t, err)
		assert.Equal(t, block, &expectedBlock.BlockHeader)
	}))
}

func TestBaseClient_GetCollection(t *testing.T) {
	const handlerName = "getCollection"

	t.Run("Success", clientTest(func(ctx context.Context, t *testing.T, handler *mockHandler, client *Client) {
		httpCollection := unittest.CollectionFlowFixture()
		expectedCollection := convert.ToCollection(&httpCollection)

		handler.
			On(handlerName, mock.Anything, expectedCollection.ID().String()).
			Return(&httpCollection, nil)

		collection, err := client.GetCollection(ctx, expectedCollection.ID())

		assert.NoError(t, err)
		assert.Equal(t, collection, expectedCollection)
	}))

	t.Run("Not Found", clientTest(func(ctx context.Context, t *testing.T, handler *mockHandler, client *Client) {
		handler.
			On(handlerName, mock.Anything, mock.Anything).
			Return(nil, HTTPError{
				Url:     "/",
				Code:    404,
				Message: "collection not found",
			})

		coll, err := client.GetCollection(ctx, flow.HexToID("0x1"))
		assert.EqualError(t, err, "collection not found")
		assert.Nil(t, coll)
	}))
}

func TestBaseClient_SendTransaction(t *testing.T) {
	const handlerName = "sendTransaction"

	t.Run("Success", clientTest(func(ctx context.Context, t *testing.T, handler *mockHandler, client *Client) {
		httpTx := unittest.TransactionFlowFixture()
		expectedTx, err := convert.ToTransaction(&httpTx)
		assert.NoError(t, err)

		sentTx, err := convert.TncodeTransaction(*expectedTx)
		assert.NoError(t, err)

		handler.
			On(handlerName, mock.Anything, sentTx).
			Return(nil)

		err = client.SendTransaction(ctx, *expectedTx)
		assert.NoError(t, err)
	}))

	t.Run("Not Found", clientTest(func(ctx context.Context, t *testing.T, handler *mockHandler, client *Client) {
		handler.On(handlerName, mock.Anything, mock.Anything).Return(HTTPError{
			Url:     "/",
			Code:    400,
			Message: "invalid payload",
		})

		tx := test.TransactionGenerator().New()
		err := client.SendTransaction(ctx, *tx)
		assert.EqualError(t, err, "invalid payload")
	}))
}

func TestBaseClient_GetTransaction(t *testing.T) {
	const handlerName = "getTransaction"

	t.Run("Success", clientTest(func(ctx context.Context, t *testing.T, handler *mockHandler, client *Client) {
		httpTx := unittest.TransactionFlowFixture()
		expectedTx, err := convert.ToTransaction(&httpTx)
		assert.NoError(t, err)

		handler.
			On(handlerName, mock.Anything, expectedTx.ID().String(), false).
			Return(&httpTx, nil)

		tx, err := client.GetTransaction(ctx, expectedTx.ID())
		assert.NoError(t, err)
		assert.Equal(t, tx, expectedTx)
	}))

	t.Run("Not Found", clientTest(func(ctx context.Context, t *testing.T, handler *mockHandler, client *Client) {
		handler.On(handlerName, mock.Anything, mock.Anything, mock.Anything).Return(nil, HTTPError{
			Url:     "/",
			Code:    404,
			Message: "tx not found",
		})

		tx, err := client.GetTransaction(ctx, flow.HexToID("0x1"))
		assert.EqualError(t, err, "tx not found")
		assert.Nil(t, tx)
	}))
}

func TestBaseClient_GetTransactionResult(t *testing.T) {
	const handlerName = "getTransaction"

	t.Run("Success", clientTest(func(ctx context.Context, t *testing.T, handler *mockHandler, client *Client) {
		httpTx := unittest.TransactionFlowFixture()
		httpTxRes := unittest.TransactionResultFlowFixture(flow.EventEncodingVersionJSONCDC)
		httpTx.Result = &httpTxRes
		expectedTx, err := convert.ToTransaction(&httpTx)
		assert.NoError(t, err)

		expectedTxRes, err := convert.ToTransactionResult(&httpTxRes, nil)
		assert.NoError(t, err)

		handler.
			On(handlerName, mock.Anything, expectedTx.ID().String(), true).
			Return(&httpTx, nil)

		txRes, err := client.GetTransactionResult(ctx, expectedTx.ID())
		assert.NoError(t, err)
		assert.Equal(t, txRes, expectedTxRes)
	}))

	t.Run("Not Found", clientTest(func(ctx context.Context, t *testing.T, handler *mockHandler, client *Client) {
		handler.On(handlerName, mock.Anything, mock.Anything, true).Return(nil, HTTPError{
			Url:     "/",
			Code:    404,
			Message: "tx result not found",
		})

		tx, err := client.GetTransactionResult(ctx, flow.HexToID("0x1"))
		assert.EqualError(t, err, "tx result not found")
		assert.Nil(t, tx)
	}))
}

func TestBaseClient_GetAccount(t *testing.T) {
	const handlerName = "getAccount"

	t.Run("Success", clientTest(func(ctx context.Context, t *testing.T, handler *mockHandler, client *Client) {
		httpAccount := unittest.AccountFlowFixture()
		expectedAccount, err := convert.ToAccount(&httpAccount)
		assert.NoError(t, err)

		handler.
			On(handlerName, mock.Anything, httpAccount.Address, "sealed").
			Return(&httpAccount, nil)

		account, err := client.GetAccount(ctx, expectedAccount.Address)
		assert.NoError(t, err)
		assert.Equal(t, account, expectedAccount)

		account, err = client.GetAccountAtLatestBlock(ctx, expectedAccount.Address)
		assert.NoError(t, err)
		assert.Equal(t, account, expectedAccount)
	}))

	t.Run("Not Found", clientTest(func(ctx context.Context, t *testing.T, handler *mockHandler, client *Client) {
		handler.On(handlerName, mock.Anything, mock.Anything, mock.Anything).Return(nil, HTTPError{
			Url:     "/",
			Code:    404,
			Message: "account not found",
		})

		acc1, err := client.GetAccount(ctx, flow.HexToAddress("0x1"))
		assert.EqualError(t, err, "account not found")
		assert.Nil(t, acc1)

		acc2, err := client.GetAccountAtLatestBlock(ctx, flow.HexToAddress("0x1"))
		assert.EqualError(t, err, "account not found")
		assert.Nil(t, acc2)
	}))
}

func TestBaseClient_GetAccountAtBlockHeight(t *testing.T) {
	const handlerName = "getAccount"

	t.Run("Success", clientTest(func(ctx context.Context, t *testing.T, handler *mockHandler, client *Client) {
		httpAccount := unittest.AccountFlowFixture()
		expectedAccount, err := convert.ToAccount(&httpAccount)
		assert.NoError(t, err)

		handler.
			On(handlerName, mock.Anything, httpAccount.Address, "10").
			Return(&httpAccount, nil)

		account, err := client.GetAccountAtBlockHeight(ctx, expectedAccount.Address, 10)
		assert.NoError(t, err)
		assert.Equal(t, account, expectedAccount)
	}))

	t.Run("Not Found", clientTest(func(ctx context.Context, t *testing.T, handler *mockHandler, client *Client) {
		handler.On(handlerName, mock.Anything, mock.Anything, mock.Anything).Return(nil, HTTPError{
			Url:     "/",
			Code:    404,
			Message: "account not found",
		})

		acc, err := client.GetAccountAtBlockHeight(ctx, flow.HexToAddress("0x1"), 10)
		assert.EqualError(t, err, "account not found")
		assert.Nil(t, acc)
	}))
}

func TestBaseClient_ExecuteScript(t *testing.T) {

	t.Run("Success Block Height", clientTest(func(ctx context.Context, t *testing.T, handler *mockHandler, client *Client) {
		script := []byte(`main() { return "Hello World" }`)
		encodedScript := base64.StdEncoding.EncodeToString(script)
		const height uint64 = 10
		response := base64.StdEncoding.EncodeToString([]byte(`{
		  "type": "String",
		  "value": "Hello World"
		}`))

		handler.
			On("executeScriptAtBlockHeight", mock.Anything, fmt.Sprintf("%d", height), encodedScript, []string{}).
			Return(response, nil)

		val, err := client.ExecuteScriptAtBlockHeight(ctx, height, script, nil)
		assert.NoError(t, err)
		assert.Equal(t, val.String(), "\"Hello World\"")
	}))

	t.Run("Success Latest Height", clientTest(func(ctx context.Context, t *testing.T, handler *mockHandler, client *Client) {
		script := []byte(`main() { return "Hello World" }`)
		encodedScript := base64.StdEncoding.EncodeToString(script)
		response := base64.StdEncoding.EncodeToString([]byte(`{
		  "type": "String",
		  "value": "Hello World"
		}`))

		handler.
			On("executeScriptAtBlockHeight", mock.Anything, "sealed", encodedScript, []string{}).
			Return(response, nil)

		val, err := client.ExecuteScriptAtLatestBlock(ctx, script, nil)
		assert.NoError(t, err)
		assert.Equal(t, val.String(), "\"Hello World\"")
	}))

	t.Run("Success Block ID", clientTest(func(ctx context.Context, t *testing.T, handler *mockHandler, client *Client) {
		script := []byte(`main() { return "Hello World" }`)
		encodedScript := base64.StdEncoding.EncodeToString(script)
		id := flow.HexToID("0x1")
		response := base64.StdEncoding.EncodeToString([]byte(`{
		  "type": "String",
		  "value": "Hello World"
		}`))

		handler.
			On("executeScriptAtBlockID", mock.Anything, id.String(), encodedScript, []string{}).
			Return(response, nil)

		val, err := client.ExecuteScriptAtBlockID(ctx, id, script, nil)
		assert.NoError(t, err)
		assert.Equal(t, val.String(), "\"Hello World\"")
	}))

	t.Run("Failure", clientTest(func(ctx context.Context, t *testing.T, handler *mockHandler, client *Client) {
		handler.
			On("executeScriptAtBlockID", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return("", HTTPError{
				Url:     "/",
				Code:    400,
				Message: "bad request",
			})

		_, err := client.ExecuteScriptAtBlockID(ctx, flow.HexToID("0x1"), nil, nil)
		assert.EqualError(t, err, "bad request")
	}))
}

func TestBaseClient_GetEvents(t *testing.T) {
	const handlerName = "getEvents"

	t.Run("Get For Height Range", clientTest(func(ctx context.Context, t *testing.T, handler *mockHandler, client *Client) {
		httpEvents := unittest.BlockEventsFlowFixture(flow.EventEncodingVersionJSONCDC)
		expectedEvents, err := convert.ToBlockEvents([]models.BlockEvents{httpEvents}, nil)
		const eType = "A.Foo.Bar"
		handler.
			On(handlerName, mock.Anything, eType, "0", "5", []string(nil)).
			Return([]models.BlockEvents{httpEvents}, nil)

		events, err := client.GetEventsForHeightRange(ctx, eType, 0, 5)
		assert.NoError(t, err)
		assert.Equal(t, events, expectedEvents)
	}))

	t.Run("Get For Block IDs", clientTest(func(ctx context.Context, t *testing.T, handler *mockHandler, client *Client) {
		httpEvents := unittest.BlockEventsFlowFixture(flow.EventEncodingVersionJSONCDC)
		expectedEvents, err := convert.ToBlockEvents([]models.BlockEvents{httpEvents}, nil)
		const eType = "A.Foo.Bar"
		handler.
			On(handlerName, mock.Anything, eType, "", "", []string{expectedEvents[0].BlockID.String()}).
			Return([]models.BlockEvents{httpEvents}, nil)

		events, err := client.GetEventsForBlockIDs(ctx, eType, []flow.Identifier{expectedEvents[0].BlockID})
		assert.NoError(t, err)
		assert.Equal(t, events, expectedEvents)
	}))

	t.Run("Get For Block IDs Not Found", clientTest(func(ctx context.Context, t *testing.T, handler *mockHandler, client *Client) {
		const eType = "A.Foo.Bar"
		id := test.IdentifierGenerator().New()
		handler.
			On(handlerName, mock.Anything, eType, "", "", []string{id.String()}).
			Return([]models.BlockEvents{}, nil)

		events, err := client.GetEventsForBlockIDs(ctx, eType, []flow.Identifier{id})
		assert.NoError(t, err)
		assert.Equal(t, events, []flow.BlockEvents{})
	}))

	t.Run("Failure", clientTest(func(ctx context.Context, t *testing.T, handler *mockHandler, client *Client) {
		handler.
			On(handlerName, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(nil, HTTPError{
				Url:     "/",
				Code:    400,
				Message: "bad request",
			})

		e, err := client.GetEventsForBlockIDs(ctx, "A.Foo", []flow.Identifier{flow.HexToID("0x1")})
		assert.EqualError(t, err, "bad request")
		assert.Nil(t, e)
	}))

	t.Run("Get For Height Range - Invalid Range", clientTest(func(ctx context.Context, t *testing.T, handler *mockHandler, client *Client) {
		tests := []struct {
			in  []uint64
			err string
		}{
			{in: []uint64{0, 0}, err: "must provide start and end height range"},
			{in: []uint64{5, 0}, err: "start height (5) must be smaller than end height (0)"},
		}

		for _, v := range tests {
			events, err := client.GetEventsForHeightRange(ctx, "A.Foo.Bar", v.in[0], v.in[1])
			assert.EqualError(t, err, v.err)
			assert.Nil(t, events)
		}
	}))

}
