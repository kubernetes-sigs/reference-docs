genref: main.go types.go
	go build -mod mod -o genref

all: genref
	rm -fr output/md/*
	./genref -o output/md
