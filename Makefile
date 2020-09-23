build_dir := build
producer_name := iap-polling-produer
migrate_name := iap-receipt-migrate

producer_dev_out := $(build_dir)/$(producer_name)
producer_prod_out := $(build_dir)/linux/$(producer_name)

migrate_dev_out := $(build_dir)/$(migrate_name)
migrate_prod_out := $(build_dir)/linux/$(migrate_name)

version := `git describe --tags`
build_time := `date +%FT%T%z`

ldflags := -ldflags "-w -s -X main.version=${version} -X main.build=${build_time}"

linux_producer := GOOS=linux GOARCH=amd64 go build -o $(producer_prod_out) $(ldflags) -v ./cmd/producer/

.PHONY: build producer consumer linux clean
# Development
build :
	go build -o $(producer_dev_out) $(ldflags) -v ./cmd/producer/
	go build -o $(migrate_dev_out) $(ldflags) -v ./cmd/migrate/

# Run development build
producer :
	./$(producer_dev_out)

migrate :
	./$(migrate_dev_out)

# Cross compiling linux on for dev.
linux :
	$(linux_producer)

clean :
	go clean -x
	rm build/*
