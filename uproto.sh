#!/bin/bash

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
    --go_out=. --go_opt=paths=source_relative \
    --csharp_out=models \
    models/account.proto models/file_info.proto models/message.proto models/http_models.proto

echo "protoc enums..."
protoc --proto_path=. \
    --go_out=. --go_opt=paths=source_relative \
    --csharp_out=consts \
    consts/enums.proto