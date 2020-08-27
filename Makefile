.PHONY: build clean deploy gomodgen

buildForAWS: gomodgen
	export GO111MODULE=on
	env GOOS=linux go build -ldflags="-s -w" -o bin/initiate initiate/initiate.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/fetch fetch/fetch.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/database database/database.go

buildForOSX: gomodgen
	export GO111MODULE=on
	env GOOS=darwin go build -ldflags="-s -w" -o bin/initiate initiate/initiate.go
	env GOOS=darwin go build -ldflags="-s -w" -o bin/fetch fetch/fetch.go
	env GOOS=darwin go build -ldflags="-s -w" -o bin/database database/database.go

clean:
	rm -rf ./bin ./vendor Gopkg.lock

test: clean buildForOSX
	go test -covermode count -coverprofile cover.out ./...

deploy-all: clean buildForAWS
	sls deploy --verbose --config serverless-initiate.yml
	sls deploy --verbose --config serverless-fetch.yml
	sls deploy --verbose --config serverless-database.yml

deploy-initate:
	export GO111MODULE=on
	env GOOS=linux go build -ldflags="-s -w" -o bin/initiate initiate/initiate.go
	sls deploy --verbose --config serverless-initiate.yml

deploy-fetch:
	export GO111MODULE=on
	env GOOS=linux go build -ldflags="-s -w" -o bin/fetch fetch/fetch.go
	sls deploy --verbose --config serverless-fetch.yml

deploy-database:
	export GO111MODULE=on
	env GOOS=linux go build -ldflags="-s -w" -o bin/database database/database.go
	sls deploy --verbose --config serverless-database.yml

gomodgen:
	chmod u+x gomod.sh
	./gomod.sh
