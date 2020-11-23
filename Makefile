version := `git tag -l --sort=-v:refname | head -n 1`
build_time := `date +%FT%T%z`

app_name := iap-migrate

ldflags := -ldflags "-w -s -X main.version=${version} -X main.build=${build_time}"

build_dir := build

executable := $(build_dir)/$(app_name)
src_dir := .

.PHONY: build
build :
	$(goos) go build -o $(executable) $(ldflags) -v $(src_dir)

.PHONY: run
run :
	./$(executable)

.PHONY: clean
clean :
	go clean -x
	rm -rf build/*
