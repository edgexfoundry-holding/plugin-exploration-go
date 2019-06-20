//
// Copyright (c) 2019 Intel Corporation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package foo

import (
	"log"
	"net/rpc"

	"github.com/hashicorp/go-plugin"
)

//  FooGRPCClient

type FooRPCClient struct {
	client *rpc.Client
}

func (m *FooRPCClient) Put(message string) error {
	var resp interface{}

	err := m.client.Call("Plugin.Put", message, &resp)
	return err
}

func (m *FooRPCClient) Get() (string, error) {
	var resp string
	err := m.client.Call("Plugin.Get", new(interface{}), &resp)
	return resp, err
}

func (m *FooRPCClient) DoCheck() (bool, error) {
	var resp bool
	err := m.client.Call("Plugin.DoCheck", new(interface{}), &resp)
	return resp, err
}

func (m *FooRPCClient) Close() error {
	pluginClient.Kill()
	return nil
}

func (m *FooRPCClient) WaitForIt(it It) error {
	var resp interface{}
	err := m.client.Call("Plugin.WaitForIt", it, &resp)
	return err
}

func (m *FooRPCClient) SetPerson(person Person) (int, error) {
	var resp int
	err := m.client.Call("Plugin.SetPerson", person, &resp)
	return resp, err
}

func (m *FooRPCClient) GetPerson(id int) (Person, error) {
	var resp Person
	err := m.client.Call("Plugin.GetPerson", id, &resp)
	return resp, err
}

//FooRPCServer

type FooRPCServer struct {
	// This is the real implementation
	Impl Client
}

func (s *FooRPCServer) Put(args string, resp *interface{}) error {
	return s.Impl.Put(args)
}

func (s *FooRPCServer) Get(args interface{}, resp *string) error {
	var err error
	*resp, err = s.Impl.Get()
	return err
}
func (s *FooRPCServer) DoCheck(args interface{}, resp *bool) error {
	var err error
	*resp, err = s.Impl.DoCheck()
	return err
}

func (m *FooRPCServer) Close(args interface{}, resp *interface{}) error {
	err := m.Impl.Close()
	return err
}

func (m *FooRPCServer) WaitForIt(args It, resp *interface{}) error {
	return m.Impl.WaitForIt(args)
}

func (s *FooRPCServer) SetPerson(args Person, resp *int) error {
	var err error
	*resp, err = s.Impl.SetPerson(args)
	return err
}

func (m *FooRPCServer) GetPerson(args int, resp *Person) error {
	var err error

	*resp, err = m.Impl.GetPerson(args)
	return err
}

// FooRPCPlugin

type FooRPCPlugin struct {
	// Concrete implementation, written in Go. This is only used for hashicorp
	// that are written in Go.
	Impl Client
}

func (p *FooRPCPlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &FooRPCServer{Impl: p.Impl}, nil
}

func (FooRPCPlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &FooRPCClient{client: c}, nil
}

// Dummy It impl
type ItRPCClient struct {
}

func (m ItRPCClient) ThisIsIt(message string) error {
	log.Println("This is it from client...")

	return nil
}

// ItGRPCServer

type ItRPCServer struct {
	// This is the real implementation
	Impl It
}

func (m ItRPCServer) ThisIsIt(args string, resp *interface{}) error {
	log.Println("This is it from server...")
	return nil
}
