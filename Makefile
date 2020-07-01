.PHONY: build clean deploy gomodgen

buildForAWS: gomodgen
	export GO111MODULE=on
	env GOOS=linux go build -ldflags="-s -w" -o bin/initiate initiate/initiate.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/fetch fetch/fetch.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/database database/database.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/cleanup cleanup/cleanup.go

buildForOSX: gomodgen
	export GO111MODULE=on
	env GOOS=darwin go build -ldflags="-s -w" -o bin/initiate initiate/initiate.go
	env GOOS=darwin go build -ldflags="-s -w" -o bin/fetch fetch/fetch.go
	env GOOS=darwin go build -ldflags="-s -w" -o bin/database database/database.go
	env GOOS=darwin go build -ldflags="-s -w" -o bin/cleanup cleanup/cleanup.go

clean:
	rm -rf ./bin ./vendor Gopkg.lock

test: clean buildForOSX
	go test -covermode count -coverprofile cover.out ./...

deploy: clean buildForAWS
	sls deploy --verbose

gomodgen:
	chmod u+x gomod.sh
	./gomod.sh
