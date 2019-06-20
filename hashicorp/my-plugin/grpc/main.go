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

package main

import (
	"errors"
	"log"
	"os"
	"time"

	"github.com/bogus/my-abstraction/foo"
	"github.com/hashicorp/go-plugin"
)

// Here is a real implementation of Greeter
type MyFoo struct {
	message string
	person  *foo.Person
}

func (g *MyFoo) Put(message string) error {
	g.message = message
	return nil
}

func (g *MyFoo) Get() (string, error) {
	log.Println("Get called...")
	return "Message" + ": " + g.message, nil
}

func (g *MyFoo) DoCheck() (bool, error) {

	var err error
	if g.message == "" {
		return false, errors.New("message was never set")
	}

	if g.message == "n/a" {
		return false, err
	}

	return true, err
}

func (g *MyFoo) Close() error {
	os.Exit(0)
	return nil
}

func (g *MyFoo) WaitForIt(it foo.It) error {
	// Doen't work are an async go func. Client side listener doesn't respond. Need to have longer life??
	//go func() {
	time.Sleep(2 * time.Second)
	it.ThisIsIt("Here IT is....")

	log.Println("Sent it...")
	//}()

	//log.Println("Started waiting for it...")

	return nil
}

func (g *MyFoo) SetPerson(person foo.Person) (int, error) {
	g.person = &person
	return 1, nil
}

func (g *MyFoo) GetPerson(id int) (foo.Person, error) {

	// Uncomment this to demonstrate plugin crashed and doesn't get restarted
	//os.Exit(0)

	if id != 1 {
		return foo.Person{}, errors.New("Invalid Person Id")
	}

	if g.person == nil {
		return foo.Person{}, errors.New("Person at id=i not set")
	}

	return *g.person, nil
}

// handshakeConfigs are used to just do a basic handshake between
// a plugin and host. If the handshake fails, a user friendly error is shown.
// This prevents users from executing bad hashicorp or executing a plugin
// directory. It is a UX feature, not a security feature.
var handshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "FOO_PLUGIN",
	MagicCookieValue: "foo",
}

func main() {

	myFoo := &MyFoo{
		message: "n/a",
	}

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: handshakeConfig,
		Plugins: map[string]plugin.Plugin{
			"foogrpc": &foo.FooGRPCPlugin{Impl: myFoo},
		},

		// A non-nil value here enables gRPC serving for this plugin...
		GRPCServer: plugin.DefaultGRPCServer,
	})
}
