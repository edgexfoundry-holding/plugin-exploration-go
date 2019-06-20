# HashiCorp plugin exploration

The code here is exploring the use of HashiCorp plug-ins using both RPC and GRPC.

The goal was to determine the viability of implementing plug-ins for the EdgeX abstract APIs for Registry, Message Bus, etc. Thus the structure of the abstract Foo interface and factory follow that implemented for go-mod-registry & go-mod-messaging

#### Use the makefile to build and run the exploration code

 * make clean
 * make build
 * make run-rpc
 * make run-grpc

#### Updating the protocol buffers file

- If you update the my-abstraction/proto/foo.proto run the following command from the my-abstraction directory  to compile it.

  `protoc -I proto/ proto/foo.proto --go_out=plugins=grpc:proto/`

  This command requires installation of the protoc compiler and the go grpc:proto plugin 

#### The following are the initial finding from this exploration:

   - They work, based on proven rpc/grpc protocols
   - Complex to implement a large interface, especially with grpc due to protocol buffers file and requires protoc compiler and go plug-in for generating the go implementation of the  proto
		- Complexity is on the plug-in definition
		- Actual plug-in implement is same as a regular interface implement
   - If plug-in dies, nothing in frame work will restart it
		- get error when calling the interface
   - Plug-ins don't always exit when service exits
         - This may have been caused by debugging. All runs from command line cleaned up even it the service (client) panics
   - Old plug-in still runs if missing implementation of interface
		- Get error when calling missing interface. 
		- rpc error message is more concise