build: gospt

gospt: $(shell find . -name '*.go')
	go build -o gospt .

completions:
	mkdir -p completions
	gospt completion zsh > completions/_gospt
	gospt completion bash > completions/gospt
	gospt completion fish > completions/gospt.fish

run:
	go run main.go

tidy:
	go mod tidy

clean:
	rm -f gospt
	rm -rf completions

uninstall:
	rm -f /usr/bin/gospt
	rm -f /usr/share/zsh/site-functions/_gospt
	rm -f /usr/share/bash-completion/completions/gospt
	rm -f /usr/share/fish/vendor_completions.d/gospt.fish

install:
	cp gospt /usr/bin
	gospt completion zsh > /usr/share/zsh/site-functions/_gospt
	gospt completion bash > /usr/share/bash-completion/completions/gospt
	gospt completion fish > /usr/share/fish/vendor_completions.d/gospt.fish
