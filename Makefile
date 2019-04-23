all: build

build:
	go install seeder

run: build
	./bin/seeder run
