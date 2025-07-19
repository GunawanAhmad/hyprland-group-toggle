BINARY_NAME=hyprland-group-toggle
DESTDIR=/usr/local/bin

build: main.go
	go build -o bin/$(BINARY_NAME) .

install: build
	sudo install -Dm755 bin/$(BINARY_NAME) $(DESTDIR)/$(BINARY_NAME)

run: main.go
	go run main.go
