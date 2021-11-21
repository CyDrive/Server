# CyDrive

## Dev Environment
You can dev within the devcontainer, without installing any requirements.

**NOTE: You may need set proxy for git, or something else, the container would resolve `host.docker.internal` to the host address like:**
```shell
git config --global http.proxy [http/socks5]://host.docker.internal:[port]
git config --global https.proxy [http/socks5]://host.docker.internal:[port]
```

## Requirements

- protobuf-compiler: You need to install protobuf compiler if you have to modify the proto files, please follow this [installation guide](https://grpc.io/docs/protoc-installation/). For generating different language codes, you may need another plugins, the makefile of this repo **only install the required plugins of itself**
- proto-gen-go: 
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

You can compile only the master/node, and even only one version, the targets are named by: `component_os`. For example, you can compile the master running on linux by:
```
make master_linux
```

Or compile the node on all OS:
```
make node_all
```

Note that these targets wouldn't update the proto-defined models and enums.


