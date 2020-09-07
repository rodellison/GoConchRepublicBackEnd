.PHONY: build clean deploy gomodgen

buildForAWS: gomodgen
	export GO111MODULE=on
	env GOOS=linux go build -ldflags="-s -w" -o bin/fetch fetch/fetch.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/database database/database.go

buildForOSX: gomodgen
	export GO111MODULE=on
	env GOOS=darwin go build -ldflags="-s -w" -o bin/fetch fetch/fetch.go
	env GOOS=darwin go build -ldflags="-s -w" -o bin/database database/database.go

buildForWin:
	set GO111MODULE=on
	set GOOS=windows go build -ldflags="-s -w" -o bin\fetch fetch\fetch.go
	set GOOS=windows go build -ldflags="-s -w" -o bin\database database\database.go

clean:
	rm -rf ./bin ./vendor Gopkg.lock

cleanWin:
	del /s  /q bin\*.*

test: clean buildForOSX
	go test -covermode count -coverprofile cover.out ./...

testWin: cleanWin buildForWin
	go test -covermode count -coverprofile cover.out ./...

deploy: clean buildForAWS
	sls deploy --verbose --config serverless.yml

gomodgen:
	chmod u+x gomod.sh
	./gomod.sh
