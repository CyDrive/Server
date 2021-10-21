
dir = $(CURDIR)

all: master_all node_all

all_update: update_proto all

check: master_check node_check

pre:
	@go install google.golang.org/protobuf/cmd/protoc-gen-go
	@go install github.com/infobloxopen/protoc-gen-gorm
	@go mod tidy

update_proto: pre
	@sh uproto.sh

master_all: master_windows master_linux master_macos

master_windows: pre
	@echo "building master_windows..."
	@cd $(dir)/master/bin && \
	CGO_ENABLED=0  GOOS=windows go build -o master_windows.exe
	@echo
	
master_linux: pre
	@echo "building master_linux..."
	@cd $(dir)/master/bin && \
	CGO_ENABLED=0  GOOS=linux go build -o master_linux
	@echo

master_macos: pre
	@echo "building master_macos..."
	@cd $(dir)/master/bin && \
	CGO_ENABLED=0  GOOS=darwin go build -o master_macos
	@echo

master_check: pre
	@cd $(dir)/master/bin && \
	go build --race -o master_race_check && \
	rm master_race_check

node_all: node_windows node_linux node_macos

node_windows: pre
	@echo "building node_windows..."
	@cd $(dir)/node/bin && \
	CGO_ENABLED=0  GOOS=windows go build -o node_windows.exe
	@echo

node_linux: pre
	@echo "building node_linux..."
	@cd $(dir)/node/bin && \
	CGO_ENABLED=0  GOOS=linux go build -o node_linux
	@echo

node_macos: pre
	@echo "building node_macos..."
	@cd $(dir)/node/bin && \
	CGO_ENABLED=0  GOOS=darwin go build -o node_macos
	@echo

node_check: pre
	@cd $(dir)/node/bin && \
	go build --race -o node_race_check && \
	rm node_race_check

clean:
	@echo "removing master_windows.exe, master_linux, node_windows.exe, node_linux"
	@rm -f master/bin/master_windows.exe master/bin/master_linux \
	node/bin/node_windows.exe node/bin/node_linux
	@rm -f master/bin/*.log node/bin/*.log