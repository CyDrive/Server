# CyDrive

## Structure
CyDrive consists of two parts: Master and Storage Nodes. They are located in the `master` dir and `node` dir, and the common parts of them are located in the root dir of this repo:
- config: configuration definitions
- consts: constants and enums
- docs: documents
- master: the master
- models: types with only data
- network: some classes that running over network
- node: the storage node
- rpc: the proto files used in the communication between master and storage nodes
- types: some widely used types
- utils: many useful functions

## Models
Many models and enums are defined by protobuf (you can identify them by the suffix with ".pb.go"), DO NOT modify them directly. You should modify the corresponding proto file, and re-generate the go file. You can do this by:
```shell
sh uproto.sh
```
The shell script will update them all.

## Build
You should compile by the makefile:
```
make all
```
or just
```
make
```

The default target wouldn't update the models and enums defined by proto files, make with `all_update` target if you want to update them:
```
make all_update
```

The makefile would compile three executions for both master and storage node, which correspond to three OS (Linux, Windows, macOS) respectively.

You can compile only the master/node, and even only one version, the targets are named by: `name_os`. For example, you can compile the master running on linux by:
```
make master_linux
```

Or compile the node on all OS:
```
make node_all
```

Note that these targets wouldn't update the proto-defined models and enums.


