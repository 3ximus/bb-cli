FLAGS=-ldflags="-s -w"

all: build

build:
	go build -o bb $(FLAGS) main.go

install:
	GOBIN=~/.local/bin go install $(FLAGS)

clean:
	rm -f bb

