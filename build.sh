# Generate proto-defined things
sh uproto.sh

root=$(pwd)

go mod tidy

cd $root/master/bin
go build -race -o master

cd $root/node
pwd
