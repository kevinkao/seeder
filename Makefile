all: build

build:
	go install seeder

run: build
	./bin/seeder run

sql: build
	./bin/seeder sql database/seeds/insert_base_data.sql