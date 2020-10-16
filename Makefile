# One of poller | migrate
APP := poller

version := `git tag -l --sort=-v:refname | head -n 1`
build_time := `date +%FT%T%z`

app_name_poller := iap-polling
app_name_migrate := migrate-receipt

ldflags := -ldflags "-w -s -X main.version=${version} -X main.build=${build_time}"

build_dir := build

poller_executable := $(build_dir)/$(app_name_poller)
poller_linux_executable := $(build_dir)/linux/$(app_name_poller)
poller_src_dir := ./cmd/poller/

migrate_executable := $(build_dir)/$(app_name_migrate)
migrate_src_dir := ./cmd/migrate/

goos := GOOS=linux GOARCH=amd64

.PHONY: build-poller
build-poller :
	go build -o $(poller_executable) $(ldflags) -v $(poller_src_dir)

.PHONY: build-migrate
build-migrate :
	go build -o $(migrate_executable) $(ldflags) -v $(migrate_src_dir)

.PHONY: run-migrate
run-migrate :
	./$(migrate_executable) -production -dir="iap_receipts"

.PHONY: linux-poller
linux-poller :
	$(goos) go build -o $(poller_linux_executable) $(ldflags) -v $(poller_src_dir)

.PHONY: publish
publish :
	rsync -v $(poller_linux_executable) tk11:/home/node/go/bin/

.PHONY: deploy
deploy : linux-poller publish
	ssh tk11 supervisorctl restart $(app_name_poller)
	@echo "deploy success"

.PHONY: clean
clean :
	go clean -x
	rm -rf build/*
