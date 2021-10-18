# Generate proto-defined things
sh uproto.sh
echo

root=$(pwd)

go mod tidy

cd $root/master/bin
echo "build master..."
go build -race -o master

echo

cd $root/node/bin
echo "build node..."
go build -race -o node
