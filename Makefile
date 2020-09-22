build_dir := build
producer_name := iap-polling-produer
consumer_name := iap-polling-consumer

producer_dev_out := $(build_dir)/$(producer_name)
producer_prod_out := $(build_dir)/linux/$(producer_name)

consumer_dev_out := $(build_dir)/$(consumer_name)
consumer_prod_out := $(build_dir)/linux/$(consumer_name)

version := `git describe --tags`
build_time := `date +%FT%T%z`

ldflags := -ldflags "-w -s -X main.version=${version} -X main.build=${build_time}"

linux_producer := GOOS=linux GOARCH=amd64 go build -o $(producer_prod_out) $(ldflags) -v ./cmd/producer/
linux_consumer := GOOS=linux GOARCH=amd64 go build -o $(consumer_prod_out) $(ldflags) -v ./cmd/consumer/

.PHONY: build producer consumer linux clean
# Development
build :
	go build -o $(producer_dev_out) $(ldflags) -v ./cmd/producer/
	go build -o $(consumer_dev_out) $(ldflags) -v ./cmd/consumer/

# Run development build
producer :
	./$(producer_dev_out)

consumer :
	./$(consumer_dev_out)

# Cross compiling linux on for dev.
linux :
	$(linux_producer)
	$(linux_consumer)

clean :
	go clean -x
	rm build/*
