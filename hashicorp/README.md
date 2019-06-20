# HashiCorp plugin exploration

The code here is exploring the use of HashiCorp plugins using both RPC and GRPC.

The goal was to determine the BLANK of implementing plugins for our abstract APIs for Registry, Message Bus, etc. Thus the structure of the abstract Foo interface and factory follow that implement for go-mod-registry & go-mod-messaging

Use the makefile to build and run the exploration code
 * make build
 * make run-rpc
 * make run-grpc
 
 The following are the initial finding from this exploration