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
	"fmt"
	"log"
	"os"

	"github.com/bogus/my-abstraction/foo"
	"github.com/bogus/my-abstraction/pkg/types"
)

type myIt struct {
	Name string
	its  chan string
}

func (it *myIt) ThisIsIt(message string) error {
	fmt.Println("Received It... " + message)
	it.its <- message
	return nil
}

const (
	grpcPlugin = "../my-plugin/grpc/grpc"
	rpcPlugin  = "../my-plugin/rpc/rpc"
)

func main() {

	// We don't want to see the plugin logs.
	//log.SetOutput(ioutil.Discard)

	config := types.Config{
		Type:       "plugin",
		Plugin:     grpcPlugin,
		PluginName: "foogrpc",
		Host:       "localhost",
	}

	if len(os.Args) > 1 {
		fmt.Printf("Arg count = %d, using RPC\n", len(os.Args))
		config.Plugin = rpcPlugin
		config.PluginName = "foorpc"
	}

	client, err := foo.NewFooClient(config)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	// Requied to kill the plugin process.
	defer client.Close()

	DoMessageCheck(client)

	log.Println("Calling Get: ")

	message, err := client.Get()
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	fmt.Println("Message received: " + message)

	if err := client.Put("Hello from Foo Plugin!"); err != nil {
		fmt.Println("Error calling Put: " + err.Error())
	}

	DoMessageCheck(client)

	message, err = client.Get()
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	fmt.Println("Message received: " + message)

	if err := client.Put(""); err != nil {
		fmt.Println("Error calling Put: " + err.Error())
	}

	// Uncomment this to test if plug-in process is cleaned up
	//panic(0)

	DoMessageCheck(client)

	if err := client.Put("Hi Foo"); err != nil {
		fmt.Println("Error calling Put: " + err.Error())

	}

	message, err = client.Get()
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	fmt.Println("Message received: " + message)

	if config.PluginName == "foogrpc" {
		it := myIt{
			its: make(chan string, 2),
		}

		go func() {
			if err := client.WaitForIt(&it); err != nil {
				fmt.Println("Error WaitFotIt: " + err.Error())
			}
		}()

		fmt.Println("Waiting for it: ")

		message = <-it.its

		fmt.Println(message)
	}

	var id int
	var person foo.Person

	person, err = client.GetPerson(0)
	if err != nil {
		fmt.Println(err)
	}

	person, err = client.GetPerson(1)
	if err != nil {
		fmt.Println(err)
	}

	id, err = client.SetPerson(foo.Person{
		Name:   "Sam",
		Age:    41,
		Salary: 60.89,
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	fmt.Println(fmt.Sprintf("Person set. Id is: %d", id))

	person, err = client.GetPerson(id)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	fmt.Println(person)

}

func DoMessageCheck(client foo.Client) {
	messageStatus, err := client.DoCheck()
	if err != nil {
		fmt.Println("client.DoCheck() failed: ", err.Error())
		return
	}

	fmt.Println("Message status check: ", messageStatus)
}
