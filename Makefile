# One of poller | migrate
APP := poller

version := `git describe --tags`
build_time := `date +%FT%T%z`

executable := iap-polling
src_dir := ./cmd/poller/

ifeq ($(APP), migrate)
    executable := migrate-receipt
    src_dir := ./cmd/migrate/
endif

ldflags := -ldflags "-w -s -X main.version=${version} -X main.build=${build_time}"

build_dir := build

dev_executable := $(build_dir)/$(executable)
linux_executable := $(build_dir)/linux/$(executable)

goos := GOOS=linux GOARCH=amd64

.PHONY: dev
dev :
	go build -o $(dev_executable) $(ldflags) -v $(src_dir)

# Cross compiling linux on for dev.
.PHONY: linux
linux :
	$(goos) go build -o $(linux_executable) $(ldflags) -v $(src_dir)

.PHONY: publish
publish :
	rsync -v $(linux_executable) tk11:/home/node/go/bin/

.PHONY: restart
restart :
	ssh ucloud supervisorctl restart $(executable)

.PHONY: deploy
deploy : linux publish restart
	@echo "deploy success"

.PHONY: clean
clean :
	go clean -x
	rm build/*
