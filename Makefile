BINARY_NAME=gospt

hashes: build
	sha256sum bin/gospt

srcinfo:
	cd aur && makepkg --printsrcinfo > .SRCINFO

build:
	go build -o bin/${BINARY_NAME}

run:
	go run main.go

tidy:
	go mod tidy

clean:
	rm -rf bin

uninstall:
	rm -f /usr/local/bin/${BINARY_NAME}

install:
	cp bin/${BINARY_NAME} /usr/local/bin

