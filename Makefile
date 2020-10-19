# One of poller | migrate
APP := poller

version := `git tag -l --sort=-v:refname | head -n 1`
build_time := `date +%FT%T%z`

app_name := iap-polling

ldflags := -ldflags "-w -s -X main.version=${version} -X main.build=${build_time}"

build_dir := build

executable := $(build_dir)/$(app_name)
linux_executable := $(build_dir)/linux/$(app_name)
src_dir := ./cmd/poller/

config_file_name := api.toml
goos := GOOS=linux GOARCH=amd64
go_version := go1.15

.PHONY: dev
dev :
	go build -o $(executable) $(ldflags) -v $(src_dir)

.PHONY: run
run :
	./$(executable)

.PHONY: install-go
install-go:
	gvm install $(go_version)
	gvm use $(go_version)

.PHONY: build
build :
	$(goos) go build -o $(linux_executable) $(ldflags) -v $(src_dir)

.PHONY: config
config :
	rsync -v tk11:/home/node/config/$(config_file_name) ./$(build_dir)
	rsync -v ./$(build_dir)/$(config_file_name) ucloud:/home/node/config

.PHONY: publish
publish :
	ssh ucloud "rm -f /home/node/go/bin/$(app_name).bak"
	rsync -v $(executable) bj32:/home/node
	ssh bj32 "rsync -v /home/node/$(app_name) ucloud:/home/node/go/bin/$(app_name).bak"

.PHONY: restart
restart :
	ssh ucloud "cd /home/node/go/bin/ && \mv $(app_name).bak $(app_name)"
	ssh ucloud supervisorctl restart $(app_name)

.PHONY: deploy
deploy : build
	rsync -v $(linux_executable) tk11:/home/node/go/bin/
	ssh tk11 supervisorctl restart $(app_name)
	@echo "deploy success"

.PHONY: clean
clean :
	go clean -x
	rm -rf build/*
