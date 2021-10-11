
dir = $(CURDIR)

pre:
	sh uproto.sh
	echo

	go mod tidy

all: master_all

windows: master

master_all: master_windows master_linux

master_windows: pre
	echo "build master..."
	cd $(dir)/master/bin && \
	CGO_ENABLED=0  GOOS=windows go build -o master_windows.exe
	
master_linux: pre
	echo "build master..."
	cd $(dir)/master/bin && \
	CGO_ENABLED=0  GOOS=linux go build -o master_linux