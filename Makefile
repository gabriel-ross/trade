up:
	docker-compose up

docker-build-cd:
	docker build -t foo .

docker-run-cd:
	docker run -p 81:81 foo

run:
	go run cmd/trade/trade.go

run-proxy:
	go run cmd/proxy/main.go

test:
	go test ./...