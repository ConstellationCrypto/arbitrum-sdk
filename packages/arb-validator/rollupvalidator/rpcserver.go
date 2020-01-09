/*
 * Copyright 2019-2020, Offchain Labs, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package rollupvalidator

import (
	"context"
	"github.com/offchainlabs/arbitrum/packages/arb-validator/arbbridge"
	"net/http"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/offchainlabs/arbitrum/packages/arb-validator/structures"
)

//go:generate bash -c "protoc -I$(go list -f '{{ .Dir }}' -m github.com/offchainlabs/arbitrum/packages/arb-validator) -I. --tstypes_out=../../arb-provider-ethers/src/lib --go_out=paths=source_relative,plugins=grpc:. *.proto"
// Server provides an interface for interacting with a a running coordinator
type RPCServer struct {
	*Server
}

// NewServer returns a new instance of the Server class
func NewRPCServer(
	auth *bind.TransactOpts,
	client arbbridge.ArbClient,
	rollupAddress common.Address,
	codeFile string,
	config structures.ChainParams,
) (*RPCServer, error) {
	server, err := NewServer(auth, client, rollupAddress, codeFile, config)
	return &RPCServer{server}, err
}

// FindLogs takes a set of parameters and return the list of all logs that match the query
func (m *RPCServer) FindLogs(r *http.Request, args *FindLogsArgs, reply *FindLogsReply) error {
	ret, err := m.Server.FindLogs(context.Background(), args)
	if ret != nil {
		*reply = *ret
	}
	return err
}

// GetMessageResult returns the value output by the VM in response to the message with the given hash
func (m *RPCServer) GetMessageResult(r *http.Request, args *GetMessageResultArgs, reply *GetMessageResultReply) error {
	ret, err := m.Server.GetMessageResult(context.Background(), args)
	if ret != nil {
		*reply = *ret
	}
	return err
}

// GetAssertionCount returns the total number of finalized assertions
func (m *RPCServer) GetAssertionCount(r *http.Request, args *GetAssertionCountArgs, reply *GetAssertionCountReply) error {
	ret, err := m.Server.GetAssertionCount(context.Background(), args)
	if ret != nil {
		*reply = *ret
	}
	return err
}

// GetVMInfo returns current metadata about this VM
func (m *RPCServer) GetVMInfo(r *http.Request, args *GetVMInfoArgs, reply *GetVMInfoReply) error {
	ret, err := m.Server.GetVMInfo(context.Background(), args)
	if ret != nil {
		*reply = *ret
	}
	return err
}

// CallMessage takes a request from a client to process in a temporary context and return the result
func (m *RPCServer) CallMessage(r *http.Request, args *CallMessageArgs, reply *CallMessageReply) error {
	ret, err := m.Server.CallMessage(context.Background(), args)
	if ret != nil {
		*reply = *ret
	}
	return err
}
