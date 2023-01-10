build:
	mkdir -p bin
	go build -o ./bin/gospt .

completions:
	mkdir -p completions
	bin/gospt completion zsh > completions/_gospt
	bin/gospt completion bash > completions/gospt
	bin/gospt completion fish > completions/gospt.fish

run:
	go run main.go

tidy:
	go mod tidy

clean:
	rm -rf bin
	rm -rf completions

uninstall:
	rm -f /usr/bin/gospt
	rm /usr/share/zsh/site-functions/_gospt

install:
	cp bin/gospt /usr/bin
	cp completions/_gospt /usr/share/zsh/site-functions/_gospt
	cp completions/gospt /usr/share/bash-completion/completions/gospt
	cp completions/gospt.fish /usr/share/fish/vendor_completions.d/gospt.fish
