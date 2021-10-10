# You need protoc-gen-gorm to generate the structs with gorm tags
# set environment variable $proto_dir to the dir of depended proto files of protoc-gen-gorm
proto_dir="$GOPATH/src/github.com/infobloxopen/protoc-gen-gorm/proto"
echo "proto_dir=$proto_dir"

echo "protoc rpc..."
protoc --proto_path=. \
    --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    rpc/manage.proto

echo "protoc models..."
protoc --proto_path=. -I$proto_dir \
    --go_out=. --go_opt=paths=source_relative --gorm_out="engine=mysql:." --gorm_opt=paths=source_relative \
    --csharp_out=model \
    model/storage.proto model/account.proto model/message.proto

echo "protoc enums..."
protoc --proto_path=. \
    --go_out=. --go_opt=paths=source_relative \
    --csharp_out=consts \
    consts/status_code.proto consts/message_type.proto