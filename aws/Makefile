.PHONY: build clean deploy

build:
	cd hello && go build -ldflags="-s -w" -o ../bin/helloBin hello.go && cd ..

clean:
	rm -rf ./bin ./vendor Gopkg.lock

deploy: clean build
	sls deploy --verbose