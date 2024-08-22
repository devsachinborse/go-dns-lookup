build:
	@go build -o bin/dns-lookup .

run:build
	@./bin/dns-lookup


	