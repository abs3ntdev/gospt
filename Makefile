hashes: build
	sha256sum gospt

srcinfo:
	cd aur && makepkg -g >> PKGBUILD

build:
	go build -o gospt .

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

