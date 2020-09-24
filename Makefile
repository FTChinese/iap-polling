build_dir := build
poller_name := iap-polling
migrate_name := iap-receipt-migrate

poller_dev_out := $(build_dir)/$(poller_name)
poller_prod_out := $(build_dir)/linux/$(poller_name)

migrate_dev_out := $(build_dir)/$(migrate_name)
migrate_prod_out := $(build_dir)/linux/$(migrate_name)

version := `git describe --tags`
build_time := `date +%FT%T%z`

ldflags := -ldflags "-w -s -X main.version=${version} -X main.build=${build_time}"

linux_poller := GOOS=linux GOARCH=amd64 go build -o $(poller_prod_out) $(ldflags) -v ./cmd/producer/

.PHONY: build poller migrate linux clean
# Development
build :
	go build -o $(poller_dev_out) $(ldflags) -v ./cmd/poller/
	go build -o $(migrate_dev_out) $(ldflags) -v ./cmd/migrate/

# Run development build
poller :
	./$(poller_dev_out)

migrate :
	./$(migrate_dev_out)

# Cross compiling linux on for dev.
linux :
	$(linux_poller)

clean :
	go clean -x
	rm build/*
