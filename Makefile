DEFAULT:=./cmd
OUTPUT:=./bin/gouno
default: build

build:
	go build -gcflags "-N -l" -o $(OUTPUT) $(DEFAULT)
	chmod +x $(OUTPUT)