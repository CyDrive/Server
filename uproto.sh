#!/bin/bash

# You need protoc-gen-gorm to generate the structs with gorm tags
# set environment variable $proto_dir to the dir of depended proto files of protoc-gen-gorm
# follow https://github.com/infobloxopen/protoc-gen-gorm to install protoc-gen-gorm
# note: new go version may download the repo into $GOPATH/pkg/mod

set -e

proto_dir="third_party/proto"
echo "proto_dir=$proto_dir"

echo "protoc rpc..."
protoc --proto_path=. \
    --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    rpc/manage.proto \
    rpc/file_stream.proto

echo "protoc models..."
protoc --proto_path=. -I$proto_dir \
    --go_out=. --go_opt=paths=source_relative --gorm_out="engine=mysql:." --gorm_opt=paths=source_relative \
    --csharp_out=models \
    models/account.proto models/file_info.proto models/http_models.proto models/message.proto

echo "protoc enums..."
protoc --proto_path=. \
    --go_out=. --go_opt=paths=source_relative \
    --csharp_out=consts \
    consts/enums.proto