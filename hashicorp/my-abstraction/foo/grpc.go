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

	"github.com/hashicorp/go-plugin"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/bogus/my-abstraction/proto"
)

//  FooGRPCClient

type FooGRPCClient struct {
	broker *plugin.GRPCBroker
	client proto.FooClient
}

func (m *FooGRPCClient) Put(message string) error {
	_, err := m.client.Put(context.Background(), &proto.PutRequest{
		Message: message,
	})
	return err
}

func (m *FooGRPCClient) Get() (string, error) {
	resp, err := m.client.Get(context.Background(), &proto.Empty{})
	if err != nil {
		return "", err
	}

	return resp.Message, nil
}

func (m *FooGRPCClient) DoCheck() (bool, error) {
	resp, err := m.client.DoCheck(context.Background(), &proto.Empty{})
	if err != nil {
		return false, err
	}

	return resp.Check, nil
}

func (m *FooGRPCClient) Close() error {
	pluginClient.Kill()
	return nil
}

func (m *FooGRPCClient) WaitForIt(it It) error {
	itServer := &ItGRPCServer{Impl: it}

	var s *grpc.Server
	serverFunc := func(opts []grpc.ServerOption) *grpc.Server {
		s = grpc.NewServer(opts...)
		proto.RegisterItServer(s, itServer)

		return s
	}

	brokerID := m.broker.NextId()
	go m.broker.AcceptAndServe(brokerID, serverFunc)

	_, err := m.client.WaitForIt(context.Background(), &proto.WaitForItRequest{ItServer: brokerID})

	s.Stop()
	return err
}

func (m *FooGRPCClient) SetPerson(person Person) (int, error) {
	response, err := m.client.SetPerson(context.Background(), &proto.SetPersonRequest{
		Person: &proto.Person{
			Age:    int32(person.Age),
			Name:   person.Name,
			Salary: person.Salary,
		},
	})

	if err != nil {
		return 0, err
	}
	return int(response.Id), nil
}

func (m *FooGRPCClient) GetPerson(id int) (Person, error) {
	response, err := m.client.GetPerson(context.Background(), &proto.GetPersonRequest{
		Id: int32(id),
	})

	if err != nil {
		return Person{}, err
	}
	return Person{
		Name:   response.Person.Name,
		Age:    int(response.Person.Age),
		Salary: response.Person.Salary,
	}, nil
}

//FooGRPCServer

type FooGRPCServer struct {
	// This is the real implementation
	Impl   Client
	broker *plugin.GRPCBroker
}

func (m *FooGRPCServer) Put(ctx context.Context, req *proto.PutRequest) (*proto.Empty, error) {
	return &proto.Empty{}, m.Impl.Put(req.Message)
}

func (m *FooGRPCServer) Get(ctx context.Context, req *proto.Empty) (*proto.GetResponse, error) {
	message, err := m.Impl.Get()
	return &proto.GetResponse{Message: message}, err
}

func (m *FooGRPCServer) DoCheck(ctx context.Context, req *proto.Empty) (*proto.DoCheckResponse, error) {
	value, err := m.Impl.DoCheck()
	return &proto.DoCheckResponse{Check: value}, err
}

func (m *FooGRPCServer) Close(ctx context.Context, req *proto.Empty) (*proto.Empty, error) {
	m.Impl.Close()
	return &proto.Empty{}, nil
}

func (m *FooGRPCServer) WaitForIt(ctx context.Context, req *proto.WaitForItRequest) (*proto.Empty, error) {
	conn, err := m.broker.Dial(req.ItServer)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	it := &ItGRPCClient{proto.NewItClient(conn)}
	return &proto.Empty{}, m.Impl.WaitForIt(it)
}

func (m *FooGRPCServer) SetPerson(ctx context.Context, req *proto.SetPersonRequest) (*proto.SetPersonResponse, error) {
	person := Person{
		Name:   req.Person.Name,
		Age:    int(req.Person.Age),
		Salary: req.Person.Salary,
	}
	id, err := m.Impl.SetPerson(person)
	return &proto.SetPersonResponse{Id: int32(id)}, err
}

func (m *FooGRPCServer) GetPerson(ctx context.Context, req *proto.GetPersonRequest) (*proto.GetPersonResponse, error) {

	person, err := m.Impl.GetPerson(int(req.Id))
	return &proto.GetPersonResponse{Person: &proto.Person{
		Name:   person.Name,
		Age:    int32(person.Age),
		Salary: person.Salary,
	}}, err
}

// FooGRPCPlugin

type FooGRPCPlugin struct {
	// GRPCPlugin must still implement the Plugin interface
	plugin.Plugin
	// Concrete implementation, written in Go. This is only used for hashicorp
	// that are written in Go.
	Impl Client
}

func (p *FooGRPCPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	proto.RegisterFooServer(s, &FooGRPCServer{
		Impl:   p.Impl,
		broker: broker,
	})
	return nil
}

func (p *FooGRPCPlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &FooGRPCClient{
		client: proto.NewFooClient(c),
		broker: broker,
	}, nil
}

var _ plugin.GRPCPlugin = &FooGRPCPlugin{}

// ItGRPCClient
type ItGRPCClient struct {
	client proto.ItClient
}

func (m *ItGRPCClient) ThisIsIt(message string) error {
	log.Println("This is it from client...")

	_, err := m.client.ThisIsIt(context.Background(), &proto.ThisIsItRequest{
		Message: message,
	})
	return err
}

// ItGRPCServer

type ItGRPCServer struct {
	// This is the real implementation
	Impl It
}

func (m *ItGRPCServer) ThisIsIt(ctx context.Context, req *proto.ThisIsItRequest) (*proto.Empty, error) {
	log.Println("This is it from server...")
	return &proto.Empty{}, m.Impl.ThisIsIt(req.Message)
}
