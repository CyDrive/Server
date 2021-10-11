# Generate proto-defined things
sh uproto.sh
echo

root=$(pwd)

go mod tidy

cd $root/master/bin
echo "build master..."
go build -race -o master

echo

cd $root/node
echo "build node..."
pwd
