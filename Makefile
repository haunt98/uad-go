.PHONY: all test test-color coverage coverage-cli coverate-html lint format build clean upstream

all:
	$(MAKE) upstream
	go mod tidy
	$(MAKE) test-color
	$(MAKE) lint
	$(MAKE) format
	$(MAKE) build
	$(MAKE) clean

test:
	go test -race -failfast ./...

test-color:
	go install github.com/haunt98/go-test-color@latest
	go-test-color -race -failfast ./...

coverage:
	go test -coverprofile=coverage.out ./...

coverage-cli:
	$(MAKE) coverage
	go tool cover -func=coverage.out

coverage-html:
	$(MAKE) coverage
	go tool cover -html=coverage.out

lint:
	golangci-lint run ./...

format:
	go install github.com/haunt98/gofimports/cmd/gofimports@latest
	go install mvdan.cc/gofumpt@latest
	gofimports -w --company github.com/make-go-great,github.com/haunt98 .
	gofumpt -w -extra .
	deno fmt uad_lists.json

build:
	$(MAKE) clean
	go build -o guad .

clean:
	rm -rf guad

upstream:
	curl https://raw.githubusercontent.com/0x192/universal-android-debloater/main/resources/assets/uad_lists.json --output uad_lists.json
	curl https://raw.githubusercontent.com/MuntashirAkon/android-debloat-list/master/aosp.json --output adl_aosp.json
	curl https://raw.githubusercontent.com/MuntashirAkon/android-debloat-list/master/carrier.json --output adl_carrier.json
	curl https://raw.githubusercontent.com/MuntashirAkon/android-debloat-list/master/google.json --output adl_google.json
	curl https://raw.githubusercontent.com/MuntashirAkon/android-debloat-list/master/misc.json --output adl_misc.json
	curl https://raw.githubusercontent.com/MuntashirAkon/android-debloat-list/master/oem.json --output adl_oem.json
